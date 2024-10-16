package users

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/users/grpc_client/services"
	"go.uber.org/fx"
)

var (
	// Module provided to fx
	Module = fx.Module(
		"userfx",
		userProviders,
	)

	userProviders = fx.Options(
		fx.Provide(
			services.NewUserGrpcClientService,
		),
	)
)
