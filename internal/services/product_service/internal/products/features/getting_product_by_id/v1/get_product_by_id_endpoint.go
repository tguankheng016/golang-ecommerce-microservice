package v1

import (
	"context"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/dtos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/services"
)

// Request
type GetProductByIdRequest struct {
	Id int `path:"id"`
}

// Result
type GetProductByIdResult struct {
	Body struct {
		Product dtos.CreateOrEditProductDto
	}
}

// Validator
func (e GetProductByIdRequest) Schema() v.Schema {
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
			OperationID:   "GetProductById",
			Method:        http.MethodGet,
			Path:          "/products/product/{id}",
			Summary:       "Get Product By Id",
			Tags:          []string{"Products"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesProducts),
			},
		},
		getProductById(pool),
	)
}

func getProductById(pool *pgxpool.Pool) func(context.Context, *GetProductByIdRequest) (*GetProductByIdResult, error) {
	return func(ctx context.Context, request *GetProductByIdRequest) (*GetProductByIdResult, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		var productEditDto dtos.CreateOrEditProductDto

		if request.Id > 0 {
			productManager := services.NewProductManager(pool)

			product, err := productManager.GetProductById(ctx, request.Id)
			if err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}
			if product == nil {
				return nil, huma.Error404NotFound("product not found")
			}
			if err := copier.Copy(&productEditDto, &product); err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}
		} else {
			productEditDto = dtos.CreateOrEditProductDto{}
		}

		result := GetProductByIdResult{}
		result.Body.Product = productEditDto

		return &result, nil
	}
}
