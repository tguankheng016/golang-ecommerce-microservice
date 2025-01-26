package services

import (
	"context"

	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	userService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/services"
)

type userPermissionDbManager struct {
	userRolePermissionManager userService.IUserRolePermissionManager
}

func NewPermissionDbManager(db userService.IUserRolePermissionManager) permissions.IPermissionDbManager {
	return &userPermissionDbManager{
		userRolePermissionManager: db,
	}
}

func (m *userPermissionDbManager) GetGrantedPermissionsFromDb(ctx context.Context, userId int64) (map[string]struct{}, error) {
	return m.userRolePermissionManager.SetUserPermissions(ctx, userId)
}
