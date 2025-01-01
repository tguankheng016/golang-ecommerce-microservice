package v1

import (
	"context"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/dtos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/services"
)

// Request
type HumaUpdateProductRequest struct {
	Body struct {
		dtos.EditProductDto
	}
}

// Result
type HumaUpdateProductResult struct {
	Body struct {
		Product dtos.ProductDto
	}
}

// Validator
func (e HumaUpdateProductRequest) Schema() v.Schema {
	return v.Schema{
		v.F("id", e.Body.Id): v.All(
			v.Nonzero[*int]().Msg("Invalid product id"),
			v.Nested(func(ptr *int) v.Validator { return v.Value(*ptr, v.Gt(0).Msg("Invalid product id")) }),
		),
		v.F("product_name", e.Body.Name):        v.Nonzero[string]().Msg("Please enter the product name"),
		v.F("product_desc", e.Body.Description): v.Nonzero[string]().Msg("Please enter the product description"),
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
			OperationID:   "UpdateProduct",
			Method:        http.MethodPut,
			Path:          "/products/product",
			Summary:       "Update Product",
			Tags:          []string{"Products"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesProductsEdit),
				postgres.SetupTransaction(api, pool),
			},
		},
		updateProduct(),
	)
}

func updateProduct() func(context.Context, *HumaUpdateProductRequest) (*HumaUpdateProductResult, error) {
	return func(ctx context.Context, request *HumaUpdateProductRequest) (*HumaUpdateProductResult, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		tx, err := postgres.GetTxFromCtx(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		productManager := services.NewProductManager(tx)

		product, err := productManager.GetProductById(ctx, *request.Body.Id)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		if product == nil {
			return nil, huma.Error404NotFound("product not found")
		}

		if err := copier.Copy(&product, &request.Body); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if err := productManager.UpdateProduct(ctx, product); err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		var productDto dtos.ProductDto
		if err := copier.Copy(&productDto, &product); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		result := HumaUpdateProductResult{}
		result.Body.Product = productDto

		return &result, nil
	}
}
