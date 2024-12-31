package seeds

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	userGrpc "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/grpc_client/protos"
	"go.uber.org/fx"
)

var (
	// Module provided to fx
	Module = fx.Module(
		"seeds_fx",
		seedInvoke,
	)

	seedInvoke = fx.Options(
		fx.Invoke(func(ctx context.Context, pool *pgxpool.Pool, userGrpcClientService userGrpc.UserGrpcServiceClient) error {
			return SeedData(ctx, pool, userGrpcClientService)
		}),
	)
)
