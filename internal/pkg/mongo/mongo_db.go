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

func NewMongoDb(cfg *MongoDbOptions) (*mongo.Client, *mongo.Database, error) {
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
		return nil, nil, err
	}

	database := client.Database(cfg.Database)

	return client, database, nil
}

func RunMongoDB(lc fx.Lifecycle, logger *zap.Logger, client *mongo.Client) error {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("starting mongo ...")

			if err := client.Ping(ctx, readpref.Primary()); err != nil {
				return err
			}

			logger.Info("mongo connected ...")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("disconnecting mongo...")

			if err := client.Disconnect(ctx); err != nil {
				logger.Info("error when closing mongo...", zap.Error(err))
			}

			logger.Info("mongo connection disconnected")

			return nil
		},
	})

	return nil
}
