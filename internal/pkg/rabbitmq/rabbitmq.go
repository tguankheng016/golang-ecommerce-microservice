package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/streadway/amqp"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logging"
	"go.uber.org/zap"
)

func NewRabbitMQConn(cfg *RabbitMQOptions, ctx context.Context) (*amqp.Connection, error) {
	connAddr := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 10 * time.Second // Maximum time to retry
	maxRetries := 5                      // Number of retries (including the initial attempt)

	var conn *amqp.Connection
	var err error

	err = backoff.Retry(func() error {

		conn, err = amqp.Dial(connAddr)
		if err != nil {
			logging.Logger.Error(fmt.Sprintf("Failed to connect to RabbitMQ: %v. Connection information: %s", err, connAddr), zap.Error(err))
			return err
		}

		return nil
	}, backoff.WithMaxRetries(bo, uint64(maxRetries-1)))

	logging.Logger.Info("Connected to RabbitMQ")

	go func() {
		<-ctx.Done()
		err := conn.Close()
		if err != nil {
			logging.Logger.Error("Failed to close RabbitMQ connection", zap.Error(err))
		}
		logging.Logger.Info("RabbitMQ connection is closed")
	}()

	return conn, err
}
