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
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/models"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/services"
)

// Request
type HumaCreateCategoryRequest struct {
	Body struct {
		dtos.CreateCategoryDto
	}
}

// Result
type HumaCreateCategoryResult struct {
	Body struct {
		Category dtos.CategoryDto
	}
}

// Validator
func (e HumaCreateCategoryRequest) Schema() v.Schema {
	return v.Schema{
		v.F("id", e.Body.Id): v.Any(
			v.Zero[*int](),
			v.Nested(func(ptr *int) v.Validator { return v.Value(*ptr, v.Eq(0).Msg("Invalid category id")) }),
		).LastError(),
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
			OperationID:   "CreateCategory",
			Method:        http.MethodPost,
			Path:          "/products/category",
			Summary:       "Create Category",
			Tags:          []string{"Categories"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesCategoriesCreate),
				postgres.SetupTransaction(api, pool),
			},
		},
		createCategory(),
	)
}

func createCategory() func(context.Context, *HumaCreateCategoryRequest) (*HumaCreateCategoryResult, error) {
	return func(ctx context.Context, request *HumaCreateCategoryRequest) (*HumaCreateCategoryResult, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		tx, err := postgres.GetTxFromCtx(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		categoryManager := services.NewCategoryManager(tx)

		var category models.Category
		if err := copier.Copy(&category, &request.Body); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if err := categoryManager.CreateCategory(ctx, &category); err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		var categoryDto dtos.CategoryDto
		if err := copier.Copy(&categoryDto, &category); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		result := HumaCreateCategoryResult{}
		result.Body.Category = categoryDto

		return &result, nil
	}
}
