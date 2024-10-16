package configurations

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	user_service "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/grpc_server/protos"
	user_grpc_service "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/grpc_server/services"
	identity_grpc_service "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/identities/grpc_server/services"
	identity_service "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/identities/grpc_server/protos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	"gorm.io/gorm"
)

func ConfigUserGrpcServer(db *gorm.DB, grpcServer *grpc.GrpcServer) {
	userGrpcService := user_grpc_service.NewUserGrpcServerService(db)
	user_service.RegisterUserGrpcServiceServer(grpcServer.Grpc, userGrpcService)
}

func ConfigIdentityGrpcServer(securityStampvalidator jwt.IJwtSecurityStampDbValidator, tokenKeyValidator jwt.IJwtTokenKeyDbValidator, grpcServer *grpc.GrpcServer) {
	identityGrpcService := identity_grpc_service.NewIdentityGrpcServerService(securityStampvalidator, tokenKeyValidator)
	identity_service.RegisterIdentityGrpcServiceServer(grpcServer.Grpc, identityGrpcService)
}