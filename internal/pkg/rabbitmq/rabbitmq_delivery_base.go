package rabbitmq

import (
	"context"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type RabbitMQDeliveryBase struct {
	RabbitmqPublisher IPublisher
	ConnRabbitmq      *amqp.Connection
	Gorm              *gorm.DB
	Ctx               context.Context
}
