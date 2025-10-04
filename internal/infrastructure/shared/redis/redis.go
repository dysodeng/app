package redis

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/dysodeng/app/internal/infrastructure/config"
)

// Client redis客户端
type Client interface {
	redis.UniversalClient
	Close() error
	Ping(ctx context.Context) *redis.StatusCmd
}

var (
	mainClient  Client // 主redis连接
	cacheClient Client // redis缓存连接
)

func createRedisClient(cfg config.RedisItem) (Client, error) {
	switch strings.ToLower(cfg.Mode) {
	case "cluster":
		return createClusterClient(cfg)
	case "sentinel":
		return createSentinelClient(cfg)
	case "standalone", "":
		return createStandaloneClient(cfg)
	default:
		return nil, fmt.Errorf("unsupported redis mode: %s", cfg.Mode)
	}
}

func createStandaloneClient(cfg config.RedisItem) (Client, error) {
	addr := cfg.Host + ":" + cfg.Port
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MinIdleConns: cfg.Pool.MinIdleConns,
		MaxRetries:   cfg.Pool.MaxRetries,
		PoolSize:     cfg.Pool.PoolSize,
	})
	return client, nil
}

func createClusterClient(cfg config.RedisItem) (Client, error) {
	if len(cfg.Cluster.Addrs) == 0 {
		return nil, fmt.Errorf("cluster addrs cannot be empty")
	}

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        cfg.Cluster.Addrs,
		Password:     cfg.Cluster.Password,
		MinIdleConns: cfg.Pool.MinIdleConns,
		MaxRetries:   cfg.Pool.MaxRetries,
		PoolSize:     cfg.Pool.PoolSize,
	})
	return client, nil
}

func createSentinelClient(cfg config.RedisItem) (Client, error) {
	if cfg.Sentinel.MasterName == "" {
		return nil, fmt.Errorf("sentinel master name cannot be empty")
	}
	if len(cfg.Sentinel.SentinelAddrs) == 0 {
		return nil, fmt.Errorf("sentinel addrs cannot be empty")
	}

	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       cfg.Sentinel.MasterName,
		SentinelAddrs:    cfg.Sentinel.SentinelAddrs,
		Password:         cfg.Sentinel.Password,
		SentinelPassword: cfg.Sentinel.SentinelPassword,
		DB:               cfg.Sentinel.DB,
		MinIdleConns:     cfg.Pool.MinIdleConns,
		MaxRetries:       cfg.Pool.MaxRetries,
		PoolSize:         cfg.Pool.PoolSize,
	})
	return client, nil
}

func testConnection(client Client, name string) error {
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("failed to ping %s redis: %w", name, err)
	}
	log.Printf("%s redis state: %s", name, pong)
	return nil
}

func Initialize(cfg *config.Config) (Client, error) {
	var err error

	// 初始化主Redis客户端
	mainClient, err = createRedisClient(cfg.Redis.Main)
	if err != nil {
		log.Fatalf("failed to create main redis client: %+v", err)
	}

	// 初始化缓存Redis客户端
	cacheClient, err = createRedisClient(cfg.Redis.Cache)
	if err != nil {
		log.Fatalf("failed to create cache redis client: %+v", err)
	}

	// 测试连接
	if err = testConnection(mainClient, "main"); err != nil {
		log.Fatalf("main redis connection test failed: %+v", err)
	}

	if err = testConnection(cacheClient, "cache"); err != nil {
		log.Fatalf("cache redis connection test failed: %+v", err)
	}

	log.Println("redis connections initialized successfully")

	return mainClient, nil
}

// MainClient 获取主redis实例
func MainClient() Client {
	return mainClient
}

// CacheClient 获取缓存redis实例
func CacheClient() Client {
	return cacheClient
}

func Close() {
	if err := mainClient.Close(); err != nil {
		log.Printf("failed to close main redis connection: %+v", err)
	}
	if err := cacheClient.Close(); err != nil {
		log.Printf("failed to close cache redis connection: %+v", err)
	}
	log.Println("redis connections closed")
}

// MainKey 构建安全key
func MainKey(key string) string {
	prefix := config.GlobalConfig.Redis.Main.KeyPrefix
	if prefix != "" {
		key = config.GlobalConfig.Redis.Main.KeyPrefix + ":" + key
	}
	return key
}

// CacheKey 构建缓存key
func CacheKey(key string) string {
	prefix := config.GlobalConfig.Redis.Cache.KeyPrefix
	if prefix != "" {
		key = config.GlobalConfig.Redis.Cache.KeyPrefix + ":" + key
	}
	return key
}
