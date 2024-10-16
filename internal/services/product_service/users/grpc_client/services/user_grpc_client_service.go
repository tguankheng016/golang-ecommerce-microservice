package services

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/config"
	user_service "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/users/grpc_client/protos"
)

func NewUserGrpcClientService(clientFactory *grpc.GrpcClientFactory, clientAddress *config.GrpcAddress) user_service.UserGrpcServiceClient {
	clientFactory.AddClient(clientAddress.IdentityAddress)
	conn := clientFactory.Clients[clientAddress.IdentityAddress].GetGrpcConnection()
	return user_service.NewUserGrpcServiceClient(conn)
}
