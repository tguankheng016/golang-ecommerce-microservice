package config

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/caching"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/config"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/environment"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/messaging"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/openiddict"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/otel"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
)

type Config struct {
	ServerOptions   *http.ServerOptions         `mapstructure:"serverOptions"`
	PostgresOptions *postgres.PostgresOptions   `mapstructure:"postgresOptions"`
	AuthOptions     *jwt.AuthOptions            `mapstructure:"authOptions"`
	RedisOptions    *caching.RedisOptions       `mapstructure:"redisOptions"`
	GrpcOptions     *grpc.GrpcOptions           `mapstructure:"grpcOptions"`
	WatermillOptons *messaging.WatermillOptions `mapstructure:"watermillOptions"`
	JaegerOptions   *otel.JaegerOptions         `mapstructure:"jaegerOptions"`
	OAuthOptions    *openiddict.OAuthOptions    `mapstructure:"oauthOptions"`
}

func InitConfig(env environment.Environment) (
	*Config,
	*http.ServerOptions,
	*postgres.PostgresOptions,
	*jwt.AuthOptions,
	*caching.RedisOptions,
	*grpc.GrpcOptions,
	*messaging.WatermillOptions,
	*otel.JaegerOptions,
	*openiddict.OAuthOptions,
	error,
) {
	config, err := config.BindConfig[*Config](env)
	if err != nil {
		return returnError(err)
	}

	return config, config.ServerOptions, config.PostgresOptions, config.AuthOptions, config.RedisOptions, config.GrpcOptions, config.WatermillOptons, config.JaegerOptions, config.OAuthOptions, nil
}

func returnError(err error) (
	*Config,
	*http.ServerOptions,
	*postgres.PostgresOptions,
	*jwt.AuthOptions,
	*caching.RedisOptions,
	*grpc.GrpcOptions,
	*messaging.WatermillOptions,
	*otel.JaegerOptions,
	*openiddict.OAuthOptions,
	error,
) {
	return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
}
