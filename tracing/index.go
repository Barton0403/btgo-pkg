package tracing

import (
	"barton.top/btgo/pkg/grpc"
	"context"
	"github.com/opentracing/opentracing-go"
	zipkintracer "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"strconv"

	zipkin "github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

// GRPCServerTrace enables native Zipkin tracing of a Go kit gRPC transport
// Server.
//
// Go kit creates gRPC transport servers per gRPC method. This middleware can be
// set-up individually by adding the method name for each of the Go kit method
// servers using the Name() TracerOption.
// If wanting to use the gRPC FullMethod (/service/method) as Span name you can
// create a global server tracer omitting the Name() TracerOption, which you can
// then feed to each Go kit method server. For this to work you will need to
// wire the Go kit gRPC Interceptor too.
// If instrumenting a service to external (not on your platform) clients, you
// will probably want to disallow propagation of a client SpanContext using
// the AllowPropagation TracerOption and setting it to false.
func GRPCServerTrace(tracer *zipkin.Tracer, options ...TracerOption) grpc.HandlerOption {
	config := tracerOptions{
		tags:      make(map[string]string),
		name:      "",
		logger:    log.NewNopLogger(),
		propagate: true,
	}

	for _, option := range options {
		option(&config)
	}

	serverBefore := grpc.ServerBefore(
		func(ctx context.Context, md metadata.MD) context.Context {
			var (
				spanContext model.SpanContext
				name        string
				tags        = make(map[string]string)
			)

			rpcMethod, ok := ctx.Value(kitgrpc.ContextKeyRequestMethod).(string)
			if !ok {
				config.logger.Log("err", "unable to retrieve method name: missing gRPC interceptor hook")
			} else {
				tags["grpc.method"] = rpcMethod
			}

			if config.name != "" {
				name = config.name
			} else {
				name = rpcMethod
			}

			if config.propagate {
				spanContext = tracer.Extract(b3.ExtractGRPC(&md))
				if spanContext.Err != nil {
					config.logger.Log("err", spanContext.Err)
				}
			}

			span := tracer.StartSpan(
				name,
				zipkin.Kind(model.Server),
				zipkin.Tags(config.tags),
				zipkin.Tags(tags),
				zipkin.Parent(spanContext),
				zipkin.FlushOnFinish(false),
			)

			return zipkin.NewContext(ctx, span)
		},
	)

	serverAfter := grpc.ServerAfter(
		func(ctx context.Context, _ *metadata.MD, _ *metadata.MD) context.Context {
			if span := zipkin.SpanFromContext(ctx); span != nil {
				span.Finish()
			}

			return ctx
		},
	)

	serverFinalizer := grpc.ServerFinalizer(
		func(ctx context.Context, err error) {
			if span := zipkin.SpanFromContext(ctx); span != nil {
				if err != nil {
					if status, ok := status.FromError(err); ok {
						statusCode := strconv.FormatUint(uint64(status.Code()), 10)
						zipkin.TagGRPCStatusCode.Set(span, statusCode)
						zipkin.TagError.Set(span, status.Message())
					} else {
						zipkin.TagError.Set(span, err.Error())
					}
				}

				// calling span.Finish() a second time is a noop, if we didn't get to
				// ServerAfter we can at least time the early bail out by calling it
				// here.
				span.Finish()
				// send span to the Reporter
				span.Flush()
			}
		},
	)

	return func(s *grpc.Handler) {
		serverBefore(s)
		serverAfter(s)
		serverFinalizer(s)
	}
}

var DefaultTracer *zipkin.Tracer
var DefaultOpenTracer opentracing.Tracer

func TracerInit(addr string, serviceName string, port string) error {
	reporter := zipkinhttp.NewReporter(addr, zipkinhttp.MaxBacklog(10000))

	zEP, _ := zipkin.NewEndpoint(serviceName, port)
	var err error
	DefaultTracer, err = zipkin.NewTracer(
		reporter, zipkin.WithLocalEndpoint(zEP),
	)
	if err != nil {
		return err
	}

	DefaultOpenTracer = zipkintracer.Wrap(DefaultTracer)

	return nil
}
