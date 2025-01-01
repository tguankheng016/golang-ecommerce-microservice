package v1

import (
	"context"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/internal/categories/services"
)

// Request
type DeleteCategoryRequest struct {
	Id int `path:"id"`
}

// Validator
func (e DeleteCategoryRequest) Schema() v.Schema {
	return v.Schema{
		v.F("id", e.Id): v.Gte(int64(0)).Msg("Invalid category id"),
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
			OperationID:   "DeleteCategory",
			Method:        http.MethodDelete,
			Path:          "/products/category/{id}",
			Summary:       "Delete Category",
			Tags:          []string{"Categories"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
			Middlewares: huma.Middlewares{
				permissions.Authorize(api, permissions.PagesCategoriesDelete),
				postgres.SetupTransaction(api, pool),
			},
		},
		deleteCategory(),
	)
}

func deleteCategory() func(context.Context, *DeleteCategoryRequest) (*struct{}, error) {
	return func(ctx context.Context, request *DeleteCategoryRequest) (*struct{}, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		tx, err := postgres.GetTxFromCtx(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		categoryManager := services.NewCategoryManager(tx)

		category, err := categoryManager.GetCategoryById(ctx, request.Id)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}
		if category == nil {
			return nil, huma.Error404NotFound("category not found")
		}

		if err := categoryManager.DeleteCategory(ctx, category.Id); err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return nil, nil
	}
}
