package v1

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/jinzhu/copier"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core/pagination"
	httpServer "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/carts/dtos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/carts/models"
)

// Request
type GetCartsRequest struct {
	pagination.PageRequest
}

// Result
type GetCartsResult struct {
	Body struct {
		Items []dtos.CartDto
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
			OperationID:   "GetCarts",
			Method:        http.MethodGet,
			Path:          "/carts/carts",
			Summary:       "Get Carts",
			Tags:          []string{"Carts"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, ""),
			},
		},
		getCarts(database),
	)
}

func getCarts(database *mongo.Database) func(context.Context, *struct{}) (*GetCartsResult, error) {
	return func(ctx context.Context, request *struct{}) (*GetCartsResult, error) {
		userId, ok := httpServer.GetCurrentUser(ctx)
		if !ok {
			return nil, huma.Error401Unauthorized("unable to get current user id")
		}

		cartCollection := database.Collection("carts")

		filter := bson.D{bson.E{Key: "user_id", Value: userId}}

		cursor, err := cartCollection.Find(ctx, filter)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		var carts []models.Cart
		if err = cursor.All(ctx, &carts); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		var cartDtos []dtos.CartDto
		if err := copier.Copy(&cartDtos, &carts); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		result := GetCartsResult{}
		result.Body.Items = cartDtos

		return &result, nil
	}
}
