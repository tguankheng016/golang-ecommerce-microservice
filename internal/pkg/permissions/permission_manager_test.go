package permissions

import (
	"context"
	"fmt"
	"testing"

	miniRedis "github.com/alicebob/miniredis/v2"
	redisCli "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/redis"
)

var (
	ctx = context.TODO()
)

type mockPermissionDbManager struct {
	mock.Mock
}

func TestValidatePermissionName(t *testing.T) {
	tests := []struct {
		name               string
		grantedPermissions []string
		wantErr            bool
		wantErrMessage     string
	}{
		{
			name:               "valid permission name",
			grantedPermissions: []string{"Pages.Administration.Users", "Pages.Administration.Roles"},
			wantErr:            false,
			wantErrMessage:     "",
		},
		{
			name:               "invalid permission name",
			grantedPermissions: []string{"Pages.Administration.Users", "Roles"},
			wantErr:            true,
			wantErrMessage:     "invalid permission name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePermissionName(tt.grantedPermissions)
			assert.Equal(t, err != nil, tt.wantErr)

			if tt.wantErr {
				assert.Equal(t, err.Error(), tt.wantErrMessage)
			}
		})
	}
}

func TestGetGrantedPermissions(t *testing.T) {
	tests := []struct {
		name              string
		userId            int64
		permissionToCheck string
		wantPermissions   map[string]struct{}
		wantIsGranted     bool
	}{
		{
			name:              "get caches user granted permissions",
			userId:            1,
			permissionToCheck: "Pages.Administration.Users",
			wantPermissions: map[string]struct{}{
				"Pages.Administration.Users": {},
				"Pages.Administration.Roles": {},
			},
			wantIsGranted: true,
		},
		{
			name:              "get caches role granted permissions",
			userId:            2,
			permissionToCheck: "Pages.Administration.Users",
			wantPermissions: map[string]struct{}{
				"Pages.Administration.Users": {},
				"Pages.Administration.Roles": {},
			},
			wantIsGranted: true,
		},
		{
			name:              "get caches role and user granted permissions with prohibited permissions",
			userId:            3,
			permissionToCheck: "Pages.Administration.Users",
			wantPermissions: map[string]struct{}{
				"Pages.Administration.Roles": {},
			},
			wantIsGranted: false,
		},
		{
			name:              "get non caches user permissions",
			userId:            4,
			permissionToCheck: "Pages.Administration.Users",
			wantPermissions: map[string]struct{}{
				"Pages.Administration.Users": {},
			},
			wantIsGranted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.userId == 4 {
				fmt.Print("test")
			}
			client := getMockRedisClientByUserId(t, tt.userId)
			dbManager := getMockDbManagerByUserId(tt.userId)
			permissionManager := NewPermissionManager(client, dbManager)

			gotPermissions, err := permissionManager.GetGrantedPermissions(ctx, tt.userId)
			assert.Equal(t, err, nil)
			assert.Equal(t, tt.wantPermissions, gotPermissions)

			isGranted, err := permissionManager.IsGranted(ctx, tt.userId, tt.permissionToCheck)
			assert.Equal(t, err, nil)
			assert.Equal(t, tt.wantIsGranted, isGranted)
		})
	}
}

func getMockRedisClientByUserId(t *testing.T, userId int64) *redisCli.Client {
	// Create a mock Redis client
	s := miniRedis.RunT(t)

	switch userId {
	case 1:
		// No role but have granted some permissions
		// Have caches
		roleIds := make([]int64, 0)
		bytes, _ := redis.MarshalCacheItem(NewUserRoleCacheItem(userId, roleIds))
		s.Set(GenerateUserRoleCacheKey(userId), string(bytes))

		grantedUserPermissions := map[string]struct{}{
			"Pages.Administration.Users": {},
			"Pages.Administration.Roles": {},
		}
		prohibitedUserPermissions := make(map[string]struct{})
		bytes, _ = redis.MarshalCacheItem(NewUserPermissionCacheItem(userId, grantedUserPermissions, prohibitedUserPermissions))
		s.Set(GenerateUserPermissionCacheKey(userId), string(bytes))
	case 2:
		// Have role and role permissions
		// No user permissions
		// Have caches
		var roleId int64 = 1
		roleIds := make([]int64, 0)
		roleIds = append(roleIds, roleId)
		bytes, _ := redis.MarshalCacheItem(NewUserRoleCacheItem(userId, roleIds))
		s.Set(GenerateUserRoleCacheKey(userId), string(bytes))

		grantedRolePermissions := map[string]struct{}{
			"Pages.Administration.Users": {},
			"Pages.Administration.Roles": {},
		}
		bytes, _ = redis.MarshalCacheItem(NewRolePermissionCacheItem(roleId, grantedRolePermissions))
		s.Set(GenerateRolePermissionCacheKey(roleId), string(bytes))

		grantedUserPermissions := make(map[string]struct{})
		prohibitedUserPermissions := make(map[string]struct{})
		bytes, _ = redis.MarshalCacheItem(NewUserPermissionCacheItem(userId, grantedUserPermissions, prohibitedUserPermissions))
		s.Set(GenerateUserPermissionCacheKey(userId), string(bytes))
	case 3:
		// Have role and role permissions and user permissions
		// Have prohibited user permissions
		// Have caches
		var roleId int64 = 1
		roleIds := make([]int64, 0)
		roleIds = append(roleIds, roleId)
		bytes, _ := redis.MarshalCacheItem(NewUserRoleCacheItem(userId, roleIds))
		s.Set(GenerateUserRoleCacheKey(userId), string(bytes))

		grantedRolePermissions := map[string]struct{}{
			"Pages.Administration.Users": {},
			"Pages.Administration.Roles": {},
		}
		bytes, _ = redis.MarshalCacheItem(NewRolePermissionCacheItem(roleId, grantedRolePermissions))
		s.Set(GenerateRolePermissionCacheKey(roleId), string(bytes))

		prohibitedUserPermissions := map[string]struct{}{
			"Pages.Administration.Users": {},
		}
		bytes, _ = redis.MarshalCacheItem(NewUserPermissionCacheItem(userId, make(map[string]struct{}), prohibitedUserPermissions))
		s.Set(GenerateUserPermissionCacheKey(userId), string(bytes))
	case 4:
		// No Caches
	default:
	}

	client := redisCli.NewClient(&redisCli.Options{
		Addr: s.Addr(),
	})

	return client
}

func getMockDbManagerByUserId(userId int64) *mockPermissionDbManager {
	dbManager := new(mockPermissionDbManager)

	switch userId {
	case 4:
		dbManager.On("GetGrantedPermissionsFromDb", ctx, userId).Return(map[string]struct{}{
			"Pages.Administration.Users": {},
		}, nil)
	default:
	}

	return dbManager
}

func (m *mockPermissionDbManager) GetGrantedPermissionsFromDb(ctx context.Context, userId int64) (map[string]struct{}, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(map[string]struct{}), args.Error(1)
}
