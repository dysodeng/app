package mq

import (
	"log"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/telemetry/metrics"

	"github.com/dysodeng/mq"
	mqConfig "github.com/dysodeng/mq/config"
	"github.com/dysodeng/mq/contract"
)

const QueuePrefix = "app"

var factory *mq.Factory
var mqInstance contract.MQ

func init() {
	cfg := mqConfig.Config{
		KeyPrefix: QueuePrefix,
	}

	switch config.MessageQueue.Driver {
	case "redis":
		cfg.Adapter = mqConfig.AdapterRedis
		cfg.Redis = createRedisConfig()
	case "amqp":
		cfg.Adapter = mqConfig.AdapterRabbitMQ
		cfg.RabbitMQ = createRabbitMQConfig()
	default:
		cfg.Adapter = mqConfig.AdapterMemory
		cfg.Memory = createMemoryConfig()
	}

	factory = mq.NewFactory(
		cfg,
		mq.WithObserver(&metricsObserver{ // 可观测性
			meter:  metrics.Meter(),
			logger: logger.ZapLogger(),
		}),
	)

	var err error
	mqInstance, err = factory.CreateMQ()
	if err != nil {
		log.Fatal("Failed to create MQ:", err)
	}
}

// Instance MQ实例
func Instance() contract.MQ {
	return mqInstance
}

// Consumer 消费者
func Consumer() contract.Consumer {
	return mqInstance.Consumer()
}

// Producer 生产者
func Producer() contract.Producer {
	return mqInstance.Producer()
}

func QueueKey(queueKey string) string {
	return QueuePrefix + queueKey
}
