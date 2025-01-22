package job

import "github.com/dysodeng/mq/message"

// Handler 队列任务处理接口
type Handler interface {
	// QueueKey 队列类型key
	QueueKey() string
	// IsDelay 是否延时队列
	IsDelay() bool
	// Handle 队列处理器
	Handle(message message.Message) error
}
