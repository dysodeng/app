package job

import "github.com/dysodeng/mq/message"

// Interface 队列任务接口
type Interface interface {
	// QueueKey 队列类型key
	QueueKey() string
	// IsDelay 是否延时队列
	IsDelay() bool
	// Handle 队列处理器
	Handle(message message.Message) error
}
