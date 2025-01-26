package permissions

import "go.uber.org/fx"

var (
	// Module provided to fx
	Module = fx.Module(
		"permission_fx",
		permissionProviders,
	)

	permissionProviders = fx.Options(
		fx.Provide(
			NewPermissionManager,
		),
	)
)

var (
	DefaultModule = fx.Module(
		"defaultpermission_fx",
		defaultPermissionProviders,
	)

	defaultPermissionProviders = fx.Options(
		fx.Provide(
			NewDefaultPermissionDbManager,
		),
	)
)
