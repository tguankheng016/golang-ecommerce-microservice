package endpoints

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/events"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/openiddict"
	postgresGorm "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/postgres_gorm"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/rabbitmq"
	identityConsts "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/constants"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/identities/services"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/models"
	userService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/users/services"
	"gorm.io/gorm"
)

type OAuthAuthenticateRequest struct {
	Code        string `json:"code" validate:"required"`
	RedirectUri string `json:"redirect_uri" validate:"required"`
} // @name OAuthAuthenticateRequest

type OAuthAuthenticateResult struct {
	AccessToken                 string `json:"access_token"`
	ExpireInSeconds             int    `json:"expire_in_seconds"`
	RefreshToken                string `json:"refresh_token"`
	RefreshTokenExpireInSeconds int    `json:"refresh_token_expire_in_seconds"`
} // @name OAuthAuthenticateResult

func MapRoute(
	echo *echo.Echo,
	validator *validator.Validate,
	jwtTokenGenerator services.IJwtTokenGenerator,
	oAuthApiClient openiddict.IOAuthApiClient,
	rabbitMQPublisher rabbitmq.IPublisher,
) {
	group := echo.Group("/api/v1/identities/oauth-authenticate")
	group.POST("", oauthAuthenticate(validator, jwtTokenGenerator, oAuthApiClient, rabbitMQPublisher))
}

// OAuthAuthenticate
// @Tags Identities
// @Summary OAuthAuthenticate
// @Description OAuthAuthenticate
// @Accept json
// @Produce json
// @Param OAuthAuthenticateRequest body OAuthAuthenticateRequest true "OAuthAuthenticateRequest"
// @Success 200 {object} OAuthAuthenticateResult
// @Security ApiKeyAuth
// @Router /api/v1/identities/oauth-authenticate [post]
func oauthAuthenticate(
	validator *validator.Validate,
	jwtTokenGenerator services.IJwtTokenGenerator,
	oAuthApiClient openiddict.IOAuthApiClient,
	rabbitMQPublisher rabbitmq.IPublisher,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		request := &OAuthAuthenticateRequest{}

		if err := c.Bind(request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := validator.StructCtx(ctx, request); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		tx, err := postgresGorm.GetTxFromContext(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		userManager := userService.NewUserManager(tx)

		tokenResponse, err := oAuthApiClient.ConnectToken(ctx, request.Code, request.RedirectUri)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		userInfo, err := oAuthApiClient.ConnectUserInfo(ctx, tokenResponse.AccessToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		externalUserId, err := uuid.FromString(userInfo.Sub)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		var user models.User
		if err := tx.Where("external_user_id = ?", externalUserId).First(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		if user.Id == 0 && userInfo.PreferredUsername != identityConsts.DefaultAdminUsername {
			// New User
			newUser := &models.User{
				UserName:  userInfo.PreferredUsername,
				FirstName: userInfo.GivenName,
				LastName:  userInfo.FamilyName,
				Email:     userInfo.Email,
			}

			randomPassword, err := userService.GenerateRandomPassword(10)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			if err := userManager.CreateUser(newUser, randomPassword); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			var userCreatedEvent events.UserCreatedEvent
			if err := copier.Copy(&userCreatedEvent, &user); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
			if err := rabbitMQPublisher.PublishMessage(&userCreatedEvent); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		} else {
			if user.Id == 0 && userInfo.PreferredUsername == identityConsts.DefaultAdminUsername {
				if err := tx.Where("user_name = ?", identityConsts.DefaultAdminUsername).First(&user).Error; err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, err)
				}

				user.ExternalUserId = externalUserId
			}

			user.UserName = userInfo.PreferredUsername
			user.Email = userInfo.Email
			user.FirstName = userInfo.GivenName
			user.LastName = userInfo.FamilyName

			if err := userManager.UpdateUser(&user, ""); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}

			var userUpdatedEvent events.UserUpdatedEvent
			if err := copier.Copy(&userUpdatedEvent, &user); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
			if err := rabbitMQPublisher.PublishMessage(&userUpdatedEvent); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}
		}

		refreshToken, refreshTokenKey, refreshTokenSeconds, err := jwtTokenGenerator.GenerateRefreshToken(ctx, &user)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		accessToken, accessTokenSeconds, err := jwtTokenGenerator.GenerateAccessToken(ctx, &user, refreshTokenKey)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		result := &OAuthAuthenticateResult{
			AccessToken:                 accessToken,
			ExpireInSeconds:             accessTokenSeconds,
			RefreshToken:                refreshToken,
			RefreshTokenExpireInSeconds: refreshTokenSeconds,
		}

		return c.JSON(http.StatusOK, result)
	}
}
