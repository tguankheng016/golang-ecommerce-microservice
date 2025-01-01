package configurations

import "go.uber.org/fx"

var (
	// Module provided to fx
	Module = fx.Module(
		"configuration_fx",
		configurationInvokes,
	)

	configurationInvokes = fx.Options(
		fx.Invoke(ConfigureMessageHandler),
		//fx.Invoke(ConfigureEndpoints),
	)
)
