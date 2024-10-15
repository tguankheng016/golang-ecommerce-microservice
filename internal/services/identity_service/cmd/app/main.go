package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/config/environment"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http"
	echoServer "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http/echo"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
	gormDb "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres_gorm"
	redis "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/redis"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/config"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/configurations"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/data/seeds"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/identities/services"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/server"
	"go.uber.org/fx"
	"gorm.io/gorm"
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
				logger.InitLogger,
				http.NewContext,
				gormDb.NewGormDB,
				redis.NewRedisClient,
				echoServer.NewEchoServer,
				jwt.NewTokenHandler,
				jwt.NewTokenKeyValidator,
				jwt.NewSecurityStampValidator,
				services.NewCustomStampDBValidator,
				services.NewCustomTokenKeyDBValidator,
				services.NewJwtTokenGenerator,
				validator.New,
			),
			fx.Invoke(logger.RunLogger),
			fx.Invoke(server.RunServers),
			fx.Invoke(redis.RegisterRedisServer),
			fx.Invoke(configurations.ConfigMiddlewares),
			fx.Invoke(configurations.ConfigSwagger),
			fx.Invoke(func(db *gorm.DB) error {
				if err := gormDb.RunGooseMigration(db); err != nil {
					return err
				}
				return seeds.DataSeeder(db)
			}),
			fx.Invoke(configurations.ConfigEndpoints),
		),
	).Run()
}
