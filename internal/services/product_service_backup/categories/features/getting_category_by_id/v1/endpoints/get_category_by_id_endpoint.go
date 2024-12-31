package endpoints

import (
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	postgresGorm "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres_gorm"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/categories/dtos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/categories/models"
)

type GetCategoryByIdResult struct {
	Category dtos.CreateOrEditCategoryDto `json:"category"`
} // @name GetCategoryByIdResult

func MapRoute(echo *echo.Echo) {
	group := echo.Group("/api/v1/products/category/:categoryId")
	group.GET("", getCategoryById(), middlewares.Authorize(permissions.PagesCategories))
}

// @ID GetCategoryById
// @Tags Categories
// @Summary Get category by id
// @Description Get category by id
// @Accept json
// @Produce json
// @Param categoryId path int true "Category Id"
// @Success 200 {object} GetCategoryByIdResult
// @Security ApiKeyAuth
// @Router /api/v1/products/category/{categoryId} [get]
func getCategoryById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var categoryId int64
		if err := echo.PathParamsBinder(c).Int64("categoryId", &categoryId).BindError(); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var categoryEditDto dtos.CreateOrEditCategoryDto

		if categoryId == 0 {
			// Create
			categoryEditDto = dtos.CreateOrEditCategoryDto{}
		} else {
			// Edit
			tx, err := postgresGorm.GetTxFromContext(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			var category models.Category
			if err := tx.First(&category, categoryId).Error; err != nil {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}

			if err := copier.Copy(&categoryEditDto, &category); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		}

		result := &GetCategoryByIdResult{
			Category: categoryEditDto,
		}

		return c.JSON(http.StatusOK, result)
	}
}
