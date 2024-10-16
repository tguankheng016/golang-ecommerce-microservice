package http

import (
	echoServer "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http/echo"
	"go.uber.org/fx"
)

var (
	// Module provided to fx
	Module = fx.Module(
		"httpfx",
		httpProviders,
	)

	httpProviders = fx.Options(
		fx.Provide(
			NewContext,
			echoServer.NewEchoServer,
		),
	)
)
