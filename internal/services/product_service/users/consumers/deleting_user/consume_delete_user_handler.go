package deleting_user

import (
	"encoding/json"
	"fmt"

	"github.com/anchore/go-logger"
	"github.com/streadway/amqp"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/events"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/rabbitmq"
	userModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/users/models"
)

func HandleConsumeDeleteUser(queue string, msg amqp.Delivery, deliveryBase *rabbitmq.RabbitMQDeliveryBase) error {
	logger.Logger.Info(fmt.Sprintf("Message received on queue: %s with message: %s", queue, string(msg.Body)))

	var userDeletedEvent events.UserDeletedEvent
	err := json.Unmarshal(msg.Body, &userDeletedEvent)
	if err != nil {
		return err
	}

	var count int64
	if err := deliveryBase.Gorm.Model(&userModel.User{}).Where("id = ?", userDeletedEvent.Id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		// User Exist
		var userToDelete userModel.User
		if err := deliveryBase.Gorm.First(&userToDelete, userDeletedEvent.Id).Error; err != nil {
			return err
		}
	}

	return nil
}
