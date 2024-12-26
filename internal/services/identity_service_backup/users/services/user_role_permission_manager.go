package services

import (
	"context"
	"strings"
	"time"

	"github.com/anchore/go-logger"
	redisCli "github.com/redis/go-redis/v9"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/redis"
	identityConsts "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/constants"
	roleModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/models"
	userModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	DefaultCacheExpiration = 1 * time.Hour
)

type IUserRolePermissionManager interface {
	IsGranted(ctx context.Context, userId int64, permissionName string) (bool, error)
	SetUserPermissions(ctx context.Context, userId int64) (map[string]struct{}, error)
	SetRolePermissions(ctx context.Context, roleId int64) (map[string]struct{}, error)
	RemoveUserRoleCaches(ctx context.Context, userId int64)
}

type userRolePermissionManager struct {
	db     *gorm.DB
	client redisCli.UniversalClient
}

func NewUserRolePermissionManager(db *gorm.DB, client redisCli.UniversalClient) IUserRolePermissionManager {
	return &userRolePermissionManager{
		db:     db,
		client: client,
	}
}

func (u *userRolePermissionManager) IsGranted(ctx context.Context, userId int64, permissionName string) (bool, error) {
	grantedPermissions, err := u.SetUserPermissions(ctx, userId)
	if err != nil {
		return false, err
	}

	_, ok := grantedPermissions[permissionName]
	return ok, nil
}

func (u *userRolePermissionManager) SetUserPermissions(ctx context.Context, userId int64) (map[string]struct{}, error) {
	var user userModel.User
	if err := u.db.Model(&userModel.User{}).Preload("Roles").First(&user, userId).Error; err != nil {
		return nil, err
	}

	// Cache User Role Ids
	roleIds := make([]int64, 0)
	for _, role := range user.Roles {
		roleIds = append(roleIds, role.Id)
	}

	userRoleCacheItem, err := redis.MarshalCacheItem(permissions.NewUserRoleCacheItem(userId, roleIds))
	if err != nil {
		return nil, err
	}
	if err := u.client.Set(ctx, permissions.GenerateUserRoleCacheKey(userId), userRoleCacheItem, DefaultCacheExpiration).Err(); err != nil {
		// Dont return just log
		logger.Logger.Error("error in setting user role caches", zap.Error(err))
	}

	var userPermissions []userModel.UserRolePermission

	if err := u.db.Model(&userModel.UserRolePermission{}).Where("user_id = ?", userId).Find(&userPermissions).Error; err != nil {
		return nil, err
	}

	grantedPermissions := make(map[string]struct{})
	grantedUserPermissions := make(map[string]struct{})
	prohibitedUserPermissions := make(map[string]struct{})
	for _, permission := range userPermissions {
		if permission.IsGranted {
			grantedUserPermissions[permission.Name] = struct{}{}
			grantedPermissions[permission.Name] = struct{}{}
		} else {
			prohibitedUserPermissions[permission.Name] = struct{}{}
		}
	}

	userPermissionCacheItem, err := redis.MarshalCacheItem(permissions.NewUserPermissionCacheItem(userId, grantedUserPermissions, prohibitedUserPermissions))
	if err != nil {
		return nil, err
	}
	if err := u.client.Set(ctx, permissions.GenerateUserPermissionCacheKey(userId), userPermissionCacheItem, DefaultCacheExpiration).Err(); err != nil {
		// Dont return just log
		logger.Logger.Error("error in setting user permissions caches", zap.Error(err))
	}

	for _, role := range user.Roles {
		rolePermissions, err := u.SetRolePermissions(ctx, role.Id)
		if err != nil {
			return nil, err
		}

		for permission := range rolePermissions {
			if _, ok := prohibitedUserPermissions[permission]; !ok {
				if _, ok := grantedPermissions[permission]; !ok {
					// key does not exists
					grantedPermissions[permission] = struct{}{}
				}
			}
		}
	}

	return grantedPermissions, nil
}

func (u *userRolePermissionManager) SetRolePermissions(ctx context.Context, roleId int64) (map[string]struct{}, error) {
	var role roleModel.Role
	if err := u.db.Model(&roleModel.Role{}).First(&role, roleId).Error; err != nil {
		return nil, err
	}

	isAdmin := strings.EqualFold(role.Name, identityConsts.DefaultAdminRoleName)

	var rolePermissions []userModel.UserRolePermission

	if err := u.db.Model(&userModel.UserRolePermission{}).Where("role_id = ?", roleId).Find(&rolePermissions).Error; err != nil {
		return nil, err
	}

	if isAdmin {
		allPermissions := permissions.GetAppPermissions().Items

		// Get all prohibited permissions for admin role
		prohibitedPermissions := make(map[string]struct{})
		for _, permission := range rolePermissions {
			if !permission.IsGranted {
				prohibitedPermissions[permission.Name] = struct{}{}
			}
		}

		// Excluded prohibited permissions
		grantedPermissions := make(map[string]struct{})
		for key := range allPermissions {
			if _, ok := prohibitedPermissions[key]; !ok {
				grantedPermissions[key] = struct{}{}
			}
		}

		rolePermissionCacheItem, err := redis.MarshalCacheItem(permissions.NewRolePermissionCacheItem(roleId, grantedPermissions))
		if err != nil {
			return nil, err
		}
		if err := u.client.Set(ctx, permissions.GenerateRolePermissionCacheKey(roleId), rolePermissionCacheItem, DefaultCacheExpiration).Err(); err != nil {
			// Dont return just log
			logger.Logger.Error("error in setting role permission caches", zap.Error(err))
		}

		return grantedPermissions, nil
	} else {
		grantedPermissions := make(map[string]struct{})
		for _, permission := range rolePermissions {
			if permission.IsGranted {
				grantedPermissions[permission.Name] = struct{}{}
			}
		}

		return grantedPermissions, nil
	}
}

func (u *userRolePermissionManager) RemoveUserRoleCaches(ctx context.Context, userId int64) {
	if err := u.client.Del(ctx, permissions.GenerateUserRoleCacheKey(userId)).Err(); err != nil {
		logger.Logger.Error("error in removing user roles caches", zap.Error(err))
	}
}
