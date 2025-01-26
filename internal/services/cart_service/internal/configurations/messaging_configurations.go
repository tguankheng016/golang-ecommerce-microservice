package configurations

import (
	"github.com/ThreeDotsLabs/watermill/message"
	updating_product_out_of_stock "github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/carts/consumers/updating_product_out_of_stock"
	creating_product "github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/products/consumers/creating_product"
	deleting_product "github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/products/consumers/deleting_product"
	updating_product "github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/products/consumers/updating_product"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func ConfigureMessageHandler(
	router *message.Router,
	subscriber message.Subscriber,
	database *mongo.Database,
) {
	creating_product.MapHandler(router, subscriber, database)
	updating_product.MapHandler(router, subscriber, database)
	deleting_product.MapHandler(router, subscriber, database)

	updating_product_out_of_stock.MapHandler(router, subscriber, database)
}
