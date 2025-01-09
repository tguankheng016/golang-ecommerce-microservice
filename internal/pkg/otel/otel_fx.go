package otel

import "go.uber.org/fx"

var (
	// Module provided to fx
	Module = fx.Module(
		"otel_fx",
		otelProviders,
		otelInvokes,
	)

	otelProviders = fx.Options(
		fx.Provide(
			NewJaegerExporter,
			NewTraceProvider,
		),
	)

	otelInvokes = fx.Options(
		fx.Invoke(RunTracing),
	)
)
