package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis redis缓存
type Redis struct {
	client    *redis.Client
	keyPrefix string
}

// NewRedis 创建redis缓存
// 参数addr格式为 host:port 127.0.0.1:6379
func NewRedis(addr, password, keyPrefix string, db int) Cache {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		MinIdleConns: 10,
	})

	pong, err := client.Ping(context.Background()).Result()
	fmt.Println(pong, err)

	return &Redis{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// NewRedisWithClient 使用redis连接创建缓存实例
func NewRedisWithClient(redisClient *redis.Client, keyPrefix string) Cache {
	return &Redis{
		client:    redisClient,
		keyPrefix: keyPrefix,
	}
}

func (redis *Redis) key(key string) string {
	if redis.keyPrefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", redis.keyPrefix, key)
}

func (redis *Redis) IsExist(key string) bool {
	if v, err := redis.client.Exists(context.Background(), redis.key(key)).Result(); err == nil && v > 0 {
		return true
	}
	return false
}

func (redis *Redis) Get(key string) (string, error) {
	return redis.client.Get(context.Background(), redis.key(key)).Result()
}

func (redis *Redis) Put(key string, value string, expiration time.Duration) error {
	_, err := redis.client.Set(context.Background(), redis.key(key), value, expiration).Result()
	return err
}

func (redis *Redis) Delete(key string) error {
	_, err := redis.client.Del(context.Background(), redis.key(key)).Result()
	return err
}

// BatchDelete 批量删除
func (redis *Redis) BatchDelete(prefix string) error {
	var cursor uint64
	var keys []string
	var err error

	ctx := context.Background()

	for {
		keys, cursor, err = redis.client.Scan(ctx, cursor, prefix+"*", 10).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			err = redis.client.Del(ctx, key).Err()
			if err != nil {
				fmt.Printf("failed to delete key: %s\n", key)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return err
}
