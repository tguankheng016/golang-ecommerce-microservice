package endpoints

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core/pagination"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	postgresGorm "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres_gorm"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/categories/dtos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/product_service/categories/models"
)

type GetCategoriesRequest struct {
	*pagination.PageRequest
}

type GetCategoriesResult struct {
	*pagination.PageResultDto[dtos.CategoryDto]
} // @name GetCategoriesResult

func MapRoute(echo *echo.Echo, validator *validator.Validate) {
	group := echo.Group("/api/v1/products/categories")
	group.GET("", getAllCategories(validator), middlewares.Authorize(permissions.PagesCategories))
}

// @ID GetAllCategories
// @Tags Categories
// @Summary Get all categories
// @Description Get all categories
// @Accept json
// @Produce json
// @Param GetCategoriesRequest query GetCategoriesRequest false "GetCategoriesRequest"
// @Success 200 {object} GetCategoriesResult
// @Security ApiKeyAuth
// @Router /api/v1/products/categories [get]
func getAllCategories(validator *validator.Validate) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		tx, err := postgresGorm.GetTxFromContext(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		var categories []models.Category

		pageRequest, err := pagination.GetPageRequestFromCtx(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := validator.StructCtx(ctx, pageRequest); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		fields := []string{"name"}

		if err := pageRequest.SanitizeSorting(fields...); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		categoryPageRequest := &GetCategoriesRequest{PageRequest: pageRequest}

		query := tx
		countQuery := tx.Model(&models.Category{})

		if categoryPageRequest.Filters != "" {
			likeExpr := categoryPageRequest.BuildFiltersExpr(fields...)
			query = query.Where(likeExpr)
			countQuery = countQuery.Where(likeExpr)
		}

		if categoryPageRequest.Sorting != "" {
			query = query.Order(categoryPageRequest.Sorting)
		}

		if categoryPageRequest.SkipCount > 0 || categoryPageRequest.MaxResultCount > 0 {
			query = categoryPageRequest.Paginate(query)
		}

		var count int64

		if err := countQuery.Count(&count).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := query.Find(&categories).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var categoryDtos []dtos.CategoryDto
		if err := copier.Copy(&categoryDtos, &categories); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		result := &GetCategoriesResult{
			pagination.NewPageResultDto(categoryDtos, count),
		}

		return c.JSON(http.StatusOK, result)
	}
}
