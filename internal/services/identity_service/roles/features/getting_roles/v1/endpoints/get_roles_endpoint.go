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
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/dtos"
	roleModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/models"
)

type GetRolesRequest struct {
	*pagination.PageRequest
}

type GetRolesResult struct {
	*pagination.PageResultDto[dtos.RoleDto]
} // @name GetRolesResult

func MapRoute(echo *echo.Echo, validator *validator.Validate) {
	group := echo.Group("/api/v1/identities/roles")
	group.GET("", getAllRoles(validator), middlewares.Authorize(permissions.PagesAdministrationRoles))
}

// @ID GetAllRoles
// @Tags Roles
// @Summary Get all roles
// @Description Get all roles
// @Accept json
// @Produce json
// @Param GetRolesRequest query GetRolesRequest false "GetRolesRequest"
// @Success 200 {object} GetRolesResult
// @Security ApiKeyAuth
// @Router /api/v1/identities/roles [get]
func getAllRoles(validator *validator.Validate) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		tx, err := postgresGorm.GetTxFromContext(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		var roles []roleModel.Role

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

		rolePageRequest := &GetRolesRequest{PageRequest: pageRequest}

		query := tx
		countQuery := tx.Model(&roleModel.Role{})

		if rolePageRequest.Filters != "" {
			likeExpr := rolePageRequest.BuildFiltersExpr(fields...)
			query = query.Where(likeExpr)
			countQuery = countQuery.Where(likeExpr)
		}

		if rolePageRequest.Sorting != "" {
			query = query.Order(rolePageRequest.Sorting)
		}

		if rolePageRequest.SkipCount > 0 || rolePageRequest.MaxResultCount > 0 {
			query = rolePageRequest.Paginate(query)
		}

		var count int64

		if err := countQuery.Count(&count).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := query.Find(&roles).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var roleDtos []dtos.RoleDto
		if err := copier.Copy(&roleDtos, &roles); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		result := &GetRolesResult{
			pagination.NewPageResultDto(roleDtos, count),
		}

		return c.JSON(http.StatusOK, result)
	}
}
