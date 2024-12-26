package identities

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/internal/identities/services"
	"go.uber.org/fx"
)

var (
	// Module provided to fx
	Module = fx.Module(
		"identity_fx",
		identityProviders,
	)

	identityProviders = fx.Options(
		fx.Provide(
			services.NewTokenKeyDBValidator,
			services.NewSecurityStampDbValidator,
			services.NewJwtTokenGenerator,
			services.NewPermissionDbManager,
		),
	)
)
