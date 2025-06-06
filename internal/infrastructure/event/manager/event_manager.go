package manager

import (
	"github.com/dysodeng/app/internal/infrastructure/event/bus"
)

// EventManager 事件管理器
type EventManager struct {
	eventBus bus.EventBus
}

// NewEventManager 接受任意数量的 EventHandler 接口
func NewEventManager(
	eventBus bus.EventBus,
	handlers ...bus.EventHandler, // 使用可变参数和接口
) *EventManager {
	// 注册所有事件处理器
	for _, handler := range handlers {
		eventBus.Subscribe(handler)
	}

	return &EventManager{
		eventBus: eventBus,
	}
}

func (em *EventManager) EventBus() bus.EventBus {
	return em.eventBus
}

func (em *EventManager) Register(handler bus.EventHandler) {
	em.eventBus.Subscribe(handler)
}
