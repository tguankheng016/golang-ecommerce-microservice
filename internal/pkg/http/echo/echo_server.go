package echo

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	MaxHeaderBytes = 1 << 20
	ReadTimeout    = 15 * time.Second
	WriteTimeout   = 15 * time.Second
)

func NewEchoServer() *echo.Echo {
	e := echo.New()
	return e
}

func RunHttpServer(ctx context.Context, echo *echo.Echo, cfg *EchoOptions) error {
	echo.Server.ReadTimeout = ReadTimeout
	echo.Server.WriteTimeout = WriteTimeout
	echo.Server.MaxHeaderBytes = MaxHeaderBytes

	go func() {
		<-ctx.Done()
		logger.Logger.Info("shutting down Http PORT: " + cfg.Port)
		err := echo.Shutdown(ctx)
		if err != nil {
			logger.Logger.Error("(Shutdown) err: ", zap.Error(err))
			return
		}
		logger.Logger.Info("echo server exited properly")
	}()

	err := echo.Start(cfg.Port)

	return err
}

func RunServers(lc fx.Lifecycle, e *echo.Echo, ctx context.Context, cfg *EchoOptions) error {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := RunHttpServer(ctx, e, cfg); !errors.Is(err, http.ErrServerClosed) {
					logger.Logger.Fatal("error running http server", zap.Error(err))
				}
			}()

			e.GET("/", func(c echo.Context) error {
				return c.String(http.StatusOK, cfg.Name)
			})

			return nil
		},
		OnStop: func(_ context.Context) error {
			logger.Logger.Info("all echo servers shutdown gracefully...")
			return nil
		},
	})

	return nil
}
