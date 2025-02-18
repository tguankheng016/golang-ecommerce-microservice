package messaging

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"

	wotelfloss "github.com/dentech-floss/watermill-opentelemetry-go-extra/pkg/opentelemetry"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logging"
	wotel "github.com/voi-oss/watermill-opentelemetry/pkg/opentelemetry"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewWatermillRouter(logger watermill.LoggerAdapter) (*message.Router, error) {
	router, err := message.NewRouter(message.RouterConfig{}, logger)

	router.AddMiddleware(
		// CorrelationID will copy the correlation id from the incoming message's metadata to the produced messages
		middleware.CorrelationID,

		// The handler function is retried if it returns an error.
		// After MaxRetries, the message is Nacked and it's up to the PubSub to resend it.
		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Millisecond * 100,
			Logger:          logger,
		}.Middleware,

		// Recoverer handles panics from handlers.
		// In this case, it passes them as errors to the Retry middleware.
		middleware.Recoverer,
	)

	router.AddMiddleware(wotelfloss.ExtractRemoteParentSpanContext())
	router.AddMiddleware(wotel.Trace())

	return router, err
}

func RunWatermillRouter(lc fx.Lifecycle, router *message.Router, subscriber message.Subscriber, ctx context.Context) error {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := router.Run(ctx); err != nil {
					logging.Logger.Fatal("error running watermill router", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(_ context.Context) error {
			logging.Logger.Info("closing watermill subscriber...")

			if err := subscriber.Close(); err != nil {
				logging.Logger.Info("error when closing watermill subscriber...", zap.Error(err))
			}

			logging.Logger.Info("watermill router shutdown gracefully")

			return nil
		},
	})

	return nil
}
