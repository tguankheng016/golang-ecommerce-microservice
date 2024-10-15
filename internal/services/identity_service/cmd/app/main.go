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
	identityService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/identities/services"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/server"
	userService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/services"
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
				redis.NewRedisUniversalClient,
				echoServer.NewEchoServer,
				jwt.NewTokenHandler,
				jwt.NewTokenKeyValidator,
				jwt.NewSecurityStampValidator,
				identityService.NewCustomStampDBValidator,
				identityService.NewCustomTokenKeyDBValidator,
				identityService.NewJwtTokenGenerator,
				userService.NewUserRolePermissionManager,
				userService.NewUserPermissionDbManager,
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
