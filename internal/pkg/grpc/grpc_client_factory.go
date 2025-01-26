package grpc

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logging"
	"go.uber.org/zap"
)

type GrpcClientFactory struct {
	Clients map[string]GrpcClient
}

func NewGrpcClientFactory() *GrpcClientFactory {
	return &GrpcClientFactory{Clients: make(map[string]GrpcClient)}
}

func (f *GrpcClientFactory) AddClient(address string) error {
	client, err := NewGrpcClient(address)
	if err != nil {
		return err
	}

	f.Clients[address] = client

	return nil
}

func (f *GrpcClientFactory) RemoveClients() {
	for key, client := range f.Clients {
		if err := client.Close(); err != nil {
			logging.Logger.Error("error in closing grpc client "+key, zap.Error(err))
		}
	}
}
