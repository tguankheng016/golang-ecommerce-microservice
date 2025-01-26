package configurations

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	permission_service "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions/grpc_client/protos"
	permission_grpc_service "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions/grpc_client/services"
	identity_service "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt/grpc_client/protos"
	identity_grpc_service "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt/grpc_client/services"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/config"
	productGrpc "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/grpc_server/protos"
	productGrpcService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/grpc_server/services"
)

func ConfigureProductGrpcServer(db *pgxpool.Pool, grpcServer *grpc.GrpcServer) {
	productGrpcService := productGrpcService.NewProductGrpcServerService(db)
	productGrpc.RegisterProductGrpcServiceServer(grpcServer.Grpc, productGrpcService)
}

func ConfigIdentityGrpcClientService(clientFactory *grpc.GrpcClientFactory, clientAddress *config.GrpcAddress) identity_service.IdentityGrpcServiceClient {
	return identity_grpc_service.NewIdentityGrpcClientService(clientFactory, clientAddress.IdentityAddress)
}

func ConfigPermissionGrpcClientService(clientFactory *grpc.GrpcClientFactory, clientAddress *config.GrpcAddress) permission_service.PermissionGrpcServiceClient {
	return permission_grpc_service.NewPermissionGrpcClientService(clientFactory, clientAddress.IdentityAddress)
}
