package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logging"
	"go.uber.org/zap"
)

type ResponseWriter struct {
	middleware.WrapResponseWriter
	statusCode int
	body       string
}

func NewResponseWriter(w middleware.WrapResponseWriter) *ResponseWriter {
	return &ResponseWriter{WrapResponseWriter: w, statusCode: http.StatusOK}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.WrapResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseWriter) Write(p []byte) (n int, err error) {
	if rw.statusCode >= 400 {
		// Capture the error message in the body for error responses
		rw.body = string(p)
	}
	return rw.WrapResponseWriter.Write(p)
}

func SetupLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := NewResponseWriter(middleware.NewWrapResponseWriter(w, r.ProtoMajor))

			t1 := time.Now()
			defer func() {
				fields := []zap.Field{
					zap.String("proto", r.Proto),
					zap.String("path", r.URL.Path),
					zap.Duration("lat", time.Since(t1)),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
					zap.String("reqId", middleware.GetReqID(r.Context())),
				}

				if ww.Status() >= 400 {
					fields = append(fields, zap.String("error_message", ww.body))
					logging.Logger.Error("Error", fields...)
				} else {
					logging.Logger.Info("Served", fields...)
				}
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
