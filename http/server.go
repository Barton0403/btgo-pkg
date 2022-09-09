package http

import (
	nethttp "net/http"
)

type Context interface {
	JSON(int, interface{})
	Next()
	Abort()
	Set(key string, value interface{})
	Get(key string) (value interface{}, exists bool)
	PostForm(key string) string
	Query(key string) string
	GetHeader(key string) string
	Error(code int, msg string)
	GetUnionId() string

	Request() *nethttp.Request
	Writer() nethttp.ResponseWriter
	Parent() interface{}
}

type HandlerFunc func(Context)

type Route struct {
	Method   string
	Path     string
	Handlers []HandlerFunc
}

var DefaultRoutes []Route

func Handle(m string, p string, h ...HandlerFunc) {
	DefaultRoutes = append(DefaultRoutes, Route{
		Method:   m,
		Path:     p,
		Handlers: h,
	})
}

var DefaultMiddlewares []HandlerFunc

func Use(h HandlerFunc) {
	DefaultMiddlewares = append(DefaultMiddlewares, h)
}
