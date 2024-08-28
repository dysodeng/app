package mq

import (
	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/mq"
	"github.com/dysodeng/mq/contract"
	"github.com/dysodeng/mq/driver/amqp"
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
	switch config.MQ.Driver {
	case string(mq.Amqp):
		return mq.NewQueueConsumer(mq.Amqp, QueueKey(queueKey), &amqp.Config{
			Host:     config.MQ.Amqp.Host + ":" + config.MQ.Amqp.Port,
			Username: config.MQ.Amqp.Username,
			Password: config.MQ.Amqp.Password,
			VHost:    config.MQ.Amqp.Vhost,
		})
	}
	return nil, errors.New("mq driver not found.")
}

// NewMessageQueueProducer 创建消息队列生产者
func NewMessageQueueProducer(pool *contract.Pool) (contract.Producer, error) {
	switch config.MQ.Driver {
	case string(mq.Amqp):
		return mq.NewQueueProducer(mq.Amqp, &amqp.Config{
			Host:     config.MQ.Amqp.Host + ":" + config.MQ.Amqp.Port,
			Username: config.MQ.Amqp.Username,
			Password: config.MQ.Amqp.Password,
			VHost:    config.MQ.Amqp.Vhost,
			Pool:     pool,
		})
	}
	return nil, errors.New("mq driver not found.")
}
