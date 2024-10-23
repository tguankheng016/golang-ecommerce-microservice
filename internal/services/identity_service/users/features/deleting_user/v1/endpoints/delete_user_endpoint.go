package endpoints

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/events"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	postgresGorm "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres_gorm"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/rabbitmq"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/constants"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/models"
)

func MapRoute(echo *echo.Echo, rabbitMQPublisher rabbitmq.IPublisher) {
	group := echo.Group("/api/v1/identities/user/:userId")
	group.DELETE("", deleteUser(rabbitMQPublisher), middlewares.Authorize(permissions.PagesAdministrationUsersDelete))
}

// @ID DeleteUser
// @Tags Users
// @Summary Delete user
// @Description Delete user
// @Accept json
// @Produce json
// @Param userId path int true "User Id"
// @Success 200
// @Security ApiKeyAuth
// @Router /api/v1/identities/user/{userId} [delete]
func deleteUser(rabbitMQPublisher rabbitmq.IPublisher) echo.HandlerFunc {
	return func(c echo.Context) error {
		var userId int64
		if err := echo.PathParamsBinder(c).Int64("userId", &userId).BindError(); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		tx, err := postgresGorm.GetTxFromContext(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		var user models.User
		if err := tx.First(&user, userId).Error; err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}

		if user.NormalizedUserName == strings.ToUpper(constants.DefaultAdminUsername) {
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("You cannot delete the default admin user"))
		}

		if err := tx.Delete(&user).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		userDeletedEvent := &events.UserDeletedEvent{
			Id: userId,
		}
		if err := rabbitMQPublisher.PublishMessage(&userDeletedEvent); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.NoContent(http.StatusOK)
	}
}
