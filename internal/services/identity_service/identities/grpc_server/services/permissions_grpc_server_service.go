package services

import (
	"context"

	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core/helpers"
	permission_service "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions/grpc_client/protos"
)

type PermissionGrpcServerService struct {
	permissionManager permissions.IPermissionDbManager
}

func NewPermissionGrpcServerService(permissionManager permissions.IPermissionDbManager) *PermissionGrpcServerService {
	return &PermissionGrpcServerService{
		permissionManager: permissionManager,
	}
}

func (i *PermissionGrpcServerService) GetUserPermissions(ctx context.Context, req *permission_service.GetUserPermissionsRequest) (*permission_service.GetUserPermissionsResponse, error) {
	grantedPermissions, err := i.permissionManager.GetGrantedPermissionsFromDb(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	var result = &permission_service.GetUserPermissionsResponse{
		Permissions: helpers.MapKeysToSlice(grantedPermissions),
	}

	return result, nil
}