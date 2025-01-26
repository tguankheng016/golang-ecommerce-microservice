package updatingproduct

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/events"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logging"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/products/models"
	wotel "github.com/voi-oss/watermill-opentelemetry/pkg/opentelemetry"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func MapHandler(router *message.Router, subscriber message.Subscriber, database *mongo.Database) {
	router.AddNoPublisherHandler(
		"cart_updating_product_v1",
		events.ProductUpdatedTopicV1,
		subscriber,
		wotel.TraceNoPublishHandler(updateProduct(database)),
	)
}

func updateProduct(database *mongo.Database) func(msg *message.Message) error {
	return func(msg *message.Message) error {
		event := events.ProductUpdatedEvent{}
		if err := json.Unmarshal(msg.Payload, &event); err != nil {
			return err
		}

		logging.Logger.Info("updating product")

		productCollection := database.Collection(models.ProductCollectionName)

		filter := bson.D{bson.E{Key: "id", Value: event.Id}}

		var dbProduct models.Product
		err := productCollection.FindOne(msg.Context(), filter).Decode(&dbProduct)

		if err != nil && err != mongo.ErrNoDocuments {
			return err
		}

		if err == mongo.ErrNoDocuments {
			newProduct := models.Product{
				Id:            int(event.Id),
				Name:          event.Name,
				Description:   event.Description,
				Price:         event.Price,
				StockQuantity: event.StockQuantity,
			}
			_, err := productCollection.InsertOne(msg.Context(), newProduct)
			if err != nil {
				return err
			}
		} else {
			update := bson.D{
				{
					Key: "$set",
					Value: bson.D{
						bson.E{Key: "name", Value: event.Name},
						bson.E{Key: "description", Value: event.Description},
						bson.E{Key: "price", Value: event.Price},
						bson.E{Key: "stock_quantity", Value: event.StockQuantity},
					},
				},
			}

			_, err := productCollection.UpdateOne(msg.Context(), filter, update)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
