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
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/dtos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/services"
)

// Request
type GetCategoriesRequest struct {
	pagination.PageRequest
}

// Result
type GetCategoriesResult struct {
	Body struct {
		pagination.PageResultDto[dtos.CategoryDto]
	}
}

// Validator
func (e GetCategoriesRequest) Schema() v.Schema {
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
			OperationID:   "GetCategories",
			Method:        http.MethodGet,
			Path:          "/products/categories",
			Summary:       "Get Categories",
			Tags:          []string{"Categories"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesCategories),
			},
		},
		getCategories(pool),
	)
}

func getCategories(pool *pgxpool.Pool) func(context.Context, *GetCategoriesRequest) (*GetCategoriesResult, error) {
	return func(ctx context.Context, request *GetCategoriesRequest) (*GetCategoriesResult, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		categoryManager := services.NewCategoryManager(pool)

		categories, count, err := categoryManager.GetCategories(ctx, &request.PageRequest)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		var categoryDtos []dtos.CategoryDto
		if err := copier.Copy(&categoryDtos, &categories); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		result := GetCategoriesResult{}
		result.Body.PageResultDto = pagination.NewPageResultDto(categoryDtos, count)

		return &result, nil
	}
}
