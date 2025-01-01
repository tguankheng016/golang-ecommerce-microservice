package v1

import (
	"context"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/dtos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/services"
)

// Request
type GetCategoryByIdRequest struct {
	Id int `path:"id"`
}

// Result
type GetCategoryByIdResult struct {
	Body struct {
		Category dtos.CreateOrEditCategoryDto
	}
}

// Validator
func (e GetCategoryByIdRequest) Schema() v.Schema {
	return v.Schema{
		v.F("id", e.Id): v.Gte(0).Msg("Invalid category id"),
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
			OperationID:   "GetCategoryById",
			Method:        http.MethodGet,
			Path:          "/products/category/{id}",
			Summary:       "Get Category By Id",
			Tags:          []string{"Categories"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesCategories),
			},
		},
		getCategoryById(pool),
	)
}

func getCategoryById(pool *pgxpool.Pool) func(context.Context, *GetCategoryByIdRequest) (*GetCategoryByIdResult, error) {
	return func(ctx context.Context, request *GetCategoryByIdRequest) (*GetCategoryByIdResult, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		var categoryEditDto dtos.CreateOrEditCategoryDto

		if request.Id > 0 {
			categoryManager := services.NewCategoryManager(pool)

			category, err := categoryManager.GetCategoryById(ctx, request.Id)
			if err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}
			if category == nil {
				return nil, huma.Error404NotFound("category not found")
			}
			if err := copier.Copy(&categoryEditDto, &category); err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}
		} else {
			categoryEditDto = dtos.CreateOrEditCategoryDto{}
		}

		result := GetCategoryByIdResult{}
		result.Body.Category = categoryEditDto

		return &result, nil
	}
}
