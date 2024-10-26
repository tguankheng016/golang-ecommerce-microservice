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
	group.POST("", createCategory(validator), middlewares.Authorize(permissions.PagesCategoriesCreate))
}

// @ID CreateCategory
// @Tags Categories
// @Summary Create new category
// @Description Create new category
// @Accept json
// @Produce json
// @Param CreateCategoryDto body CreateCategoryDto false "CreateCategoryDto"
// @Success 200 {object} CategoryDto
// @Security ApiKeyAuth
// @Router /api/v1/products/category [post]
func createCategory(validator *validator.Validate) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		tx, err := postgresGorm.GetTxFromContext(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		var createCategoryDto dtos.CreateCategoryDto

		if err := c.Bind(&createCategoryDto); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		if err := validator.StructCtx(ctx, createCategoryDto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var category models.Category
		if err := copier.Copy(&category, &createCategoryDto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := tx.Create(&category); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var categoryDto dtos.CategoryDto
		if err := copier.Copy(&categoryDto, &category); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, categoryDto)
	}
}
