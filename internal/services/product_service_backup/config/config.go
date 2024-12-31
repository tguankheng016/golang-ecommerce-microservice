package config

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/config"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/config/environment"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	echoserver "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http/echo"
	gormdb "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres_gorm"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/rabbitmq"
	redis "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/redis"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
)

type GrpcAddress struct {
	IdentityAddress string `mapstructure:"identityAddress"`
}

type Config struct {
	GormOptions     *gormdb.GormOptions       `mapstructure:"gormOptions"`
	EchoOptions     *echoserver.EchoOptions   `mapstructure:"echoOptions"`
	AuthOptions     *jwt.AuthOptions          `mapstructure:"authOptions"`
	RedisOptions    *redis.RedisOptions       `mapstructure:"redisOptions"`
	GrpcOptions     *grpc.GrpcOptions         `mapstructure:"grpcOptions"`
	GrpcAddresses   *GrpcAddress              `mapstructure:"grpcAddresses"`
	RabbitMQOptions *rabbitmq.RabbitMQOptions `mapstructure:"rabbitMQOptions"`
}

func InitConfig(env environment.Environment) (*Config, *gormdb.GormOptions,
	*echoserver.EchoOptions, *jwt.AuthOptions, *redis.RedisOptions, *grpc.GrpcOptions,
	*GrpcAddress, *rabbitmq.RabbitMQOptions, error) {

	cfg, err := config.BindConfig[*Config](env)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	return cfg, cfg.GormOptions, cfg.EchoOptions, cfg.AuthOptions, cfg.RedisOptions, cfg.GrpcOptions, cfg.GrpcAddresses, cfg.RabbitMQOptions, nil
}
