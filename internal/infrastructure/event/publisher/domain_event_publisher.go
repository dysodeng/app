package publisher

import (
	"context"

	domainEvent "github.com/dysodeng/app/internal/domain/shared/event"
	"github.com/dysodeng/app/internal/infrastructure/event/bus"
)

// DomainEventPublisher 领域事件发布器
type DomainEventPublisher struct {
	eventBus bus.EventBus
}

func NewDomainEventPublisher(eventBus bus.EventBus) *DomainEventPublisher {
	return &DomainEventPublisher{
		eventBus: eventBus,
	}
}

func (p *DomainEventPublisher) Publish(ctx context.Context, event domainEvent.DomainEvent) error {
	return p.eventBus.Publish(ctx, event)
}

func (p *DomainEventPublisher) PublishAll(ctx context.Context, events []domainEvent.DomainEvent) error {
	for _, event := range events {
		if err := p.Publish(ctx, event); err != nil {
			return err
		}
	}
	return nil
}
