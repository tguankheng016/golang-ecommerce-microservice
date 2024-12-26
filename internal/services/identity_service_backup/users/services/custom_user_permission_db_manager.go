package services

import (
	"context"

	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
)

type customUserPermissionDbManager struct {
	dbManager IUserRolePermissionManager
}

func NewUserPermissionDbManager(db IUserRolePermissionManager) permissions.IPermissionDbManager {
	return &customUserPermissionDbManager{
		dbManager: db,
	}
}

func (m *customUserPermissionDbManager) GetGrantedPermissionsFromDb(ctx context.Context, userId int64) (map[string]struct{}, error) {
	return m.dbManager.SetUserPermissions(ctx, userId)
}
