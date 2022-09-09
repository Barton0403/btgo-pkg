package engine

import (
	"barton.top/btgo/pkg/common"
	"barton.top/btgo/pkg/http"
	"github.com/gin-gonic/gin"
	nethttp "net/http"
)

func GinRun(addr ...string) error {
	server := gin.New()
	// 中间件
	for _, m := range http.DefaultMiddlewares {
		server.Use(MakeGinHandlerFunc(m))
	}
	// 路由
	for _, r := range http.DefaultRoutes {
		var hs gin.HandlersChain
		for _, h := range r.Handlers {
			hs = append(hs, MakeGinHandlerFunc(h))
		}
		server.Handle(r.Method, r.Path, hs...)
	}
	return server.Run(addr...)
}

type GinContext struct {
	*gin.Context
}

func (c *GinContext) Request() *nethttp.Request {
	return c.Context.Request
}

func (c *GinContext) Writer() nethttp.ResponseWriter {
	return c.Context.Writer
}

func (c *GinContext) Error(code int, msg string) {
	c.Context.AbortWithStatusJSON(code, gin.H{"error": msg})
}

func (c *GinContext) Parent() interface{} {
	return c.Context
}

func (c *GinContext) GetUnionId() string {
	return c.Request().Header.Get(common.HeaderUnionId)
}

func MakeGinHandlerFunc(handler http.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handler(&GinContext{
			Context: ctx,
		})
	}
}
