package config

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/caching"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/config"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/environment"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/messaging"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/mongo"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
)

type Config struct {
	ServerOptions   *http.ServerOptions         `mapstructure:"serverOptions"`
	MongoDbOptions  *mongo.MongoDbOptions       `mapstructure:"mongoDbOptions"`
	AuthOptions     *jwt.AuthOptions            `mapstructure:"authOptions"`
	RedisOptions    *caching.RedisOptions       `mapstructure:"redisOptions"`
	GrpcOptions     *grpc.GrpcOptions           `mapstructure:"grpcOptions"`
	WatermillOptons *messaging.WatermillOptions `mapstructure:"watermillOptions"`
}

func InitConfig(env environment.Environment) (
	*Config,
	*http.ServerOptions,
	*mongo.MongoDbOptions,
	*jwt.AuthOptions,
	*caching.RedisOptions,
	*grpc.GrpcOptions,
	*messaging.WatermillOptions,
	error,
) {
	config, err := config.BindConfig[*Config](env)
	if err != nil {
		return returnError(err)
	}

	return config, config.ServerOptions, config.MongoDbOptions, config.AuthOptions, config.RedisOptions, config.GrpcOptions, config.WatermillOptons, nil
}

func returnError(err error) (
	*Config,
	*http.ServerOptions,
	*mongo.MongoDbOptions,
	*jwt.AuthOptions,
	*caching.RedisOptions,
	*grpc.GrpcOptions,
	*messaging.WatermillOptions,
	error,
) {
	return nil, nil, nil, nil, nil, nil, nil, err
}
