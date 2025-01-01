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
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/dtos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/services"
)

// Request
type HumaUpdateCategoryRequest struct {
	Body struct {
		dtos.EditCategoryDto
	}
}

// Result
type HumaUpdateCategoryResult struct {
	Body struct {
		Category dtos.CategoryDto
	}
}

// Validator
func (e HumaUpdateCategoryRequest) Schema() v.Schema {
	return v.Schema{
		v.F("id", e.Body.Id): v.All(
			v.Nonzero[*int64]().Msg("Invalid category id"),
			v.Nested(func(ptr *int64) v.Validator { return v.Value(*ptr, v.Gt(int64(0)).Msg("Invalid category id")) }),
		),
		v.F("categoryname", e.Body.Name): v.Nonzero[string]().Msg("Please enter the category name"),
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
			OperationID:   "UpdateCategory",
			Method:        http.MethodPut,
			Path:          "/products/category",
			Summary:       "Update Category",
			Tags:          []string{"Categories"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesCategoriesEdit),
				postgres.SetupTransaction(api, pool),
			},
		},
		updateCategory(),
	)
}

func updateCategory() func(context.Context, *HumaUpdateCategoryRequest) (*HumaUpdateCategoryResult, error) {
	return func(ctx context.Context, request *HumaUpdateCategoryRequest) (*HumaUpdateCategoryResult, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		tx, err := postgres.GetTxFromCtx(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		categoryManager := services.NewCategoryManager(tx)

		category, err := categoryManager.GetCategoryById(ctx, *request.Body.Id)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		if category == nil {
			return nil, huma.Error404NotFound("category not found")
		}

		if err := copier.Copy(&category, &request.Body); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if err := categoryManager.UpdateCategory(ctx, category); err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		var categoryDto dtos.CategoryDto
		if err := copier.Copy(&categoryDto, &category); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		result := HumaUpdateCategoryResult{}
		result.Body.Category = categoryDto

		return &result, nil
	}
}
