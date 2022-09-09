package middleware

import (
	"github.com/Barton0403/btgo-pkg/common"
	"github.com/Barton0403/btgo-pkg/http"
	"github.com/go-redis/redis/v8"
)

type Cache = redis.Client

func CacheMiddleware(addr string, password string, db int) http.HandlerFunc {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
		//PoolSize: 50,
	})

	return func(c http.Context) {
		c.Set(common.CacheKey, rdb)
	}
}

func GetDefaultCache(ctx http.Context) *redis.Client {
	if v, e := ctx.Get(common.CacheKey); e {
		return v.(*redis.Client)
	}
	panic("no cache found")
}
