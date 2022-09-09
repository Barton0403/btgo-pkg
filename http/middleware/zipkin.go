package middleware

import (
	"github.com/Barton0403/btgo-pkg/http"
	"github.com/openzipkin/zipkin-go"
	zipkinmiddlewarehttp "github.com/openzipkin/zipkin-go/middleware/http"
	nethttp "net/http"
)

func ZipKinMiddleware(tracer *zipkin.Tracer, spanName string) http.HandlerFunc {
	return func(ctx http.Context) {
		mid := zipkinmiddlewarehttp.NewServerMiddleware(tracer,
			zipkinmiddlewarehttp.SpanName(spanName),
			zipkinmiddlewarehttp.TagResponseSize(true),
		)
		h := mid(nethttp.HandlerFunc(func(writer nethttp.ResponseWriter, request *nethttp.Request) {}))
		h.ServeHTTP(ctx.Writer(), ctx.Request())
	}
}
