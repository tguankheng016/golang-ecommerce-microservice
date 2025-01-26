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
	productModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/products/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Request
type UpdateCartDto struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type HumaUpdateCartRequest struct {
	Body struct {
		UpdateCartDto
	}
}

// Validator
func (e HumaUpdateCartRequest) Schema() v.Schema {
	return v.Schema{
		v.F("id", e.Body.Id):             v.Nonzero[string]().Msg("invalid cart id"),
		v.F("quantity", e.Body.Quantity): v.Gte(0).Msg("invalid quantity"),
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
			OperationID:   "UpdateCart",
			Method:        http.MethodPut,
			Path:          "/carts/cart",
			Summary:       "Update Cart",
			Tags:          []string{"Carts"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, ""),
			},
		},
		updateCart(database, publisher),
	)
}

func updateCart(database *mongo.Database, publisher message.Publisher) func(context.Context, *HumaUpdateCartRequest) (*struct{}, error) {
	return func(ctx context.Context, request *HumaUpdateCartRequest) (*struct{}, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		userId, ok := httpServer.GetCurrentUser(ctx)
		if !ok {
			return nil, huma.Error401Unauthorized("unable to get current user id")
		}

		productCollection := database.Collection(productModel.ProductCollectionName)
		cartCollection := database.Collection(cartModel.CartCollectionName)

		filter := bson.D{
			bson.E{Key: "id", Value: request.Body.Id},
			bson.E{Key: "user_id", Value: userId},
		}

		quantityChanged := 0

		var cart cartModel.Cart
		err := cartCollection.FindOne(ctx, filter).Decode(&cart)

		if err != nil && err != mongo.ErrNoDocuments {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if err == mongo.ErrNoDocuments {
			return nil, huma.Error400BadRequest("invalid cart id")
		}

		productFilter := bson.D{bson.E{Key: "id", Value: cart.ProductId}}

		var product productModel.Product
		err = productCollection.FindOne(ctx, productFilter).Decode(&product)

		if err != nil && err != mongo.ErrNoDocuments {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if err == mongo.ErrNoDocuments {
			return nil, huma.Error400BadRequest("invalid product id")
		}

		if request.Body.Quantity == 0 {
			quantityChanged = cart.Quantity

			_, err := cartCollection.DeleteOne(ctx, filter)
			if err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}
		} else if cart.Quantity != request.Body.Quantity {
			quantityChanged = cart.Quantity - request.Body.Quantity

			if quantityChanged < 0 {
				if product.StockQuantity == 0 {
					return nil, huma.Error400BadRequest("this product is out of stock.")
				}

				if (-1 * quantityChanged) > product.StockQuantity {
					return nil, huma.Error400BadRequest("the selected quantity exceeds the remaining stock.")
				}
			}

			update := bson.D{
				{
					Key: "$set",
					Value: bson.D{
						bson.E{Key: "quantity", Value: request.Body.Quantity},
					},
				},
			}

			_, err = cartCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}
		}

		if quantityChanged != 0 {
			changeProductQuantityEvent := events.ChangeProductQuantityEvent{
				Id:              cart.ProductId,
				QuantityChanged: quantityChanged,
			}

			payload, err := json.Marshal(changeProductQuantityEvent)
			if err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}

			msg := message.NewMessage(watermill.NewUUID(), payload)
			publisher.Publish(events.ChangeProductQuantityTopicV1, msg)
		}

		return nil, nil
	}
}
