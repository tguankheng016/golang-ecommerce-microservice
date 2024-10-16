package redis

import "go.uber.org/fx"

var (
	// Module provided to fx
	Module = fx.Module(
		"redisfx",
		redisProviders,
		redisInvokes,
	)

	redisProviders = fx.Options(
		fx.Provide(
			NewRedisClient,
			NewRedisUniversalClient,
		),
	)

	redisInvokes = fx.Options(
		fx.Invoke(RegisterRedisServer),
	)
)
