package event

import (
	"context"
	"sync"
	"time"
)

// Event 领域事件接口
type Event interface {
	// Name 事件名称
	Name() string
	// OccurredTime 事件发生时间
	OccurredTime() time.Time
	// Payload 事件负载
	Payload() interface{}
}

// Handler 事件处理器接口
type Handler interface {
	// Handle 处理事件
	Handle(ctx context.Context, event Event) error
}

// Bus 事件总线
type Bus struct {
	handlers map[string][]Handler
	mu       sync.RWMutex
}

// NewEventBus 创建事件总线
func NewEventBus() *Bus {
	return &Bus{
		handlers: make(map[string][]Handler),
	}
}

// Subscribe 订阅事件
func (b *Bus) Subscribe(eventName string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.handlers[eventName]; !ok {
		b.handlers[eventName] = make([]Handler, 0)
	}
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

// Publish 发布事件
func (b *Bus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	handlers, ok := b.handlers[event.Name()]
	b.mu.RUnlock()

	if !ok {
		return nil
	}

	for _, handler := range handlers {
		if err := handler.Handle(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

// BaseEvent 基础事件
type BaseEvent struct {
	name         string
	occurredTime time.Time
	payload      interface{}
}

// NewBaseEvent 创建基础事件
func NewBaseEvent(name string, payload interface{}) *BaseEvent {
	return &BaseEvent{
		name:         name,
		occurredTime: time.Now(),
		payload:      payload,
	}
}

// Name 事件名称
func (e *BaseEvent) Name() string {
	return e.name
}

// OccurredTime 事件发生时间
func (e *BaseEvent) OccurredTime() time.Time {
	return e.occurredTime
}

// Payload 事件负载
func (e *BaseEvent) Payload() interface{} {
	return e.payload
}
