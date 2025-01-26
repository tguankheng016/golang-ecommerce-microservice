package grpc

import "go.uber.org/fx"

var (
	// Module provided to fx
	Module = fx.Module(
		"grpcserverfx",
		grpcServerProviders,
		grpcInvokes,
	)

	grpcServerProviders = fx.Options(
		fx.Provide(
			NewGrpcServer,
			NewGrpcClientFactory,
		),
	)

	grpcInvokes = fx.Options(
		fx.Invoke(RunServers),
	)
)
