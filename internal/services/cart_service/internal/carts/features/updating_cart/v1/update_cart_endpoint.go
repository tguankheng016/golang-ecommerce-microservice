package v1

import (
	"context"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	httpServer "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	cartModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/carts/models"
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
		updateCart(database),
	)
}

func updateCart(database *mongo.Database) func(context.Context, *HumaUpdateCartRequest) (*struct{}, error) {
	return func(ctx context.Context, request *HumaUpdateCartRequest) (*struct{}, error) {
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
			bson.E{Key: "id", Value: request.Body.Id},
			bson.E{Key: "userid", Value: userId},
		}

		err := cartCollection.FindOne(ctx, filter).Err()

		if err != nil && err != mongo.ErrNoDocuments {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if err == mongo.ErrNoDocuments {
			return nil, huma.Error400BadRequest("invalid cart id")
		}

		if request.Body.Quantity == 0 {
			_, err := cartCollection.DeleteOne(ctx, filter)
			if err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}
		} else {
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

		return nil, nil
	}
}
