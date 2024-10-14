package redis

import (
	"context"
	"fmt"

	redisCli "github.com/redis/go-redis/v9"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewRedisClient(options *RedisOptions) *redisCli.Client {
	client := redisCli.NewClient(&redisCli.Options{
		Addr:     fmt.Sprintf("%s:%d", options.Host, options.Port),
		Password: options.Password,
		DB:       options.Database,
		PoolSize: options.PoolSize,
	})

	return client
}

func RegisterRedisServer(lc fx.Lifecycle, client *redisCli.Client, ctx context.Context) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Ping(ctx).Err()
		},
		OnStop: func(ctx context.Context) error {
			if err := client.Close(); err != nil {
				logger.Logger.Error("error in closing redis", zap.Error(err))
			} else {
				logger.Logger.Info("redis closed gracefully")
			}

			return nil
		},
	})

	return nil
}
