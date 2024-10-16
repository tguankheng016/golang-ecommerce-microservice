package configurations

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	user_service "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/grpc_server/protos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/grpc_server/services"
	"gorm.io/gorm"
)

func ConfigUserGrpcServer(db *gorm.DB, grpcServer *grpc.GrpcServer) {
	userGrpcService := services.NewUserGrpcService(db)
	user_service.RegisterUserGrpcServiceServer(grpcServer.Grpc, userGrpcService)
}
