package configurations

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	identity_service "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt/grpc_client/protos"
	identity_grpc_service "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt/grpc_client/services"
	permission_service "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions/grpc_client/protos"
	permission_grpc_service "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions/grpc_client/services"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/config"
)

func ConfigIdentityGrpcClientService(clientFactory *grpc.GrpcClientFactory, clientAddress *config.GrpcAddress) identity_service.IdentityGrpcServiceClient {
	return identity_grpc_service.NewIdentityGrpcClientService(clientFactory, clientAddress.IdentityAddress)
}

func ConfigPermissionGrpcClientService(clientFactory *grpc.GrpcClientFactory, clientAddress *config.GrpcAddress) permission_service.PermissionGrpcServiceClient {
	return permission_grpc_service.NewPermissionGrpcClientService(clientFactory, clientAddress.IdentityAddress)
}