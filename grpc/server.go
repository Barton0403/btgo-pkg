package grpc

import (
	"context"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
)

type Service struct {
	Desc *grpc.ServiceDesc
	Impl interface{}
}

var DefaultServices []Service

func Register(desc *grpc.ServiceDesc, impl interface{}) {
	DefaultServices = append(DefaultServices, Service{
		Desc: desc,
		Impl: impl,
	})
}

var DefaultMiddlewares []grpc.UnaryServerInterceptor

func Use(h grpc.UnaryServerInterceptor) {
	DefaultMiddlewares = append(DefaultMiddlewares, h)
}

type HandlerFunc func(ctx context.Context, in interface{}) (resp interface{}, err error)

type HandlerOption func(handler *Handler)

// ServerBefore functions are executed on the gRPC request object before the
// request is decoded.
func ServerBefore(before ...kitgrpc.ServerRequestFunc) HandlerOption {
	return func(s *Handler) { s.before = append(s.before, before...) }
}

// ServerAfter functions are executed on the gRPC response writer after the
// endpoint is invoked, but before anything is written to the client.
func ServerAfter(after ...kitgrpc.ServerResponseFunc) HandlerOption {
	return func(s *Handler) { s.after = append(s.after, after...) }
}

// ServerFinalizer is executed at the end of every gRPC request.
// By default, no finalizer is registered.
func ServerFinalizer(f ...kitgrpc.ServerFinalizerFunc) HandlerOption {
	return func(s *Handler) { s.finalizer = append(s.finalizer, f...) }
}

type Handler struct {
	f         HandlerFunc
	before    []kitgrpc.ServerRequestFunc
	after     []kitgrpc.ServerResponseFunc
	finalizer []kitgrpc.ServerFinalizerFunc
}

func (s *Handler) Serve(ctx context.Context, in interface{}) (resp interface{}, err error) {
	// Retrieve gRPC metadata.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.MD{}
	}

	if len(s.finalizer) > 0 {
		defer func() {
			for _, f := range s.finalizer {
				f(ctx, err)
			}
		}()
	}

	for _, f := range s.before {
		ctx = f(ctx, md)
	}

	resp, err = s.f(ctx, in)
	if err != nil {
		return nil, err
	}

	var mdHeader, mdTrailer metadata.MD
	for _, f := range s.after {
		ctx = f(ctx, &mdHeader, &mdTrailer)
	}

	return resp, nil
}

func NewHandler(f HandlerFunc, options ...HandlerOption) *Handler {
	h := &Handler{
		f: f,
	}

	for _, option := range options {
		option(h)
	}

	return h
}

func Run(addr string) error {
	lis, e := net.Listen("tcp", addr)
	if e != nil {
		return e
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(DefaultMiddlewares...),
		grpc.UnaryInterceptor(kitgrpc.Interceptor),
	)

	// 服务注册
	for _, r := range DefaultServices {
		s.RegisterService(r.Desc, r.Impl)
	}

	log.Printf("grpc server listening at %v", lis.Addr())
	return s.Serve(lis)
}
