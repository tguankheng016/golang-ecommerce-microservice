package products

import (
	"context"

	productGrpc "github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/products/grpc_client/protos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/products/grpc_client/services"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/products/seed"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"go.uber.org/fx"
)

var (
	// Module provided to fx
	Module = fx.Module(
		"product_fx",
		productProviders,
		productInvokes,
	)

	productProviders = fx.Options(
		fx.Provide(
			services.NewProductGrpcClientService,
		),
	)

	productInvokes = fx.Options(
		fx.Invoke(func(ctx context.Context, productGrpcClientService productGrpc.ProductGrpcServiceClient, database *mongo.Database) error {
			return seed.SeedProducts(ctx, productGrpcClientService, database)
		}),
	)
)
