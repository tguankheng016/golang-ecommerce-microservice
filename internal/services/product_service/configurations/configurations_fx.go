package configurations

import "go.uber.org/fx"

var (
	// Module provided to fx
	Module = fx.Module(
		"configurationfx",
		configurationProviders,
	)

	configurationProviders = fx.Options(
		fx.Provide(
			ConfigIdentityGrpcClientService,
		),
	)
)
