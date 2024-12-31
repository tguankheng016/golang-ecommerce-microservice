package configurations

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	permissionService "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions/grpc_client/protos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	identityGrpc "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt/grpc_client/protos"
	identityGrpcService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/identities/grpc_server/services"
	userGrpc "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/grpc_server/protos"
	userGrpcService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/grpc_server/services"
)

func ConfigUserGrpcServer(db *pgxpool.Pool, grpcServer *grpc.GrpcServer) {
	userGrpcService := userGrpcService.NewUserGrpcServerService(db)
	userGrpc.RegisterUserGrpcServiceServer(grpcServer.Grpc, userGrpcService)
}

func ConfigIdentityGrpcServer(securityStampvalidator jwt.IJwtSecurityStampDbValidator, tokenKeyValidator jwt.IJwtTokenKeyDbValidator, grpcServer *grpc.GrpcServer) {
	identityGrpcService := identityGrpcService.NewIdentityGrpcServerService(securityStampvalidator, tokenKeyValidator)
	identityGrpc.RegisterIdentityGrpcServiceServer(grpcServer.Grpc, identityGrpcService)
}

func ConfigPermissionGrpcServer(permissionManager permissions.IPermissionDbManager, grpcServer *grpc.GrpcServer) {
	permissionGrpcService := identityGrpcService.NewPermissionGrpcServerService(permissionManager)
	permissionService.RegisterPermissionGrpcServiceServer(grpcServer.Grpc, permissionGrpcService)
}
