package cache

import (
	"context"
	"errors"
	"time"

	"github.com/goairix/wx/v2/core/cache"
	"github.com/redis/go-redis/v9"
)

// Redis redis缓存
type Redis struct {
	client redis.UniversalClient
}

// NewRedis 创建redis缓存
func NewRedis(client redis.UniversalClient) cache.Cache {
	return &Redis{
		client: client,
	}
}

func (r *Redis) Get(ctx context.Context, key string) (string, bool, error) {
	value, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return value, true, nil
}

func (r *Redis) Put(ctx context.Context, key string, value string, expiration time.Duration) error {
	_, err := r.client.Set(ctx, key, value, expiration).Result()
	return err
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	return err
}
