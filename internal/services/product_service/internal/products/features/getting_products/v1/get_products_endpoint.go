package v1

import (
	"context"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core/pagination"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/dtos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/products/services"
)

// Request
type GetProductsRequest struct {
	pagination.PageRequest
}

// Result
type GetProductsResult struct {
	Body struct {
		pagination.PageResultDto[dtos.ProductDto]
	}
}

// Validator
func (e GetProductsRequest) Schema() v.Schema {
	return v.Schema{
		v.F("skip_count", e.SkipCount):            v.Gte(0).Msg("Page should at least greater than or equal to 0."),
		v.F("max_result_count", e.MaxResultCount): v.Gte(0).Msg("Page size should at least greater than or equal to 0."),
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
			OperationID:   "GetProducts",
			Method:        http.MethodGet,
			Path:          "/products/products",
			Summary:       "Get Products",
			Tags:          []string{"Products"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesProducts),
			},
		},
		getProducts(pool),
	)
}

func getProducts(pool *pgxpool.Pool) func(context.Context, *GetProductsRequest) (*GetProductsResult, error) {
	return func(ctx context.Context, request *GetProductsRequest) (*GetProductsResult, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		productManager := services.NewProductManager(pool)

		products, count, err := productManager.GetProductsWithCategory(ctx, &request.PageRequest)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		var productDtos []dtos.ProductDto

		for _, p := range products {
			var productDto dtos.ProductDto
			if err := copier.Copy(&productDto, &p); err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}

			productDto.CategoryName = p.CategoryFK.Name

			productDtos = append(productDtos, productDto)
		}

		result := GetProductsResult{}
		result.Body.PageResultDto = pagination.NewPageResultDto(productDtos, count)

		return &result, nil
	}
}
