package services

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	permission_service "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions/grpc_client/protos"
)

func NewPermissionGrpcClientService(clientFactory *grpc.GrpcClientFactory, address string) permission_service.PermissionGrpcServiceClient {
	clientFactory.AddClient(address)
	conn := clientFactory.Clients[address].GetGrpcConnection()
	return permission_service.NewPermissionGrpcServiceClient(conn)
}
