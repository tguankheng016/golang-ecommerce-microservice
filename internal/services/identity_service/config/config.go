package config

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/caching"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/config"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/environment"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
)

type Config struct {
	ServerOptions   *http.ServerOptions       `mapstructure:"serverOptions"`
	PostgresOptions *postgres.PostgresOptions `mapstructure:"postgresOptions"`
	AuthOptions     *jwt.AuthOptions          `mapstructure:"authOptions"`
	RedisOptions    *caching.RedisOptions     `mapstructure:"redisOptions"`
	GrpcOptions     *grpc.GrpcOptions         `mapstructure:"grpcOptions"`
}

func InitConfig(env environment.Environment) (
	*Config,
	*http.ServerOptions,
	*postgres.PostgresOptions,
	*jwt.AuthOptions,
	*caching.RedisOptions,
	*grpc.GrpcOptions,
	error,
) {
	config, err := config.BindConfig[*Config](env)
	if err != nil {
		return returnError(err)
	}

	return config, config.ServerOptions, config.PostgresOptions, config.AuthOptions, config.RedisOptions, config.GrpcOptions, nil
}

func returnError(err error) (
	*Config,
	*http.ServerOptions,
	*postgres.PostgresOptions,
	*jwt.AuthOptions,
	*caching.RedisOptions,
	*grpc.GrpcOptions,
	error,
) {
	return nil, nil, nil, nil, nil, nil, err
}
