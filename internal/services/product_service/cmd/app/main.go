package main

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/caching"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/environment"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logging"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/messaging"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/config"

	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/configurations"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/data/seeds"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Options(
			fx.Provide(
				environment.ConfigureAppEnv,
				config.InitConfig,
			),
			logging.Module,
			postgres.Module,
			caching.Module,
			security.Module,
			security.DefaultModule,
			permissions.Module,
			permissions.DefaultModule,
			users.Module,
			seeds.Module,
			configurations.Module,
			http.Module,
			grpc.Module,
			messaging.Module,
		),
	).Run()
}
