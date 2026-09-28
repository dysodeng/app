package event

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	mqv2 "github.com/goairix/mq/v2"
)

// MQEventBus 基于MQ的事件总线实现
type MQEventBus struct {
	producer mqv2.Publisher
}

// NewMQEventBus 创建基于MQ的事件总线
func NewMQEventBus(producer mqv2.Publisher) *MQEventBus {
	return &MQEventBus{
		producer: producer,
	}
}

// Publish 发布事件
func (b *MQEventBus) Publish(ctx context.Context, eventType string, eventData []byte) error {
	msg, err := mqv2.NewMessage(eventType, eventData)
	if err != nil {
		return fmt.Errorf("create MQ event message: %w", err)
	}
	return b.producer.Publish(ctx, msg)
}

// PublishEvent 发布事件
func (b *MQEventBus) PublishEvent(ctx context.Context, event any) error {
	// 检查是否实现了EventType方法
	if e, ok := event.(interface{ EventType() string }); ok {
		data, err := sonic.Marshal(event)
		if err != nil {
			return fmt.Errorf("marshal event failed: %w", err)
		}
		return b.Publish(ctx, e.EventType(), data)
	}
	return fmt.Errorf("unsupported event type: %T", event)
}

// SubscribeHandler 订阅事件处理器
func (b *MQEventBus) SubscribeHandler(handler any) error {
	if _, ok := handler.(interface{ InterestedEventTypes() []string }); ok {
		// 实际订阅由ConsumerService处理
		return nil
	}
	return fmt.Errorf("unsupported handler type: %T", handler)
}
