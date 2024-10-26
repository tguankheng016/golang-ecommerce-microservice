package configurations

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	creating_category "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/categories/features/creating_category/v1/endpoints"
	deleting_category "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/categories/features/deleting_category/v1/endpoints"
	getting_categories "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/categories/features/getting_categories/v1/endpoints"
	getting_category_by_id "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/categories/features/getting_category_by_id/v1/endpoints"
	updating_category "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/categories/features/updating_category/v1/endpoints"
)

func ConfigEndpoints(
	validator *validator.Validate,
	echo *echo.Echo,
) {
	// Categories
	getting_category_by_id.MapRoute(echo)
	getting_categories.MapRoute(echo, validator)
	creating_category.MapRoute(echo, validator)
	updating_category.MapRoute(echo, validator)
	deleting_category.MapRoute(echo)
}
