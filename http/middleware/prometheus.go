package middleware

import (
	"barton.top/btgo/pkg/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func PrometheusMiddleware() http.HandlerFunc {
	return func(ctx http.Context) {
		h := promhttp.Handler()
		h.ServeHTTP(ctx.Writer(), ctx.Request())
	}
}
