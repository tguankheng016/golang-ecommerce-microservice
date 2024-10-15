package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	echoServer "github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/http/echo"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/services/identity_service/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func RunServers(lc fx.Lifecycle, e *echo.Echo, ctx context.Context, cfg *config.Config) error {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := echoServer.RunHttpServer(ctx, e, cfg.EchoOptions); !errors.Is(err, http.ErrServerClosed) {
					logger.Logger.Fatal("error running http server", zap.Error(err))
				}
			}()

			e.GET("/", func(c echo.Context) error {
				return c.String(http.StatusOK, "Golang ECommerce Identity Service")
			})

			return nil
		},
		OnStop: func(_ context.Context) error {
			logger.Logger.Info("all servers shutdown gracefully...")
			return nil
		},
	})

	return nil
}
