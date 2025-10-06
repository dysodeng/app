package mq

import (
	"log"

	"github.com/dysodeng/mq"
	mqConfig "github.com/dysodeng/mq/config"
	"github.com/dysodeng/mq/contract"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/shared/logger"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/metrics"
)

const QueuePrefix = "app-service"

var factory *mq.Factory
var mqInstance contract.MQ

func Init(cfg *config.Config) (contract.MQ, error) {
	mqCfg := mqConfig.Config{
		KeyPrefix: QueuePrefix,
	}

	if cfg.MessageQueue.Enabled {
		switch cfg.MessageQueue.Driver {
		case "redis":
			mqCfg.Adapter = mqConfig.AdapterRedis
			mqCfg.Redis = createRedisConfig(cfg)
		case "amqp":
			mqCfg.Adapter = mqConfig.AdapterRabbitMQ
			mqCfg.RabbitMQ = createRabbitMQConfig(cfg)
		default:
			mqCfg.Adapter = mqConfig.AdapterMemory
			mqCfg.Memory = createMemoryConfig()
		}
	} else {
		mqCfg.Adapter = mqConfig.AdapterMemory
		mqCfg.Memory = createMemoryConfig()
	}

	factory = mq.NewFactory(
		mqCfg,
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

	return mqInstance, nil
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
	return queueKey
}
