package changingproductquantity

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/events"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logging"
	productService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/services"
	wotel "github.com/voi-oss/watermill-opentelemetry/pkg/opentelemetry"
)

func MapHandler(router *message.Router, subscriber message.Subscriber, pool *pgxpool.Pool, publisher message.Publisher) {
	router.AddHandler(
		"product_changing_product_quantity_v1",
		events.ChangeProductQuantityTopicV1,
		subscriber,
		events.ProductOutOfStockTopicV1,
		publisher,
		wotel.TraceHandler(updateProductQuantity(pool)),
	)
}

func updateProductQuantity(pool *pgxpool.Pool) func(msg *message.Message) ([]*message.Message, error) {
	return func(msg *message.Message) ([]*message.Message, error) {
		event := events.ChangeProductQuantityEvent{}
		if err := json.Unmarshal(msg.Payload, &event); err != nil {
			return nil, err
		}

		logging.Logger.Info("updating product quantity")

		productManager := productService.NewProductManager(pool)

		product, err := productManager.GetProductById(msg.Context(), event.Id)
		if err != nil {
			return nil, err
		}

		product.StockQuantity += event.QuantityChanged

		if product.StockQuantity < 0 {
			product.StockQuantity = 0
		}

		if err := productManager.UpdateProduct(msg.Context(), product); err != nil {
			return nil, err
		}

		callbackEvent := events.ProductOutOfStockEvent{
			Id:            product.Id,
			StockQuantity: product.StockQuantity,
		}

		newPayload, err := json.Marshal(callbackEvent)
		if err != nil {
			return nil, err
		}

		newMessage := message.NewMessage(watermill.NewUUID(), newPayload)

		return []*message.Message{newMessage}, nil
	}
}
