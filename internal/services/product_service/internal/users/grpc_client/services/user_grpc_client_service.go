package services

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/config"
	userGrpc "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/grpc_client/protos"
)

func NewUserGrpcClientService(clientFactory *grpc.GrpcClientFactory, clientAddress *config.GrpcAddress) userGrpc.UserGrpcServiceClient {
	clientFactory.AddClient(clientAddress.IdentityAddress)
	conn := clientFactory.Clients[clientAddress.IdentityAddress].GetGrpcConnection()
	return userGrpc.NewUserGrpcServiceClient(conn)
}
