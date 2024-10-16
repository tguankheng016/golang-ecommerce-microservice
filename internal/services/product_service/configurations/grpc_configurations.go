package configurations

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
)

type GrpcAddress struct {
	IdentityAddress string `mapstructure:"identityAddress"`
}

func ConfigGrpcClients(clientFactory *grpc.GrpcClientFactory, addresses *GrpcAddress) {
	clientFactory.AddClient(addresses.IdentityAddress)
}
