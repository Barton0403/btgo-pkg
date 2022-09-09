package middleware

import (
	"barton.top/btgo/pkg/common"
	"barton.top/btgo/pkg/http"
	"barton.top/btgo/pkg/http/middleware"
	"barton.top/btgo/pkg/oauth2"
	"context"
	"log"
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
