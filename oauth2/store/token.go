package store

import (
	"barton.top/btgo/pkg/common"
	"barton.top/btgo/pkg/http/middleware"
	"context"
	"encoding/json"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/google/uuid"
	"time"
)

type TokenStore struct {
}

// Create create and store the new token information
func (t *TokenStore) Create(ctx context.Context, info oauth2.TokenInfo) error {
	cache := ctx.Value(common.CacheKey).(*middleware.Cache)
	nctx := ctx
	ct := time.Now()
	jv, err := json.Marshal(info)
	if err != nil {
		return err
	}

	if code := info.GetCode(); code != "" {
		return cache.Set(ctx, code, string(jv), info.GetCodeExpiresIn()).Err()
	}

	basicID := uuid.Must(uuid.NewRandom()).String()
	aexp := info.GetAccessExpiresIn()
	rexp := aexp
	if refresh := info.GetRefresh(); refresh != "" {
		if info.GetRefreshExpiresIn() == 0 {
			rexp = 0
			aexp = 0
		} else {
			rexp = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).Sub(ct)
			if aexp.Seconds() > rexp.Seconds() {
				aexp = rexp
			}
		}

		err = cache.Set(nctx, refresh, basicID, rexp).Err()
		if err != nil {
			return err
		}
	}

	err = cache.Set(nctx, basicID, string(jv), rexp).Err()
	if err != nil {
		return err
	}
	err = cache.Set(nctx, info.GetAccess(), basicID, aexp).Err()
	return err
}

// remove key
func (t *TokenStore) remove(ctx context.Context, key string) error {
	cache := ctx.Value(common.CacheKey).(*middleware.Cache)
	return cache.Del(ctx, key).Err()
}

// RemoveByCode use the authorization code to delete the token information
func (t *TokenStore) RemoveByCode(ctx context.Context, code string) error {
	return t.remove(ctx, code)
}

func (t *TokenStore) RemoveByAccess(ctx context.Context, access string) error {
	return t.remove(ctx, access)
}

func (t *TokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	return t.remove(ctx, refresh)
}

func (t *TokenStore) getData(ctx context.Context, key string) (oauth2.TokenInfo, error) {
	cache := ctx.Value(common.CacheKey).(*middleware.Cache)
	jv, err := cache.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var tm models.Token
	err = json.Unmarshal([]byte(jv), &tm)
	if err != nil {
		return nil, err
	}

	return &tm, nil
}

func (t *TokenStore) getBasicID(ctx context.Context, key string) (string, error) {
	cache := ctx.Value(common.CacheKey).(*middleware.Cache)

	basicID, err := cache.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return basicID, nil
}

func (t *TokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	return t.getData(ctx, code)
}

func (t *TokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	basicID, err := t.getBasicID(ctx, access)
	if err != nil {
		return nil, err
	}
	return t.getData(ctx, basicID)
}

func (t *TokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	basicID, err := t.getBasicID(ctx, refresh)
	if err != nil {
		return nil, err
	}
	return t.getData(ctx, basicID)
}
