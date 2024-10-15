package middlewares

import (
	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logger"
	"go.uber.org/zap"
)

func SetupLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogUserAgent: true,
		LogRequestID: true,
		LogError:     true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			fields := []zap.Field{
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
				zap.Duration("latency", v.Latency),
				zap.String("remote_ip", v.RemoteIP),
				zap.String("user_agent", v.UserAgent),
				zap.String("request_id", v.RequestID),
			}

			if v.Error != nil {
				logger.Logger.Error("echo-error-request", append(fields, zap.Error(v.Error))...)
			} else {
				logger.Logger.Info("echo-success-request", fields...)
			}

			return nil
		},
	})
}
