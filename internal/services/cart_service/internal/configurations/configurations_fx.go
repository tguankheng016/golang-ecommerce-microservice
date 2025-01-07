package configurations

import "go.uber.org/fx"

var (
	// Module provided to fx
	Module = fx.Module(
		"configuration_fx",
		configurationProviders,
		configurationInvokes,
	)

	configurationProviders = fx.Options(
		fx.Provide(
			ConfigIdentityGrpcClientService,
			ConfigPermissionGrpcClientService,
		),
	)

	configurationInvokes = fx.Options(
		//fx.Invoke(ConfigureMessageHandler),
		fx.Invoke(ConfigureEndpoints),
	)
)
