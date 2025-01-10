package updatinguser

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/events"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logging"
	userService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/services"
	wotel "github.com/voi-oss/watermill-opentelemetry/pkg/opentelemetry"
)

func MapHandler(router *message.Router, subscriber message.Subscriber, pool *pgxpool.Pool) {
	router.AddNoPublisherHandler(
		"product_updating_user_v1",
		events.UserUpdatedTopicV1,
		subscriber,
		wotel.TraceNoPublishHandler(updateUser(pool)),
	)
}

func updateUser(pool *pgxpool.Pool) func(msg *message.Message) error {
	return func(msg *message.Message) error {
		event := events.UserUpdatedEvent{}
		if err := json.Unmarshal(msg.Payload, &event); err != nil {
			return err
		}

		logging.Logger.Info("updating user")

		userManager := userService.NewUserManager(pool)

		user, err := userManager.GetUserById(msg.Context(), event.Id)
		if err != nil {
			return err
		}

		if user == nil {
			if err := copier.Copy(&user, &event); err != nil {
				return err
			}
			if err := userManager.CreateUser(msg.Context(), user); err != nil {
				return err
			}
		} else {
			if err := copier.Copy(&user, &event); err != nil {
				return err
			}
			if err := userManager.UpdateUser(msg.Context(), user); err != nil {
				return err
			}
		}

		return nil
	}
}
