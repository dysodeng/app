package mq

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/dysodeng/app/internal/infrastructure/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

func newRedisClient(cfg config.RedisItem) (redis.UniversalClient, error) {
	poolSize := cfg.Pool.PoolSize
	if poolSize <= 0 {
		poolSize = 100
	}
	minIdle := cfg.Pool.MinIdleConns
	if minIdle < 0 {
		minIdle = 0
	}
	switch cfg.Mode {
	case "", "standalone":
		host := cfg.Host
		if host == "" {
			host = "127.0.0.1"
		}
		port := cfg.Port
		if port == "" {
			port = "6379"
		}
		return redis.NewClient(&redis.Options{
			Addr:         net.JoinHostPort(host, port),
			Password:     cfg.Password,
			DB:           cfg.DB,
			PoolSize:     poolSize,
			MinIdleConns: minIdle,
			MaxRetries:   cfg.Pool.MaxRetries,
		}), nil
	case "cluster":
		if len(cfg.Cluster.Addrs) == 0 {
			return nil, fmt.Errorf("redis MQ cluster requires addresses")
		}
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        cfg.Cluster.Addrs,
			Password:     cfg.Cluster.Password,
			PoolSize:     poolSize,
			MinIdleConns: minIdle,
			MaxRetries:   cfg.Pool.MaxRetries,
		}), nil
	case "sentinel":
		if cfg.Sentinel.MasterName == "" || len(cfg.Sentinel.SentinelAddrs) == 0 {
			return nil, fmt.Errorf("redis MQ sentinel requires master name and addresses")
		}
		return redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       cfg.Sentinel.MasterName,
			SentinelAddrs:    cfg.Sentinel.SentinelAddrs,
			SentinelPassword: cfg.Sentinel.SentinelPassword,
			Password:         cfg.Sentinel.Password,
			DB:               cfg.Sentinel.DB,
			PoolSize:         poolSize,
			MinIdleConns:     minIdle,
			MaxRetries:       cfg.Pool.MaxRetries,
		}), nil
	default:
		return nil, fmt.Errorf("unsupported Redis MQ mode %q", cfg.Mode)
	}
}

func rabbitURL(cfg *config.Config) string {
	amqpCfg := cfg.MessageQueue.Amqp
	host := amqpCfg.Host
	if host == "" {
		host = "127.0.0.1"
	}
	port := amqpCfg.Port
	if port == "" {
		port = "5672"
	}
	vhost := amqpCfg.Vhost
	if vhost == "" {
		vhost = "/"
	}
	uri := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(amqpCfg.Username, amqpCfg.Password),
		Host:   net.JoinHostPort(host, port),
		Path:   "/" + vhost,
	}
	return uri.String()
}

func newRabbitConnection(cfg *config.Config) (*amqp.Connection, error) {
	conn, err := amqp.DialConfig(rabbitURL(cfg), amqp.Config{
		Heartbeat: 30 * time.Second,
		Dial:      amqp.DefaultDial(10 * time.Second),
	})
	if err != nil {
		return nil, fmt.Errorf("connect RabbitMQ: %w", err)
	}
	return conn, nil
}
