package configurations

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jackc/pgx/v5/pgxpool"
	creating_user "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/consumers/creating_user"
	deleting_user "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/consumers/deleting_user"
	updating_user "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/users/consumers/updating_user"
)

func ConfigureMessageHandler(
	router *message.Router,
	subscriber message.Subscriber,
	pool *pgxpool.Pool,
) {
	creating_user.MapHandler(router, subscriber, pool)
	updating_user.MapHandler(router, subscriber, pool)
	deleting_user.MapHandler(router, subscriber, pool)
}
