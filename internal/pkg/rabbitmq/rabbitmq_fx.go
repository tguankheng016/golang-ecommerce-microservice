package rabbitmq

import (
	"github.com/streadway/amqp"
	"go.uber.org/fx"
)

var (
	// Module provided to fx
	Module = fx.Module(
		"rabbitmqfx",
		rabbitMQProviders,
		rabbitMQInvokes,
	)

	rabbitMQProviders = fx.Options(
		fx.Provide(
			NewRabbitMQConn,
			NewPublisher,
		),
	)

	rabbitMQInvokes = fx.Options(
		fx.Invoke(func(conn *amqp.Connection) error {
			return nil
		}),
	)
)
