package middleware

import (
	"context"
	"github.com/Barton0403/btgo-pkg/common"
	"github.com/Barton0403/btgo-pkg/http"
	"github.com/Barton0403/btgo-pkg/http/middleware"
	"github.com/Barton0403/btgo-pkg/oauth2"
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
