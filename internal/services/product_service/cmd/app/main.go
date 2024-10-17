package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/config/environment"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	gormDb "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres_gorm"
	redis "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/redis"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/config"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/configurations"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/data/seeds"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/users"
	"go.uber.org/fx"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	fx.New(
		fx.Options(
			fx.Provide(
				environment.ConfigAppEnv,
				config.InitConfig,
				validator.New,
			),
			logger.Module,
			http.Module,
			grpc.Module,
			gormDb.Module,
			redis.Module,
			security.Module,
			security.DefaultModule,
			users.Module,
			permissions.Module,
			permissions.DefaultModule,
			seeds.Module,
			configurations.Module,
		),
	).Run()
}
