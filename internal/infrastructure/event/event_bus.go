package event

import (
	"sync"
)

// Bus 事件总线
type Bus struct {
	handlers map[string][]interface{}
	mu       sync.RWMutex
}

// NewEventBus 创建事件总线
func NewEventBus() *Bus {
	return &Bus{
		handlers: make(map[string][]interface{}),
	}
}

// RegisterHandler 注册事件处理器
func (b *Bus) RegisterHandler(handler interface{}) {
	// 简化实现，实际应该根据handler类型注册到对应事件
	b.mu.Lock()
	defer b.mu.Unlock()

	// 这里简化处理，实际应该通过反射获取handler处理的事件类型
	eventType := "default"
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish 发布事件
func (b *Bus) Publish(event interface{}) {
	// 简化实现，实际应该根据event类型找到对应的handlers并调用
	b.mu.RLock()
	defer b.mu.RUnlock()

	// 这里简化处理，实际应该通过反射获取event类型并找到对应的handlers
	eventType := "default"
	if handlers, ok := b.handlers[eventType]; ok {
		for range handlers {
			// 这里简化处理，实际应该通过反射调用handler的方法
			// handler.Handle(event)
		}
	}
}
