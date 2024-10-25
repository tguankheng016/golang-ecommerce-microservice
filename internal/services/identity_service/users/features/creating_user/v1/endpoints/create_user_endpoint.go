package endpoints

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/events"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http/echo/middlewares"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/permissions"
	postgresGorm "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres_gorm"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/rabbitmq"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/dtos"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/models"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/services"
)

func MapRoute(echo *echo.Echo, validator *validator.Validate, rabbitMQPublisher rabbitmq.IPublisher) {
	group := echo.Group("/api/v1/identities/user")
	group.POST("", createUser(validator, rabbitMQPublisher), middlewares.Authorize(permissions.PagesAdministrationUsersCreate))
}

// @ID CreateUser
// @Tags Users
// @Summary Create new user
// @Description Create new user
// @Accept json
// @Produce json
// @Param CreateUserDto body CreateUserDto false "CreateUserDto"
// @Success 200 {object} UserDto
// @Security ApiKeyAuth
// @Router /api/v1/identities/user [post]
func createUser(validator *validator.Validate, rabbitMQPublisher rabbitmq.IPublisher) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		tx, err := postgresGorm.GetTxFromContext(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		userManager := services.NewUserManager(tx)

		var createUserDto dtos.CreateUserDto

		if err := c.Bind(&createUserDto); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		if err := validator.StructCtx(ctx, createUserDto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var user models.User
		if err := copier.Copy(&user, &createUserDto); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := userManager.CreateUser(&user, createUserDto.Password); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := userManager.UpdateUserRoles(&user, createUserDto.RoleIds); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var userDto dtos.UserDto
		if err := copier.Copy(&userDto, &user); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		var userCreatedEvent events.UserCreatedEvent
		if err := copier.Copy(&userCreatedEvent, &user); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		if err := rabbitMQPublisher.PublishMessage(&userCreatedEvent); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, userDto)
	}
}
