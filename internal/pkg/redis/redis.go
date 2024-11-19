package redis

import (
	"context"
	"fmt"

	"github.com/dysodeng/app/internal/config"
	"github.com/go-redis/redis/v8"
)

var redisPoolClient *redis.Client

func init() {
	addr := config.Redis.Main.Host + ":" + config.Redis.Main.Port
	redisPoolClient = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     config.Redis.Main.Password,
		DB:           config.Redis.Main.DB,
		MinIdleConns: 10,
	})

	pong, err := redisPoolClient.Ping(context.Background()).Result()
	fmt.Println(pong, err)
}

func Initialize() {}

// Client 获取redis实例
func Client() *redis.Client {
	return redisPoolClient
}

// Key 构建安全缓存key
func Key(key string) string {
	prefix := config.Redis.Main.KeyPrefix
	if prefix != "" {
		key = config.Redis.Main.KeyPrefix + ":" + key
	}
	return key
}
