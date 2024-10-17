package configurations

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	getting_categories "github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/categories/features/getting_categories/v1/endpoints"
)

func ConfigEndpoints(
	validator *validator.Validate,
	echo *echo.Echo,
) {
	// Categories
	getting_categories.MapRoute(echo, validator)
}
