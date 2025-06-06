package bus

import (
	"context"
	"sync"

	domainEvent "github.com/dysodeng/app/internal/domain/shared/event"
	"github.com/dysodeng/app/internal/pkg/logger"
)

// EventHandler 事件处理器接口
type EventHandler interface {
	// Handle 事件处理器
	Handle(ctx context.Context, event domainEvent.DomainEvent) error
	// CanHandle 判断是否可以处理该事件
	CanHandle(eventType string) bool
}

// EventBus 事件总线接口
type EventBus interface {
	// Publish 发布事件
	Publish(ctx context.Context, event domainEvent.DomainEvent) error
	// Subscribe 订阅事件
	Subscribe(handler EventHandler)
	// Unsubscribe 取消订阅事件
	Unsubscribe(handler EventHandler)
}

// InMemoryEventBus 基于内存的事件总线
type InMemoryEventBus struct {
	handlers []EventHandler
	mu       sync.RWMutex
}

func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers: make([]EventHandler, 0),
	}
}

func (bus *InMemoryEventBus) Publish(ctx context.Context, event domainEvent.DomainEvent) error {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	for _, handler := range bus.handlers {
		if handler.CanHandle(event.EventType()) {
			go func(h EventHandler) {
				if err := h.Handle(ctx, event); err != nil {
					logger.Error(ctx, "事件处理失败", logger.ErrorField(err))
				}
			}(handler)
		}
	}
	return nil
}

func (bus *InMemoryEventBus) Subscribe(handler EventHandler) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	bus.handlers = append(bus.handlers, handler)
}

func (bus *InMemoryEventBus) Unsubscribe(handler EventHandler) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	for i, h := range bus.handlers {
		if h == handler {
			bus.handlers = append(bus.handlers[:i], bus.handlers[i+1:]...)
			break
		}
	}
}
