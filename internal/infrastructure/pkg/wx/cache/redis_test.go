package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	corecache "github.com/goairix/wx/v2/core/cache"
	"github.com/redis/go-redis/v9"
)

type fakeRedisClient struct {
	redis.UniversalClient
	get func(context.Context, string) *redis.StringCmd
	set func(context.Context, string, interface{}, time.Duration) *redis.StatusCmd
	del func(context.Context, ...string) *redis.IntCmd
}

var _ corecache.Cache = (*Redis)(nil)

func (f fakeRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return f.get(ctx, key)
}

func (f fakeRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return f.set(ctx, key, value, expiration)
}

func (f fakeRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return f.del(ctx, keys...)
}

func TestRedisCacheGetDistinguishesMissFromFailure(t *testing.T) {
	ctx := context.Background()
	backendFailure := errors.New("redis unavailable")
	responses := map[string]struct {
		value string
		err   error
	}{
		"found":   {value: "token"},
		"missing": {err: redis.Nil},
		"failed":  {err: backendFailure},
	}
	cache := NewRedis(fakeRedisClient{
		get: func(gotCtx context.Context, key string) *redis.StringCmd {
			if gotCtx != ctx {
				t.Fatal("Get did not use the caller context")
			}
			response := responses[key]
			return redis.NewStringResult(response.value, response.err)
		},
	})

	for _, test := range []struct {
		key   string
		value string
		found bool
		err   error
	}{
		{key: "found", value: "token", found: true},
		{key: "missing"},
		{key: "failed", err: backendFailure},
	} {
		value, found, err := cache.Get(ctx, test.key)
		if value != test.value || found != test.found || !errors.Is(err, test.err) {
			t.Errorf("Get(%q) = (%q, %v, %v), want (%q, %v, %v)", test.key, value, found, err, test.value, test.found, test.err)
		}
	}
}

func TestRedisCachePutAndDeletePassContextAndTTL(t *testing.T) {
	ctx := context.Background()
	ttl := 5 * time.Minute
	var put, deleted bool
	cache := NewRedis(fakeRedisClient{
		set: func(gotCtx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
			if gotCtx != ctx || key != "key" || value != "value" || expiration != ttl {
				t.Errorf("unexpected Set arguments: %v, %q, %v, %v", gotCtx, key, value, expiration)
			}
			put = true
			return redis.NewStatusResult("OK", nil)
		},
		del: func(gotCtx context.Context, keys ...string) *redis.IntCmd {
			if gotCtx != ctx || len(keys) != 1 || keys[0] != "key" {
				t.Errorf("unexpected Del arguments: %v, %v", gotCtx, keys)
			}
			deleted = true
			return redis.NewIntResult(1, nil)
		},
	})

	if err := cache.Put(ctx, "key", "value", ttl); err != nil {
		t.Fatal(err)
	}
	if err := cache.Delete(ctx, "key"); err != nil {
		t.Fatal(err)
	}
	if !put || !deleted {
		t.Fatalf("expected Put and Delete to call Redis, got put=%v deleted=%v", put, deleted)
	}
}
