package users

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/grpc_client/services"
	"go.uber.org/fx"
)

var (
	// Module provided to fx
	Module = fx.Module(
		"user_fx",
		userProviders,
	)

	userProviders = fx.Options(
		fx.Provide(
			services.NewUserGrpcClientService,
		),
	)
)
