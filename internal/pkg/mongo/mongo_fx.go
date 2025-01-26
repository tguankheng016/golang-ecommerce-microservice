package mongo

import "go.uber.org/fx"

var (
	// Module provided to fx
	Module = fx.Module(
		"mongo_fx",
		mongoProviders,
		mongoInvokes,
	)

	mongoProviders = fx.Options(
		fx.Provide(
			NewMongoDb,
		),
	)

	mongoInvokes = fx.Options(
		fx.Invoke(RunMongoDB),
	)
)
