package endpoints

import (
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core/helpers"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	postgresGorm "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres_gorm"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/dtos"
	roleModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/models"
	userService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/services"
)

type GetRoleByIdResult struct {
	Role dtos.CreateOrEditRoleDto `json:"role"`
} // @name GetRoleByIdResult

func MapRoute(echo *echo.Echo, permissionManager userService.IUserRolePermissionManager) {
	group := echo.Group("/api/v1/identities/role/:roleId")
	group.GET("", getRoleById(permissionManager), middlewares.Authorize(permissions.PagesAdministrationRoles))
}

// @ID GetRoleById
// @Tags Roles
// @Summary Get role by id
// @Description Get role by id
// @Accept json
// @Produce json
// @Param roleId path int true "Role Id"
// @Success 200 {object} GetRoleByIdResult
// @Security ApiKeyAuth
// @Router /api/v1/identities/role/{roleId} [get]
func getRoleById(permissionManager userService.IUserRolePermissionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var roleId int64
		if err := echo.PathParamsBinder(c).Int64("roleId", &roleId).BindError(); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var roleEditDto dtos.CreateOrEditRoleDto

		if roleId == 0 {
			// Create
			roleEditDto = dtos.CreateOrEditRoleDto{}
		} else {
			// Edit
			tx, err := postgresGorm.GetTxFromContext(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			var role roleModel.Role
			if err := tx.First(&role, roleId).Error; err != nil {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}

			if err := copier.Copy(&roleEditDto, &role); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}

			rolePermissions, err := permissionManager.SetRolePermissions(ctx, role.Id)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}

			roleEditDto.GrantedPermissions = helpers.MapKeysToSlice(rolePermissions)
		}

		result := &GetRoleByIdResult{
			Role: roleEditDto,
		}

		return c.JSON(http.StatusOK, result)
	}
}
