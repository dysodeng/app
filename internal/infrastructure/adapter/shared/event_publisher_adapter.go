package shared

import (
	"context"

	domainEvent "github.com/dysodeng/app/internal/domain/shared/event"
	domainPort "github.com/dysodeng/app/internal/domain/shared/port"
	infraEvent "github.com/dysodeng/app/internal/infrastructure/event"
)

// EventPublisherAdapter 事件发布器适配器
type EventPublisherAdapter struct {
	bus infraEvent.Bus
}

func NewEventPublisherAdapter(bus infraEvent.Bus) domainPort.EventPublisher {
	return &EventPublisherAdapter{bus: bus}
}

func (a *EventPublisherAdapter) Publish(ctx context.Context, e domainEvent.DomainEvent[any]) error {
	// 转换为基础设施领域事件
	return a.bus.PublishEvent(ctx, infraEvent.NewDomainEvent(e.Type, e.AggregateID, e.AggregateName, e.Payload))
}
