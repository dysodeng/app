package driver

import (
	"context"
	"fmt"
	"time"

	"github.com/dysodeng/app/internal/infrastructure/persistence/cache/contract"
	"github.com/dysodeng/app/internal/pkg/redis"
)

type redisCache struct {
	client    redis.Client
	keyPrefix string
}

func NewRedisCache(client redis.Client, keyPrefix string) contract.Cache {
	return &redisCache{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

func (r *redisCache) key(key string) string {
	if r.keyPrefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", r.keyPrefix, key)
}

func (r *redisCache) IsExist(ctx context.Context, key string) bool {
	if v, err := r.client.Exists(ctx, r.key(key)).Result(); err == nil && v > 0 {
		return true
	}
	return false
}

func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, r.key(key)).Result()
}

func (r *redisCache) Put(ctx context.Context, key string, value string, expiration time.Duration) error {
	_, err := r.client.Set(ctx, r.key(key), value, expiration).Result()
	return err
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, r.key(key)).Result()
	return err
}

func (r *redisCache) BatchDelete(ctx context.Context, prefix string) error {
	var cursor uint64
	var keys []string
	var err error

	for {
		keys, cursor, err = r.client.Scan(ctx, cursor, r.key(prefix)+"*", 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			err = r.client.Del(ctx, keys...).Err()
			if err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (r *redisCache) GetWithTTL(ctx context.Context, key string) (string, time.Duration, error) {
	value, err := r.client.Get(ctx, r.key(key)).Result()
	if err != nil {
		return "", 0, err
	}
	ttl, err := r.client.TTL(ctx, r.key(key)).Result()
	return value, ttl, err
}

func (r *redisCache) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	return r.client.IncrBy(ctx, r.key(key), delta).Result()
}

func (r *redisCache) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	return r.client.DecrBy(ctx, r.key(key), delta).Result()
}

func (r *redisCache) SetNX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, r.key(key), value, expiration).Result()
}

func (r *redisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, r.key(key), expiration).Err()
}

func (r *redisCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	return r.client.Keys(ctx, r.key(pattern)).Result()
}

func (r *redisCache) FlushAll(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}

func (r *redisCache) Close() error {
	return r.client.Close()
}
