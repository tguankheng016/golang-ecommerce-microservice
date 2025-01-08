package v1

import (
	"context"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/gofrs/uuid"
	httpServer "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	cartModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/carts/models"
	productModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/cart_service/internal/products/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Request
type AddCartDto struct {
	ProductId int `json:"productId"`
}

type HumaAddProductRequest struct {
	Body struct {
		AddCartDto
	}
}

// Validator
func (e HumaAddProductRequest) Schema() v.Schema {
	return v.Schema{
		v.F("productId", e.Body.ProductId): v.Gt(0).Msg("invalid product id"),
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
			OperationID:   "AddCart",
			Method:        http.MethodPost,
			Path:          "/carts/cart",
			Summary:       "Add Cart",
			Tags:          []string{"Carts"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, ""),
			},
		},
		addCart(database),
	)
}

func addCart(database *mongo.Database) func(context.Context, *HumaAddProductRequest) (*struct{}, error) {
	return func(ctx context.Context, request *HumaAddProductRequest) (*struct{}, error) {
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

		filter := bson.D{bson.E{Key: "id", Value: request.Body.ProductId}}

		var product productModel.Product
		err := productCollection.FindOne(ctx, filter).Decode(&product)

		if err != nil && err != mongo.ErrNoDocuments {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if err == mongo.ErrNoDocuments {
			return nil, huma.Error400BadRequest("invalid product id")
		}

		filter = bson.D{
			bson.E{Key: "productid", Value: request.Body.ProductId},
			bson.E{Key: "userid", Value: userId},
		}
		var cart cartModel.Cart
		err = cartCollection.FindOne(ctx, filter).Decode(&cart)

		if err != nil && err != mongo.ErrNoDocuments {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if err == mongo.ErrNoDocuments {
			newCartId, err := uuid.NewV6()
			if err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}

			newCart := cartModel.Cart{
				Id:           newCartId.String(),
				ProductId:    product.Id,
				ProductName:  product.Name,
				ProductDesc:  product.Description,
				ProductPrice: product.Price,
				Quantity:     1,
				UserId:       userId,
			}
			_, err = cartCollection.InsertOne(ctx, newCart)
			if err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}
		} else {
			update := bson.D{
				{
					Key: "$set",
					Value: bson.D{
						bson.E{Key: "quantity", Value: cart.Quantity + 1},
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
