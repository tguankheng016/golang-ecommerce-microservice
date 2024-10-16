package security

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	"go.uber.org/fx"
)

var (
	// Module provided to fx
	Module = fx.Module(
		"securityfx",
		securityProviders,
	)

	securityProviders = fx.Options(
		fx.Provide(
			jwt.NewTokenHandler,
			jwt.NewTokenKeyValidator,
			jwt.NewSecurityStampValidator,
		),
	)
)
