package middleware

import (
	"barton.top/btgo/pkg/common"
	"context"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

type Cache = redis.Client

func CacheUnaryServerInterceptor(addr string, password string, db int) grpc.UnaryServerInterceptor {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
		//PoolSize: 50,
	})

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(context.WithValue(ctx, common.CacheKey, rdb), req)
	}
}

func GetDefaultCache(ctx context.Context) *redis.Client {
	if v := ctx.Value(common.CacheKey); v != nil {
		return v.(*redis.Client)
	}
	panic("no cache found")
}
