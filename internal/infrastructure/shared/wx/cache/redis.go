package cache

import (
	"context"
	"time"

	"github.com/dysodeng/wx/support/cache"
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

func (redis *Redis) IsExist(key string) bool {
	if v, err := redis.client.Exists(context.Background(), key).Result(); err == nil && v > 0 {
		return true
	}
	return false
}

func (redis *Redis) Get(key string) (string, error) {
	return redis.client.Get(context.Background(), key).Result()
}

func (redis *Redis) Put(key string, value string, expiration time.Duration) error {
	_, err := redis.client.Set(context.Background(), key, value, expiration).Result()
	return err
}

func (redis *Redis) Delete(key string) error {
	_, err := redis.client.Del(context.Background(), key).Result()
	return err
}

func (redis *Redis) ClearAll() error {
	_, err := redis.client.FlushDB(context.Background()).Result()
	return err
}
