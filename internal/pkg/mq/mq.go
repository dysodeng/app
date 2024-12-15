package mq

import (
	"fmt"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/mq"
	"github.com/dysodeng/mq/contract"
	"github.com/dysodeng/mq/driver/amqp"
	"github.com/dysodeng/mq/driver/redis"
	"github.com/pkg/errors"
)

const QueuePrefix = "app_"

const (
	TaskManagerDeadline = "task_manager_deadline" // 工作任务过期处理任务
	TaskManagerNotice   = "task_manager_notice"   // 工作任务提醒任务
)

func QueueKey(queueKey string) string {
	return QueuePrefix + queueKey
}

// NewMessageQueueConsumer 创建消息队列消费者
func NewMessageQueueConsumer(queueKey string) (contract.Consumer, error) {
	switch config.MessageQueue.Driver {
	case string(mq.Amqp):
		return mq.NewQueueConsumer(mq.Amqp, QueueKey(queueKey), &amqp.Config{
			Host:     config.MessageQueue.Amqp.Host + ":" + config.MessageQueue.Amqp.Port,
			Username: config.MessageQueue.Amqp.Username,
			Password: config.MessageQueue.Amqp.Password,
			VHost:    config.MessageQueue.Amqp.Vhost,
		})
	case string(mq.Redis):
		return mq.NewQueueConsumer(mq.Redis, QueueKey(queueKey), &redis.Config{
			Addr:     fmt.Sprintf("%s:%s", config.Redis.MQ.Host, config.Redis.MQ.Port),
			DB:       config.Redis.MQ.DB,
			Password: config.Redis.MQ.Password,
		})
	}
	return nil, errors.New("mq driver not found.")
}

// NewMessageQueueProducer 创建消息队列生产者
func NewMessageQueueProducer(pool *contract.Pool) (contract.Producer, error) {
	switch config.MessageQueue.Driver {
	case string(mq.Amqp):
		return mq.NewQueueProducer(mq.Amqp, &amqp.Config{
			Host:     config.MessageQueue.Amqp.Host + ":" + config.MessageQueue.Amqp.Port,
			Username: config.MessageQueue.Amqp.Username,
			Password: config.MessageQueue.Amqp.Password,
			VHost:    config.MessageQueue.Amqp.Vhost,
			Pool:     pool,
		})
	case string(mq.Redis):
		return mq.NewQueueProducer(mq.Redis, &redis.Config{
			Addr:     fmt.Sprintf("%s:%s", config.Redis.MQ.Host, config.Redis.MQ.Port),
			DB:       config.Redis.MQ.DB,
			Password: config.Redis.MQ.Password,
		})
	}
	return nil, errors.New("mq driver not found.")
}
