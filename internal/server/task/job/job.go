package job

import (
	"context"

	"github.com/dysodeng/mq/message"
)

// Handler 队列任务处理接口
type Handler interface {
	// TopicKey 队列主题key
	TopicKey() string
	// Handle 队列处理器
	Handle(ctx context.Context, msg *message.Message) error
}
