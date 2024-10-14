package http

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
)

func NewContext() context.Context {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			<-ctx.Done()
			logger.Logger.Info("context is canceled!")
			cancel()
			return
		}
	}()

	return ctx
}
