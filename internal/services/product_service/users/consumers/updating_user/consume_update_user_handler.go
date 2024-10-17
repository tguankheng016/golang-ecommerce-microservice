package updating_user

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/streadway/amqp"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/events"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/rabbitmq"
	userModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/users/models"
)

func HandleConsumeUpdateUser(queue string, msg amqp.Delivery, deliveryBase *rabbitmq.RabbitMQDeliveryBase) error {
	logger.Logger.Info(fmt.Sprintf("Message received on queue: %s with message: %s", queue, string(msg.Body)))

	var userUpdatedEvent events.UserUpdatedEvent
	err := json.Unmarshal(msg.Body, &userUpdatedEvent)
	if err != nil {
		return err
	}

	var count int64
	if err := deliveryBase.Gorm.Model(&userModel.User{}).Where("id = ?", userUpdatedEvent.Id).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		// New User
		var newUser userModel.User
		if err := copier.Copy(&newUser, &userUpdatedEvent); err != nil {
			return err
		}
		if err := deliveryBase.Gorm.Create(&newUser).Error; err != nil {
			return err
		}
	} else {
		// Update User
		var userToUpdate userModel.User
		if err := deliveryBase.Gorm.First(&userToUpdate, userUpdatedEvent.Id).Error; err != nil {
			return err
		}

		if err := copier.Copy(&userToUpdate, &userUpdatedEvent); err != nil {
			return err
		}

		if err := deliveryBase.Gorm.Save(&userToUpdate).Error; err != nil {
			return err
		}
	}

	return nil
}
