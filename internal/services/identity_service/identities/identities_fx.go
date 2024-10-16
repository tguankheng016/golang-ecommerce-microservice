package identities

import (
	identityService "github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/identities/services"
	"go.uber.org/fx"
)

var (
	// Module provided to fx
	Module = fx.Module(
		"identityfx",
		identityProviders,
	)

	identityProviders = fx.Options(
		fx.Provide(
			identityService.NewCustomStampDBValidator,
			identityService.NewCustomTokenKeyDBValidator,
			identityService.NewJwtTokenGenerator,
		),
	)
)
