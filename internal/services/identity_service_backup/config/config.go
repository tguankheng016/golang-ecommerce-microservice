package config

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/config"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/config/environment"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	echoserver "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http/echo"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/openiddict"
	gormdb "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres_gorm"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/rabbitmq"
	redis "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/redis"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
)

type Config struct {
	AppOptions      *core.AppOptions          `mapstructure:"appOptions"`
	GormOptions     *gormdb.GormOptions       `mapstructure:"gormOptions"`
	EchoOptions     *echoserver.EchoOptions   `mapstructure:"echoOptions"`
	AuthOptions     *jwt.AuthOptions          `mapstructure:"authOptions"`
	RedisOptions    *redis.RedisOptions       `mapstructure:"redisOptions"`
	GrpcOptions     *grpc.GrpcOptions         `mapstructure:"grpcOptions"`
	RabbitMQOptions *rabbitmq.RabbitMQOptions `mapstructure:"rabbitMQOptions"`
	OAuthOptions    *openiddict.OAuthOptions  `mapstructure:"oauthOptions"`
}

func InitConfig(env environment.Environment) (*Config, *gormdb.GormOptions,
	*echoserver.EchoOptions, *jwt.AuthOptions, *redis.RedisOptions, *grpc.GrpcOptions,
	*rabbitmq.RabbitMQOptions, *openiddict.OAuthOptions, *core.AppOptions, error) {

	cfg, err := config.BindConfig[*Config](env)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	return cfg, cfg.GormOptions, cfg.EchoOptions, cfg.AuthOptions, cfg.RedisOptions, cfg.GrpcOptions, cfg.RabbitMQOptions, cfg.OAuthOptions, cfg.AppOptions, nil
}
