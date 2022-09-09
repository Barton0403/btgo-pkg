package limit

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// Black 黑名单
type black struct {
	cache *redis.Client
	key   string
}

func (b *black) Allow(ctx context.Context, v string) bool {
	r, err := b.cache.SIsMember(ctx, b.key, v).Result()
	if err != nil {
		return true
	}

	if r {
		return false
	}

	return true
}

func (b *black) Add(ctx context.Context, v string) error {
	_, e := b.cache.SAdd(ctx, b.key, v).Result()
	if e != nil {
		return e
	}

	return nil
}

func (b *black) Rem(ctx context.Context, v string) error {
	_, e := b.cache.SRem(ctx, b.key, v).Result()
	if e != nil {
		return e
	}

	return nil
}

func NewBlack(cache *redis.Client, key string) *black {
	return &black{
		cache: cache,
		key:   key,
	}
}
