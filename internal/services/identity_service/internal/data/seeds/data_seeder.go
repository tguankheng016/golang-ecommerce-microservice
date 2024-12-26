package seeds

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	roleSeeder "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/roles/seed"
	userSeeder "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/seed"
)

func SeedData(ctx context.Context, pool *pgxpool.Pool) error {
	if err := roleSeeder.NewRoleSeeder(pool).SeedRoles(ctx); err != nil {
		return err
	}

	if err := userSeeder.NewUserSeeder(pool).SeedUsers(ctx); err != nil {
		return err
	}

	return nil
}
