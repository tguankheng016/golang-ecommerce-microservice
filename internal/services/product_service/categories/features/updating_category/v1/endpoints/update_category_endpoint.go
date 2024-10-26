package endpoints

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	postgresGorm "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres_gorm"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/categories/dtos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/categories/models"
)

func MapRoute(echo *echo.Echo, validator *validator.Validate) {
	group := echo.Group("/api/v1/products/category")
	group.PUT("", updateCategory(validator), middlewares.Authorize(permissions.PagesCategoriesEdit))
}

// @ID UpdateCategory
// @Tags Categories
// @Summary Update role
// @Description Update role
// @Accept json
// @Produce json
// @Param EditCategoryDto body EditCategoryDto false "EditCategoryDto"
// @Success 200 {object} CategoryDto
// @Security ApiKeyAuth
// @Router /api/v1/products/category [put]
func updateCategory(validator *validator.Validate) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var editCategoryDto dtos.EditCategoryDto

		if err := c.Bind(&editCategoryDto); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		if err := validator.StructCtx(ctx, editCategoryDto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		tx, err := postgresGorm.GetTxFromContext(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		var category models.Category
		if err := tx.First(&category, editCategoryDto.Id).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		if err := copier.Copy(&category, &editCategoryDto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := tx.Save(&category); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var categoryDto dtos.CategoryDto
		if err := copier.Copy(&categoryDto, &category); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, categoryDto)
	}
}
