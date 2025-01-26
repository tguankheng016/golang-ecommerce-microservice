package security

import (
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/security/jwt"
	"go.uber.org/fx"
)

var (
	// Module provided to fx
	Module = fx.Module(
		"security_fx",
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

var (
	DefaultModule = fx.Module(
		"defaultsecurity_fx",
		defaultSecurityProviders,
	)

	defaultSecurityProviders = fx.Options(
		fx.Provide(
			jwt.NewDefaultStampDBValidator,
			jwt.NewDefaultTokenKeyDBValidator,
		),
	)
)
