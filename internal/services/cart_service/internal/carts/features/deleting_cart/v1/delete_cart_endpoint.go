package v1

import (
	"context"
	"encoding/json"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/danielgtaylor/huma/v2"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/events"
	httpServer "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	cartModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/carts/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Request
type DeleteCartRequest struct {
	Id string `path:"id"`
}

// Validator
func (e DeleteCartRequest) Schema() v.Schema {
	return v.Schema{
		v.F("id", e.Id): v.Nonzero[string]().Msg("Invalid cart id"),
	}
}

// Handler
func MapRoute(
	api huma.API,
	database *mongo.Database,
	publisher message.Publisher,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "DeleteCart",
			Method:        http.MethodDelete,
			Path:          "/carts/cart/{id}",
			Summary:       "Delete Cart",
			Tags:          []string{"Carts"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, ""),
			},
		},
		addCart(database, publisher),
	)
}

func addCart(database *mongo.Database, publisher message.Publisher) func(context.Context, *DeleteCartRequest) (*struct{}, error) {
	return func(ctx context.Context, request *DeleteCartRequest) (*struct{}, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		userId, ok := httpServer.GetCurrentUser(ctx)
		if !ok {
			return nil, huma.Error401Unauthorized("unable to get current user id")
		}

		cartCollection := database.Collection(cartModel.CartCollectionName)

		filter := bson.D{
			bson.E{Key: "id", Value: request.Id},
			bson.E{Key: "user_id", Value: userId},
		}

		var cart cartModel.Cart
		err := cartCollection.FindOne(ctx, filter).Decode(&cart)

		if err != nil && err != mongo.ErrNoDocuments {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if err == mongo.ErrNoDocuments {
			return nil, huma.Error400BadRequest("invalid cart id")
		}

		_, err = cartCollection.DeleteOne(ctx, filter)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		changeProductQuantityEvent := events.ChangeProductQuantityEvent{
			Id:              cart.ProductId,
			QuantityChanged: cart.Quantity,
		}

		payload, err := json.Marshal(changeProductQuantityEvent)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		msg := message.NewMessage(watermill.NewUUID(), payload)
		publisher.Publish(events.ChangeProductQuantityTopicV1, msg)

		return nil, nil
	}
}
