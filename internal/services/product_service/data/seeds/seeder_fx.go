package seeds

import (
	user_service "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/users/grpc_client/protos"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var (
	// Module provided to fx
	Module = fx.Module(
		"productseederfx",
		seedInvoke,
	)

	seedInvoke = fx.Options(
		fx.Invoke(func(db *gorm.DB, userClientGrpcService user_service.UserGrpcServiceClient) error {
			return DataSeeder(db, userClientGrpcService)
		}),
	)
)
