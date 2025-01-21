package v1

import (
	"context"
	"encoding/json"
	"net/http"

	v "github.com/RussellLuo/validating/v3"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/danielgtaylor/huma/v2"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/events"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/openiddict"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/identities/services"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/models"
	userService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/users/services"
)

// Request
type OAuthAuthenticateRequest struct {
	Code        string `json:"code"`
	RedirectUri string `json:"redirectUri"`
}
type HumaOAuthAuthenticateRequest struct {
	Body struct {
		OAuthAuthenticateRequest
	}
}

// Result
type AuthenticateResult struct {
	AccessToken                 string `json:"accessToken"`
	ExpireInSeconds             int    `json:"expireInSeconds"`
	RefreshToken                string `json:"refreshToken"`
	RefreshTokenExpireInSeconds int    `json:"refreshTokenExpireInSeconds"`
}
type HumaOAuthAuthenticateResult struct {
	Body struct {
		AuthenticateResult
	}
}

// Validator
func (e HumaOAuthAuthenticateRequest) Schema() v.Schema {
	return v.Schema{
		v.F("code", e.Body.Code):                v.Nonzero[string]().Msg("Please enter the authorization code"),
		v.F("redirect_uri", e.Body.RedirectUri): v.Nonzero[string]().Msg("Please enter the redirect uri"),
	}
}

// Handler
func MapRoute(
	api huma.API,
	pool *pgxpool.Pool,
	jwtTokenGenerator services.IJwtTokenGenerator,
	oAuthApiClient openiddict.IOAuthApiClient,
	publisher message.Publisher,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "OAuthAuthenticate",
			Method:        http.MethodPost,
			Path:          "/identities/oauth-authenticate",
			Summary:       "OAuthAuthenticate",
			Tags:          []string{"Identities"},
			DefaultStatus: http.StatusOK,
		},
		authenticate(jwtTokenGenerator, pool, oAuthApiClient, publisher),
	)
}

func authenticate(
	jwtTokenGenerator services.IJwtTokenGenerator,
	pool *pgxpool.Pool,
	oAuthApiClient openiddict.IOAuthApiClient,
	publisher message.Publisher,
) func(context.Context, *HumaOAuthAuthenticateRequest) (*HumaOAuthAuthenticateResult, error) {
	return func(ctx context.Context, request *HumaOAuthAuthenticateRequest) (*HumaOAuthAuthenticateResult, error) {
		errs := v.Validate(request.Schema())
		for _, err := range errs {
			return nil, huma.Error400BadRequest(err.Message())
		}

		userManager := userService.NewUserManager(pool)

		tokenResponse, err := oAuthApiClient.ConnectToken(ctx, request.Body.Code, request.Body.RedirectUri)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error(), err)
		}

		userInfo, err := oAuthApiClient.ConnectUserInfo(ctx, tokenResponse.AccessToken)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error(), err)
		}

		user, err := userManager.GetUserByExternalUserId(ctx, userInfo.Sub)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		if user.Id == 0 {
			// New User
			newUser := &models.User{
				UserName:  userInfo.PreferredUsername,
				FirstName: userInfo.GivenName,
				LastName:  userInfo.FamilyName,
				Email:     userInfo.Email,
			}

			randomPassword, err := uuid.NewV4()
			if err != nil {
				return nil, huma.Error500InternalServerError(err.Error(), err)
			}

			if err := userManager.CreateUser(ctx, newUser, randomPassword.String()); err != nil {
				return nil, huma.Error400BadRequest(err.Error())
			}

			var userCreatedEvent events.UserCreatedEvent
			if err := copier.Copy(&userCreatedEvent, &user); err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}

			payload, err := json.Marshal(userCreatedEvent)
			if err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}

			msg := message.NewMessage(watermill.NewUUID(), payload)
			publisher.Publish(events.UserCreatedTopicV1, msg)
		} else {
			user.UserName = userInfo.PreferredUsername
			user.Email = userInfo.Email
			user.FirstName = userInfo.GivenName
			user.LastName = userInfo.FamilyName

			if err := userManager.UpdateUser(ctx, user, ""); err != nil {
				return nil, huma.Error400BadRequest(err.Error())
			}

			var userUpdatedEvent events.UserUpdatedEvent
			if err := copier.Copy(&userUpdatedEvent, &user); err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}

			payload, err := json.Marshal(userUpdatedEvent)
			if err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}

			msg := message.NewMessage(watermill.NewUUID(), payload)
			publisher.Publish(events.UserUpdatedTopicV1, msg)
		}

		refreshToken, refreshTokenKey, refreshTokenSeconds, err := jwtTokenGenerator.GenerateRefreshToken(ctx, *user)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error(), err)
		}

		accessToken, accessTokenSeconds, err := jwtTokenGenerator.GenerateAccessToken(ctx, *user, refreshTokenKey)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error(), err)
		}

		result := HumaOAuthAuthenticateResult{}

		result.Body.AccessToken = accessToken
		result.Body.ExpireInSeconds = accessTokenSeconds
		result.Body.RefreshToken = refreshToken
		result.Body.RefreshTokenExpireInSeconds = refreshTokenSeconds

		return &result, nil
	}
}
