package deletinguser

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/events"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logging"
	userService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/services"
)

func MapHandler(router *message.Router, subscriber message.Subscriber, pool *pgxpool.Pool) {
	router.AddNoPublisherHandler(
		"product_deleting_user_v1",
		events.UserDeletedTopicV1,
		subscriber,
		deleteUser(pool),
	)
}

func deleteUser(pool *pgxpool.Pool) func(msg *message.Message) error {
	return func(msg *message.Message) error {
		event := events.UserDeletedEvent{}
		if err := json.Unmarshal(msg.Payload, &event); err != nil {
			return err
		}

		logging.Logger.Info("deleting user")

		userManager := userService.NewUserManager(pool)

		if err := userManager.DeleteUser(msg.Context(), event.Id); err != nil {
			return err
		}

		return nil
	}
}
