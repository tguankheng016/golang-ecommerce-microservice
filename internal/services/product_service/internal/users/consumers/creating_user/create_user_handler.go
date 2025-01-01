package creatinguser

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/events"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logging"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/models"
	userService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/services"
)

func MapHandler(router *message.Router, subscriber message.Subscriber, pool *pgxpool.Pool) {
	router.AddNoPublisherHandler(
		"product_creating_user_v1",
		events.UserCreatedTopicV1,
		subscriber,
		createUser(pool),
	)
}

func createUser(pool *pgxpool.Pool) func(msg *message.Message) error {
	return func(msg *message.Message) error {
		event := events.UserCreatedEvent{}
		if err := json.Unmarshal(msg.Payload, &event); err != nil {
			return err
		}

		logging.Logger.Info("creating user")

		userManager := userService.NewUserManager(pool)

		var user models.User
		if err := copier.Copy(&user, &event); err != nil {
			return err
		}

		if err := userManager.CreateUser(msg.Context(), &user); err != nil {
			return err
		}

		return nil
	}
}
