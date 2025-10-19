package driver

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/dysodeng/app/internal/infrastructure/config"
	infraRedis "github.com/dysodeng/app/internal/infrastructure/shared/redis"
)

// Redis redis驱动
type Redis struct {
	client    redis.UniversalClient
	keyPrefix string
}

func NewRedisCache() *Redis {
	return &Redis{
		client:    infraRedis.CacheClient(),
		keyPrefix: config.GlobalConfig.Redis.Cache.KeyPrefix,
	}
}

func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	v, err := r.client.Exists(ctx, r.key(key)).Result()
	return v > 0 && err == nil, err
}

func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	s, err := r.client.Get(ctx, r.key(key)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	return s, err
}

func (r *Redis) Set(ctx context.Context, key string, val []byte, ttl time.Duration) error {
	return r.client.Set(ctx, r.key(key), val, ttl).Err()
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, r.key(key)).Err()
}

func (r *Redis) ScanDeleteByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	for {
		keys, next, err := r.client.Scan(ctx, cursor, r.key(prefix)+"*", 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := r.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return nil
}

func (r *Redis) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, r.key(key)).Result()
}

func (r *Redis) key(key string) string {
	return fmt.Sprintf("%s:%s", r.keyPrefix, key)
}
