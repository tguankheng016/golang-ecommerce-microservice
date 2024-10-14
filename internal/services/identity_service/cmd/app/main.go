package main

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/config/environment"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http"
	echoserver "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http/echo"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
	gormdb "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres_gorm"
	redis "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/redis"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/config"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/server"
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
				logger.InitLogger,
				http.NewContext,
				gormdb.NewGormDB,
				redis.NewRedisClient,
				echoserver.NewEchoServer,
			),
			fx.Invoke(logger.RunLogger),
			fx.Invoke(server.RunServers),
			fx.Invoke(redis.RegisterRedisServer),
		),
	).Run()
}
