package v1

import (
	"context"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/services"
)

// Request
type DeleteProductRequest struct {
	Id int `path:"id"`
}

// Validator
func (e DeleteProductRequest) Schema() v.Schema {
	return v.Schema{
		v.F("id", e.Id): v.Gte(0).Msg("Invalid product id"),
	}
}

// Handler
func MapRoute(
	api huma.API,
	pool *pgxpool.Pool,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "DeleteProduct",
			Method:        http.MethodDelete,
			Path:          "/products/product/{id}",
			Summary:       "Delete Product",
			Tags:          []string{"Products"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesProductsDelete),
				postgres.SetupTransaction(api, pool),
			},
		},
		deleteProduct(),
	)
}

func deleteProduct() func(context.Context, *DeleteProductRequest) (*struct{}, error) {
	return func(ctx context.Context, request *DeleteProductRequest) (*struct{}, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		tx, err := postgres.GetTxFromCtx(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		productManager := services.NewProductManager(tx)

		product, err := productManager.GetProductById(ctx, request.Id)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		if product == nil {
			return nil, huma.Error404NotFound("product not found")
		}

		if err := productManager.DeleteProduct(ctx, product.Id); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return nil, nil
	}
}
