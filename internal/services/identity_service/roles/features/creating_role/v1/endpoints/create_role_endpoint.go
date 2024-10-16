package endpoints

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	postgresGorm "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres_gorm"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/dtos"
	roleModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/models"
	data "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/services"
	userModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/models"
)

func MapRoute(echo *echo.Echo, validator *validator.Validate) {
	group := echo.Group("/api/v1/identities/role")
	group.POST("", createRole(validator), middlewares.Authorize(permissions.PagesAdministrationRolesCreate))
}

// CreateRole
// @Tags Roles
// @Summary Create new role
// @Description Create new role
// @Accept json
// @Produce json
// @Param CreateRoleDto body CreateRoleDto false "CreateRoleDto"
// @Success 200 {object} RoleDto
// @Security ApiKeyAuth
// @Router /api/v1/identities/role [post]
func createRole(validator *validator.Validate) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		tx, err := postgresGorm.GetTxFromContext(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		roleManager := data.NewRoleManager(tx)

		var createRoleDto dtos.CreateRoleDto

		if err := c.Bind(&createRoleDto); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		if err := validator.StructCtx(ctx, createRoleDto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var role roleModel.Role
		if err := copier.Copy(&role, &createRoleDto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := permissions.ValidatePermissionName(createRoleDto.GrantedPermissions); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := roleManager.CreateRole(&role); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if len(createRoleDto.GrantedPermissions) > 0 {
			for _, permission := range createRoleDto.GrantedPermissions {
				newUserRolePermission := &userModel.UserRolePermission{
					RoleId:    role.Id,
					Name:      permission,
					IsGranted: true,
				}
				if err := tx.Model(&userModel.UserRolePermission{}).Create(&newUserRolePermission).Error; err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, err)
				}
			}
		}

		var roleDto dtos.RoleDto
		if err := copier.Copy(&roleDto, &role); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, roleDto)
	}
}
