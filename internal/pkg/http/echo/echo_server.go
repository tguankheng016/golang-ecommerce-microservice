package echo

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
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
		logger.Logger.Info("server exited properly")
	}()

	err := echo.Start(cfg.Port)

	return err
}
