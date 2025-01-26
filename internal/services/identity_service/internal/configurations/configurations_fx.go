package configurations

import "go.uber.org/fx"

var (
	// Module provided to fx
	Module = fx.Module(
		"configuration_fx",
		configurationInvokes,
	)

	configurationInvokes = fx.Options(
		fx.Invoke(ConfigureEndpoints),
		fx.Invoke(ConfigureUserGrpcServer),
		fx.Invoke(ConfigureIdentityGrpcServer),
		fx.Invoke(ConfigurePermissionGrpcServer),
	)
)
