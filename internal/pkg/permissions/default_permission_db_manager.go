package permissions

import (
	"context"

	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core/helpers"
	permission_service "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions/grpc_client/protos"
)

type defaultPermissionDbManager struct {
	permissionGrpcClient permission_service.PermissionGrpcServiceClient
}

func NewDefaultPermissionDbManager(permissionGrpcClient permission_service.PermissionGrpcServiceClient) IPermissionDbManager {
	return &defaultPermissionDbManager{
		permissionGrpcClient: permissionGrpcClient,
	}
}

func (d *defaultPermissionDbManager) GetGrantedPermissionsFromDb(ctx context.Context, userId int64) (map[string]struct{}, error) {
	response, err := d.permissionGrpcClient.GetUserPermissions(ctx, &permission_service.GetUserPermissionsRequest{
		UserId : userId,
	})
	if err != nil {
		return nil, err
	}

	return helpers.SliceToMap(response.Permissions), nil
}
