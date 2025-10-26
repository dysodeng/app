package mq

import (
	"strconv"
	"time"

	"github.com/dysodeng/app/internal/infrastructure/config"
	mqConfig "github.com/dysodeng/mq/config"
)

func createRedisConfig(cfg *config.Config) mqConfig.RedisConfig {
	redisCfg := mqConfig.RedisConfig{
		PoolSize:           200, // 连接池
		MinIdleConns:       50,  // 最小空闲连接
		MaxConnAge:         time.Hour,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        5 * time.Minute,
		IdleCheckFrequency: time.Minute,
		MaxRetries:         3,
		MinRetryBackoff:    8 * time.Millisecond,
		MaxRetryBackoff:    512 * time.Millisecond,
		DialTimeout:        5 * time.Second,
		ReadTimeout:        3 * time.Second,
		WriteTimeout:       3 * time.Second,

		// 消费者性能配置
		ConsumerWorkerCount:   20,                     // 消费者工作池大小
		ConsumerBufferSize:    2000,                   // 消费者缓冲区大小
		ConsumerBatchSize:     1,                      // 批处理大小
		ConsumerPollTimeout:   time.Second,            // 轮询超时
		ConsumerRetryInterval: 500 * time.Millisecond, // 重试间隔
		ConsumerMaxRetries:    5,                      // 最大重试次数

		// 生产者性能配置
		ProducerBatchSize:     200,                   // 生产者批处理大小
		ProducerFlushInterval: 50 * time.Millisecond, // 刷新间隔
		ProducerCompression:   true,                  // 启用压缩

		// 序列化配置
		SerializationType:        "msgpack", // 使用MessagePack序列化
		SerializationCompression: true,      // 启用序列化压缩

		// 对象池配置
		ObjectPoolEnabled:           true, // 启用对象池
		ObjectPoolMaxMessageObjects: 2000, // 消息对象池大小
		ObjectPoolMaxBufferObjects:  1000, // 缓冲区对象池大小
	}

	switch cfg.Redis.MQ.Mode {
	case "cluster": // 集群模式
		redisCfg.Mode = mqConfig.RedisModeCluster
		redisCfg.Addrs = cfg.Redis.MQ.Cluster.Addrs
		redisCfg.Password = cfg.Redis.MQ.Cluster.Password

	case "sentinel": // 哨兵模式
		redisCfg.Mode = mqConfig.RedisModeSentinel
		redisCfg.SentinelAddrs = cfg.Redis.MQ.Sentinel.SentinelAddrs
		redisCfg.SentinelPassword = cfg.Redis.MQ.Sentinel.SentinelPassword
		redisCfg.MasterName = cfg.Redis.MQ.Sentinel.MasterName
		redisCfg.DB = cfg.Redis.MQ.Sentinel.DB
		redisCfg.Password = cfg.Redis.MQ.Sentinel.Password

	default: // 单机模式
		redisCfg.Mode = mqConfig.RedisModeStandalone
		redisCfg.Addr = cfg.Redis.MQ.Host + ":" + cfg.Redis.MQ.Port
		redisCfg.Password = cfg.Redis.MQ.Password
		redisCfg.DB = cfg.Redis.MQ.DB
	}
	return redisCfg
}

func createRabbitMQConfig(cfg *config.Config) mqConfig.RabbitMQConfig {
	port, _ := strconv.Atoi(cfg.MessageQueue.Amqp.Port)
	if port <= 0 {
		port = 5672
	}
	return mqConfig.RabbitMQConfig{
		Host:              cfg.MessageQueue.Amqp.Host,
		Port:              port,
		Username:          cfg.MessageQueue.Amqp.Username,
		Password:          cfg.MessageQueue.Amqp.Password,
		VHost:             cfg.MessageQueue.Amqp.Vhost,
		ExchangeType:      "direct",
		QueueDurable:      true,
		QueueAutoDelete:   false,
		QueueExclusive:    false,
		QueueNoWait:       false,
		QoS:               50,               // 增大预取数量
		Heartbeat:         30 * time.Second, // 心跳间隔
		ConnectionTimeout: 10 * time.Second, // 连接超时
		ChannelMax:        200,              // 最大通道数
		FrameSize:         131072,           // 帧大小

		// 连接池配置（高性能）
		PoolSize:        20, // 连接池大小
		MinConnections:  5,  // 最小连接数
		MaxConnections:  50, // 最大连接数
		ChannelPoolSize: 10, // 通道池大小

		// 重连配置
		MaxRetries:     3,                      // 最大重试次数
		RetryInterval:  500 * time.Millisecond, // 重试间隔
		ReconnectDelay: 2 * time.Second,        // 重连延迟

		// 性能配置
		Performance: mqConfig.PerformanceConfig{
			// 消费者性能配置
			Consumer: mqConfig.ConsumerPerformanceConfig{
				WorkerCount:   20,                     // 消费者工作池大小
				BufferSize:    2000,                   // 消费者缓冲区大小
				BatchSize:     10,                     // 批处理大小
				PollTimeout:   time.Second,            // 轮询超时
				RetryInterval: 500 * time.Millisecond, // 重试间隔
				MaxRetries:    5,                      // 最大重试次数
			},
			// 生产者性能配置
			Producer: mqConfig.ProducerPerformanceConfig{
				BatchSize:     200,                   // 生产者批处理大小
				FlushInterval: 50 * time.Millisecond, // 刷新间隔
				Compression:   true,                  // 启用压缩
			},
			// 序列化配置
			Serialization: mqConfig.SerializationConfig{
				Type:        "msgpack", // 使用MessagePack序列化
				Compression: true,      // 启用序列化压缩
			},
			// 对象池配置
			ObjectPool: mqConfig.ObjectPoolConfig{
				Enabled:           true, // 启用对象池
				MaxMessageObjects: 2000, // 消息对象池大小
				MaxBufferObjects:  1000, // 缓冲区对象池大小
			},
		},
	}
}

func createMemoryConfig() mqConfig.MemoryConfig {
	return mqConfig.MemoryConfig{
		MaxQueueSize:       10000,
		MaxDelayQueueSize:  2000,
		DelayCheckInterval: 100 * time.Millisecond,
		EnableMetrics:      true,
	}
}
