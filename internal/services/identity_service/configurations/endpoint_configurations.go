package configurations

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	authenticate "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/identities/features/authenticating/v1/endpoints"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/identities/services"
)

func ConfigEndpoints(
	jwtTokenGenerator services.IJwtTokenGenerator,
	jwtTokenValidator jwt.IJwtTokenHandler,
	validator *validator.Validate,
	echo *echo.Echo,
) {
	// Identites
	authenticate.MapRoute(echo, validator, jwtTokenGenerator)
}
