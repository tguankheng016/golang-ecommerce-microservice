package services

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/config"
	productGrpc "github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/products/grpc_client/protos"
)

func NewProductGrpcClientService(clientFactory *grpc.GrpcClientFactory, clientAddress *config.GrpcAddress) productGrpc.ProductGrpcServiceClient {
	clientFactory.AddClient(clientAddress.ProductAddress)
	conn := clientFactory.Clients[clientAddress.ProductAddress].GetGrpcConnection()
	return productGrpc.NewProductGrpcServiceClient(conn)
}
