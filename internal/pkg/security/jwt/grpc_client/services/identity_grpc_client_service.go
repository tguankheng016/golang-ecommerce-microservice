package services

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	identity_service "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt/grpc_client/protos"
)

func NewIdentityGrpcClientService(clientFactory *grpc.GrpcClientFactory, address string) identity_service.IdentityGrpcServiceClient {
	clientFactory.AddClient(address)
	conn := clientFactory.Clients[address].GetGrpcConnection()
	return identity_service.NewIdentityGrpcServiceClient(conn)
}
