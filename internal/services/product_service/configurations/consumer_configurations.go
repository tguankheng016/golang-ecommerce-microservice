package configurations

import (
	"context"

	"github.com/streadway/amqp"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/events"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/rabbitmq"
	creating_user "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/users/consumers/creating_user"
	deleting_user "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/users/consumers/deleting_user"
	updating_user "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/users/consumers/updating_user"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func ConfigConsumers(
	ctx context.Context,
	gorm *gorm.DB,
	rabbitMQOptions *rabbitmq.RabbitMQOptions,
	connRabbitmq *amqp.Connection,
	rabbitmqPublisher rabbitmq.IPublisher,
) error {
	deliveryBase := rabbitmq.RabbitMQDeliveryBase{
		ConnRabbitmq:      connRabbitmq,
		RabbitmqPublisher: rabbitmqPublisher,
		Ctx:               ctx,
		Gorm:              gorm,
	}

	createUserConsumer := rabbitmq.NewConsumer(ctx, rabbitMQOptions, connRabbitmq, creating_user.HandleConsumeCreateUser)
	updateUserConsumer := rabbitmq.NewConsumer(ctx, rabbitMQOptions, connRabbitmq, updating_user.HandleConsumeUpdateUser)
	deleteUserConsumer := rabbitmq.NewConsumer(ctx, rabbitMQOptions, connRabbitmq, deleting_user.HandleConsumeDeleteUser)

	go func() {
		if err := createUserConsumer.ConsumeMessage(events.UserCreatedEvent{}, &deliveryBase); err != nil {
			logger.Logger.Error("Error when consuming create user event:", zap.Error(err))
		}

		if err := updateUserConsumer.ConsumeMessage(events.UserUpdatedEvent{}, &deliveryBase); err != nil {
			logger.Logger.Error("Error when consuming update user event:", zap.Error(err))
		}

		if err := deleteUserConsumer.ConsumeMessage(events.UserDeletedEvent{}, &deliveryBase); err != nil {
			logger.Logger.Error("Error when consuming delete user event:", zap.Error(err))
		}
	}()

	return nil
}
