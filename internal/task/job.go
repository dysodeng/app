package task

import "github.com/dysodeng/mq/message"

// JobInterface 队列任务接口
type JobInterface interface {
	// QueueKey 队列类型key
	QueueKey() string
	// IsDelay 是否延时队列
	IsDelay() bool
	// Handle 队列处理器
	Handle(message message.Message) error
}
