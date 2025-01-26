package logging

import "go.uber.org/fx"

var (
	// Module provided to fx
	Module = fx.Module(
		"loggerfx",
		loggerProvider,
		loggerInvoke,
	)

	loggerProvider = fx.Options(
		fx.Provide(
			InitLogger,
		),
	)

	loggerInvoke = fx.Options(
		fx.Invoke(RunLogger),
	)
)
