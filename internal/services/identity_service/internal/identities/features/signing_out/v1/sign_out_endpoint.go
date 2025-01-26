package v1

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	httpServer "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/identities/services"
)

// Handler
func MapRoute(
	api huma.API,
	jwtTokenGenerator services.IJwtTokenGenerator,
) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "SignOut",
			Method:        http.MethodPost,
			Path:          "/identities/sign-out",
			Summary:       "SignOut",
			Tags:          []string{"Identities"},
			DefaultStatus: http.StatusOK,
			Security: []map[string][]string{
				{"bearer": {}},
			},
		},
		signOut(jwtTokenGenerator),
	)
}

func signOut(jwtTokenGenerator services.IJwtTokenGenerator) func(context.Context, *struct{}) (*struct{}, error) {
	return func(ctx context.Context, input *struct{}) (*struct{}, error) {
		userId, ok := httpServer.GetCurrentUser(ctx)

		if ok {
			claims, ok := httpServer.GetCurrentUserClaims(ctx)

			if !ok {
				return nil, huma.Error401Unauthorized("unable to get user claims")
			}

			if err := jwtTokenGenerator.RemoveUserTokens(ctx, userId, claims); err != nil {
				return nil, huma.Error500InternalServerError(err.Error())
			}
		}

		return nil, nil
	}
}
