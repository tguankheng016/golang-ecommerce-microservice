package configurations

import "go.uber.org/fx"

var (
	// Module provided to fx
	Module = fx.Module(
		"configurationfx",
		configurationInvokes,
	)

	configurationInvokes = fx.Options(
		fx.Invoke(ConfigGrpcClients),
	)
)
