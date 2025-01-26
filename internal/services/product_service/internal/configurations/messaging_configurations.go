package configurations

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jackc/pgx/v5/pgxpool"
	changing_product_quantity "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/consumers/changing_product_quantity"
	creating_user "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/consumers/creating_user"
	deleting_user "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/consumers/deleting_user"
	updating_user "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/consumers/updating_user"
)

func ConfigureMessageHandler(
	router *message.Router,
	subscriber message.Subscriber,
	publisher message.Publisher,
	pool *pgxpool.Pool,
) {
	creating_user.MapHandler(router, subscriber, pool)
	updating_user.MapHandler(router, subscriber, pool)
	deleting_user.MapHandler(router, subscriber, pool)

	changing_product_quantity.MapHandler(router, subscriber, pool, publisher)
}
