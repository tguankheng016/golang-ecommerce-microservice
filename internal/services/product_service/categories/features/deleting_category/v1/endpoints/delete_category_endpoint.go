package endpoints

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	postgresGorm "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres_gorm"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/categories/models"
)

func MapRoute(echo *echo.Echo) {
	group := echo.Group("/api/v1/products/category/:categoryId")
	group.DELETE("", deleteCategory(), middlewares.Authorize(permissions.PagesCategoriesDelete))
}

// @ID DeleteCategory
// @Tags Categories
// @Summary Delete category
// @Description Delete category
// @Accept json
// @Produce json
// @Param categoryId path int true "Category Id"
// @Success 200
// @Security ApiKeyAuth
// @Router /api/v1/products/category/{categoryId} [delete]
func deleteCategory() echo.HandlerFunc {
	return func(c echo.Context) error {
		var categoryId int64
		if err := echo.PathParamsBinder(c).Int64("categoryId", &categoryId).BindError(); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		tx, err := postgresGorm.GetTxFromContext(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		var category models.Category
		if err := tx.First(&category, categoryId).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		if err := tx.Delete(&category).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.NoContent(http.StatusOK)
	}
}
