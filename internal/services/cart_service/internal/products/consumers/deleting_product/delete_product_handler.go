package deletingproduct

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/events"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logging"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/products/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func MapHandler(router *message.Router, subscriber message.Subscriber, database *mongo.Database) {
	router.AddNoPublisherHandler(
		"product_deleting_product_v1",
		events.ProductDeletedTopicV1,
		subscriber,
		deleteProduct(database),
	)
}

func deleteProduct(database *mongo.Database) func(msg *message.Message) error {
	return func(msg *message.Message) error {
		event := events.ProductDeletedEvent{}
		if err := json.Unmarshal(msg.Payload, &event); err != nil {
			return err
		}

		logging.Logger.Info("deleting product")

		productCollection := database.Collection(models.ProductCollectionName)
		filter := bson.D{bson.E{Key: "id", Value: event.Id}}

		_, err := productCollection.DeleteOne(msg.Context(), filter)
		if err != nil {
			return err
		}

		return nil
	}
}
