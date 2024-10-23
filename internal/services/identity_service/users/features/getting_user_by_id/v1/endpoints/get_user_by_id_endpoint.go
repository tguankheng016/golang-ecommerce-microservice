package endpoints

import (
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	postgresGorm "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres_gorm"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/dtos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/models"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/services"
)

type GetUserByIdResult struct {
	User dtos.CreateOrEditUserDto `json:"user"`
} // @name GetUserByIdResult

func MapRoute(echo *echo.Echo) {
	group := echo.Group("/api/v1/identities/user/:userId")
	group.GET("", getUserById(), middlewares.Authorize(permissions.PagesAdministrationUsers))
}

// @ID GetUserById
// @Tags Users
// @Summary Get user by id
// @Description Get user by id
// @Accept json
// @Produce json
// @Param userId path int true "User Id"
// @Success 200 {object} GetUserByIdResult
// @Security ApiKeyAuth
// @Router /api/v1/identities/user/{userId} [get]
func getUserById() echo.HandlerFunc {
	return func(c echo.Context) error {
		var userId int64
		if err := echo.PathParamsBinder(c).Int64("userId", &userId).BindError(); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var userEditDto dtos.CreateOrEditUserDto

		if userId == 0 {
			// Create
			userEditDto = dtos.CreateOrEditUserDto{}
		} else {
			// Edit
			tx, err := postgresGorm.GetTxFromContext(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			var user models.User
			if err := tx.First(&user, userId).Error; err != nil {
				return echo.NewHTTPError(http.StatusNotFound, err)
			}

			if err := copier.Copy(&userEditDto, &user); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}

			userManager := services.NewUserManager(tx)

			userEditDto.RoleIds, err = userManager.GetUserRoles(&user)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		}

		result := &GetUserByIdResult{
			User: userEditDto,
		}

		return c.JSON(http.StatusOK, result)
	}
}
