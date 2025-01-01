package seeds

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	categorySeeder "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/seed"
	productSeeder "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/seed"
	userGrpc "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/grpc_client/protos"
	userSeeder "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/seed"
)

func SeedData(ctx context.Context, pool *pgxpool.Pool, userGrpcClientService userGrpc.UserGrpcServiceClient) error {
	if err := userSeeder.NewUserSeeder(pool, userGrpcClientService).SeedUsers(ctx); err != nil {
		return err
	}

	if err := categorySeeder.NewCategorySeeder(pool).SeedCategories(ctx); err != nil {
		return err
	}

	if err := productSeeder.NewProductSeeder(pool).SeedProducts(ctx); err != nil {
		return err
	}

	return nil
}
