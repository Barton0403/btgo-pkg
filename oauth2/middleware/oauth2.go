package middleware

import (
	"context"
	"github.com/Barton0403/btgo-pkg/common"
	"github.com/Barton0403/btgo-pkg/http"
	"github.com/Barton0403/btgo-pkg/http/middleware"
	"github.com/Barton0403/btgo-pkg/oauth2"
	"github.com/go-session/session/v3"
	"log"
	nethttp "net/http"
	"net/url"
)

func MakeOauth2TokenMiddleware(srv oauth2.Server) http.HandlerFunc {
	return func(ctx http.Context) {
		rctx := ctx.Request().Context()
		rctx = context.WithValue(rctx, common.DatabaseKey, middleware.GetDefaultDb(ctx))
		rctx = context.WithValue(rctx, common.CacheKey, middleware.GetDefaultCache(ctx))
		if err := srv.HandleTokenRequest(ctx.Writer(), ctx.Request().WithContext(rctx)); err != nil {
			log.Println("Handle Token Error:", err.Error())
		} else {
			//log.Println("Handle Token Success")
		}
	}
}

func MakeAuthorizeMiddleware(srv oauth2.Server) http.HandlerFunc {
	return func(ctx http.Context) {
		sessionStore, err := session.Start(ctx.Request().Context(), ctx.Writer(), ctx.Request())
		if err != nil {
			nethttp.Error(ctx.Writer(), err.Error(), nethttp.StatusInternalServerError)
			return
		}

		// 获取未登录前的请求参数
		var form url.Values
		if v, ok := sessionStore.Get("ReturnUri"); ok {
			form = v.(url.Values)
		}
		ctx.Request().Form = form

		sessionStore.Delete("ReturnUri")
		sessionStore.Save()

		rctx := ctx.Request().Context()
		rctx = context.WithValue(rctx, common.DatabaseKey, middleware.GetDefaultDb(ctx))
		rctx = context.WithValue(rctx, common.CacheKey, middleware.GetDefaultCache(ctx))
		err = srv.HandleAuthorizeRequest(ctx.Writer(), ctx.Request().WithContext(rctx))
		if err != nil {
			log.Println("Handle Authorize Error:", err.Error())
		}
	}
}
