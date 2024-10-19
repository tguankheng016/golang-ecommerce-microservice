package openiddict

import "go.uber.org/fx"

var (
	// Module provided to fx
	Module = fx.Module(
		"openiddictfx",
		openIddictProviders,
	)

	openIddictProviders = fx.Options(
		fx.Provide(
			NewOAuthApiClient,
		),
	)
)
