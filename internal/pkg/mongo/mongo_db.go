package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	connectTimeout  = 60 * time.Second
	maxConnIdleTime = 3 * time.Minute
	minPoolSize     = 20
	maxPoolSize     = 300
)

func NewMongoDb(ctx context.Context, cfg *MongoDbOptions) (*mongo.Client, error) {
	uriAddress := fmt.Sprintf(
		"mongodb://%s:%d/?ssl=false&directConnection=true",
		cfg.Host,
		cfg.Port,
	)

	opt := options.Client().ApplyURI(uriAddress).
		SetConnectTimeout(connectTimeout).
		SetMaxConnIdleTime(maxConnIdleTime).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize)

	client, err := mongo.Connect(opt)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func RunMongoDB(lc fx.Lifecycle, logger *zap.Logger, client *mongo.Client, ctx context.Context) error {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			logger.Info("starting mongo ...")
			if err := client.Ping(ctx, readpref.Primary()); err != nil {
				return err
			}
			logger.Info("mongo connected ...")

			return nil
		},
		OnStop: func(_ context.Context) error {
			logger.Info("closing mongo...")
			if err := client.Disconnect(ctx); err != nil {
				logger.Info("error when closing mongo...", zap.Error(err))
			}

			return nil
		},
	})

	return nil
}
