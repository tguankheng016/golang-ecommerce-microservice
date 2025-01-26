package messaging

import "go.uber.org/fx"

var (
	// Module provided to fx
	Module = fx.Module(
		"messaging_fx",
		messagingProviders,
		messagingInvokes,
	)

	messagingProviders = fx.Options(
		fx.Provide(
			NewWatermillLogger,
			NewWatermillRouter,
			NewWatermillPublisher,
			NewWatermillSubscriber,
		),
	)

	messagingInvokes = fx.Options(
		fx.Invoke(RunWatermillRouter),
	)
)
