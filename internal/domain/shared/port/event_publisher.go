package port

import (
	"context"

	domainEvent "github.com/dysodeng/app/internal/domain/shared/event"
)

// EventPublisher 事件发布端口
type EventPublisher interface {
	Publish(ctx context.Context, e domainEvent.DomainEvent[any]) error
}
