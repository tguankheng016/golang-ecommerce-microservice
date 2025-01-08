package updatingproductoutofstock

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/events"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logging"
	cartModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/carts/models"
	productModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/products/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func MapHandler(router *message.Router, subscriber message.Subscriber, database *mongo.Database) {
	router.AddNoPublisherHandler(
		"cart_updating_product_out_of_stock_v1",
		events.ProductOutOfStockTopicV1,
		subscriber,
		updateProductOutOfStock(database),
	)
}

func updateProductOutOfStock(database *mongo.Database) func(msg *message.Message) error {
	return func(msg *message.Message) error {
		event := events.ProductOutOfStockEvent{}
		if err := json.Unmarshal(msg.Payload, &event); err != nil {
			return err
		}

		logging.Logger.Info("updating products out of stock status")

		cartCollection := database.Collection(cartModel.CartCollectionName)
		productCollection := database.Collection(productModel.ProductCollectionName)

		filter := bson.D{bson.E{Key: "product_id", Value: event.Id}}

		update := bson.D{
			{
				Key: "$set",
				Value: bson.D{
					bson.E{Key: "is_out_of_stock", Value: event.StockQuantity <= 0},
				},
			},
		}

		if _, err := cartCollection.UpdateMany(msg.Context(), filter, update); err != nil {
			return err
		}

		filter = bson.D{bson.E{Key: "id", Value: event.Id}}

		update = bson.D{
			{
				Key: "$set",
				Value: bson.D{
					bson.E{Key: "stock_quantity", Value: event.StockQuantity},
				},
			},
		}

		if _, err := productCollection.UpdateOne(msg.Context(), filter, update); err != nil {
			return err
		}

		return nil
	}
}
