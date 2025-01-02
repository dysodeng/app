package redis

import (
	"context"
	"log"

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
	if err != nil {
		log.Fatalf("failed to connect redis %+v", err)
	}
	log.Printf("redis state: %s", pong)
	log.Println("redis connection successful")
}

func Initialize() {}

// Client 获取redis实例
func Client() *redis.Client {
	return redisPoolClient
}

func Close() {
	err := redisPoolClient.Close()
	if err != nil {
		log.Printf("failed to close redis connection: %+v", err)
		return
	}
	log.Println("redis connection closed")
}

// Key 构建安全缓存key
func Key(key string) string {
	prefix := config.Redis.Main.KeyPrefix
	if prefix != "" {
		key = config.Redis.Main.KeyPrefix + ":" + key
	}
	return key
}
