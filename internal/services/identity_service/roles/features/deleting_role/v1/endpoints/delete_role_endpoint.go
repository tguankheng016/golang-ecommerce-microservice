package endpoints

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	postgresgorm "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres_gorm"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/constants"
	roleModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/models"
	userModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/models"
	userService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/services"
)

func MapRoute(echo *echo.Echo, permissionManager userService.IUserRolePermissionManager) {
	group := echo.Group("/api/v1/identities/role/:roleId")
	group.DELETE("", deleteRole(permissionManager), middlewares.Authorize(permissions.PagesAdministrationRolesDelete))
}

// @ID DeleteRole
// @Tags Roles
// @Summary Delete role
// @Description Delete role
// @Accept json
// @Produce json
// @Param roleId path int true "Role Id"
// @Success 200
// @Security ApiKeyAuth
// @Router /api/v1/identities/role/{roleId} [delete]
func deleteRole(permissionManager userService.IUserRolePermissionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var roleId int64
		if err := echo.PathParamsBinder(c).Int64("roleId", &roleId).BindError(); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		tx, err := postgresgorm.GetTxFromContext(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		var role roleModel.Role
		if err := tx.First(&role, roleId).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		if strings.EqualFold(role.Name, constants.DefaultAdminRoleName) {
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("You cannot delete the default admin user"))
		}

		userManager := userService.NewUserManager(tx)
		userIds, err := userManager.GetUserIdsInRole(roleId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		for _, userId := range userIds {
			if err := userManager.RemoveToRoles(&userModel.User{Id: userId}, []int64{roleId}); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			permissionManager.RemoveUserRoleCaches(ctx, userId)
		}

		if err := tx.Delete(&role).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return c.NoContent(http.StatusOK)
	}
}
