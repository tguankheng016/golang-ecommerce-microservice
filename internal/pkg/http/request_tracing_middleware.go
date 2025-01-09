package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	otelchimetric "github.com/riandyrn/otelchi/metric"
)

func SetupTracing(serviceName string, router *chi.Mux) []func(next http.Handler) http.Handler {
	// define base config for metric middlewares
	baseCfg := otelchimetric.NewBaseConfig(serviceName)

	res := make([]func(next http.Handler) http.Handler, 0)

	res = append(res, otelchi.Middleware(serviceName, otelchi.WithChiRoutes(router)))
	res = append(res, otelchimetric.NewRequestDurationMillis(baseCfg))
	res = append(res, otelchimetric.NewRequestInFlight(baseCfg))
	res = append(res, otelchimetric.NewResponseSizeBytes(baseCfg))

	return res
}
