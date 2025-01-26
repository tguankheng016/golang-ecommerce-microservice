package config

import (
	"os"

	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/azure"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/caching"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/config"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/environment"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/grpc"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/messaging"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/otel"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
)

type GrpcAddress struct {
	IdentityAddress string `mapstructure:"identityAddress"`
}

type Config struct {
	ServerOptions   *http.ServerOptions         `mapstructure:"serverOptions"`
	PostgresOptions *postgres.PostgresOptions   `mapstructure:"postgresOptions"`
	AuthOptions     *jwt.AuthOptions            `mapstructure:"authOptions"`
	RedisOptions    *caching.RedisOptions       `mapstructure:"redisOptions"`
	GrpcOptions     *grpc.GrpcOptions           `mapstructure:"grpcOptions"`
	GrpcAddresses   *GrpcAddress                `mapstructure:"grpcAddresses"`
	WatermillOptons *messaging.WatermillOptions `mapstructure:"watermillOptions"`
	JaegerOptions   *otel.JaegerOptions         `mapstructure:"jaegerOptions"`
	AzureOptions    *azure.AzureOptions         `mapstructure:"azureOptions"`
}

func InitConfig(env environment.Environment) (
	*Config,
	*http.ServerOptions,
	*postgres.PostgresOptions,
	*jwt.AuthOptions,
	*caching.RedisOptions,
	*grpc.GrpcOptions,
	*GrpcAddress,
	*messaging.WatermillOptions,
	*otel.JaegerOptions,
	error,
) {
	config, err := config.BindConfig[*Config](env)
	if err != nil {
		return returnError(err)
	}

	if err := config.loadAzureConfig(); err != nil {
		return returnError(err)
	}

	return config, config.ServerOptions, config.PostgresOptions, config.AuthOptions, config.RedisOptions, config.GrpcOptions, config.GrpcAddresses, config.WatermillOptons, config.JaegerOptions, nil
}

func (config *Config) loadAzureConfig() error {
	if !config.AzureOptions.Enabled {
		return nil
	}

	keyVaultName := os.Getenv("GoProductKeyVault")
	if keyVaultName != "" {
		config.AzureOptions.KeyVaultName = keyVaultName
	}

	azureClient, err := azure.NewAzureClient(config.AzureOptions)
	if err != nil {
		return err
	}

	config.PostgresOptions.Datasource, err = azureClient.GetSecret("postgresOptions--datasource")
	if err != nil {
		return err
	}
	config.AuthOptions.SecretKey, err = azureClient.GetSecret("authOptions--secretKey")
	if err != nil {
		return err
	}

	return nil
}

func returnError(err error) (
	*Config,
	*http.ServerOptions,
	*postgres.PostgresOptions,
	*jwt.AuthOptions,
	*caching.RedisOptions,
	*grpc.GrpcOptions,
	*GrpcAddress,
	*messaging.WatermillOptions,
	*otel.JaegerOptions,
	error,
) {
	return nil, nil, nil, nil, nil, nil, nil, nil, nil, err
}
