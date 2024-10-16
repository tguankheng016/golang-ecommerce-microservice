package endpoints

import (
	"net/http"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	postgresGorm "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres_gorm"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/constants"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/dtos"
	roleModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/models"
	roleService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/roles/services"
	userModel "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/models"
	userService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/services"
	"gorm.io/gorm"
)

func MapRoute(echo *echo.Echo, validator *validator.Validate, permissionManager userService.IUserRolePermissionManager) {
	group := echo.Group("/api/v1/identities/role")
	group.PUT("", updateRole(validator, permissionManager), middlewares.Authorize(permissions.PagesAdministrationRolesEdit))
}

// UpdateRole
// @Tags Roles
// @Summary Update role
// @Description Update role
// @Accept json
// @Produce json
// @Param EditRoleDto body EditRoleDto false "EditRoleDto"
// @Success 200 {object} RoleDto
// @Security ApiKeyAuth
// @Router /api/v1/identities/role [put]
func updateRole(validator *validator.Validate, permissionManager userService.IUserRolePermissionManager) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var editRoleDto dtos.EditRoleDto

		if err := c.Bind(&editRoleDto); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		if err := validator.StructCtx(ctx, editRoleDto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		tx, err := postgresGorm.GetTxFromContext(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		roleManager := roleService.NewRoleManager(tx)

		var role roleModel.Role
		if err := tx.First(&role, editRoleDto.Id).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		if err := copier.Copy(&role, &editRoleDto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := permissions.ValidatePermissionName(editRoleDto.GrantedPermissions); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := roleManager.UpdateRole(&role); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		oldPermissions, err := permissionManager.SetRolePermissions(ctx, role.Id)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		isAdmin := strings.EqualFold(role.Name, constants.DefaultAdminRoleName)

		// Prohibit
		for oldPermission := range oldPermissions {
			if !slices.Contains(editRoleDto.GrantedPermissions, oldPermission) {
				if err := tx.Where("role_id = ? AND name = ?", role.Id, oldPermission).Delete(&userModel.UserRolePermission{}).Error; err != nil && err != gorm.ErrRecordNotFound {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}

				if isAdmin {
					if err := tx.Create(&userModel.UserRolePermission{RoleId: role.Id, Name: oldPermission, IsGranted: false}).Error; err != nil {
						return echo.NewHTTPError(http.StatusInternalServerError, err)
					}
				}
			}
		}

		// Granted
		for _, newPermission := range editRoleDto.GrantedPermissions {
			if _, ok := oldPermissions[newPermission]; !ok {
				var rolePermissionToGrant userModel.UserRolePermission
				if err := tx.First(&rolePermissionToGrant, userModel.UserRolePermission{RoleId: role.Id, Name: newPermission}).Error; err != nil && err != gorm.ErrRecordNotFound {
					return echo.NewHTTPError(http.StatusInternalServerError, err)
				}
				if rolePermissionToGrant.Id == 0 {
					if err := tx.Create(&userModel.UserRolePermission{RoleId: role.Id, Name: newPermission, IsGranted: true}).Error; err != nil {
						return echo.NewHTTPError(http.StatusInternalServerError, err)
					}
				} else if !rolePermissionToGrant.IsGranted {
					if isAdmin {
						if err := tx.Where("role_id = ? AND name = ?", role.Id, newPermission).Delete(&userModel.UserRolePermission{}).Error; err != nil && err != gorm.ErrRecordNotFound {
							return echo.NewHTTPError(http.StatusInternalServerError, err)
						}
					} else {
						// Unlikely will hit this case
						rolePermissionToGrant.IsGranted = true
						if err := tx.Save(&rolePermissionToGrant).Error; err != nil {
							return echo.NewHTTPError(http.StatusInternalServerError, err)
						}
					}
				}
			}
		}

		// Commit because permission manager tx is different
		tx.Commit()

		permissionManager.SetRolePermissions(ctx, role.Id)

		var roleDto dtos.RoleDto
		if err := copier.Copy(&roleDto, &role); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, roleDto)
	}
}
