package event

import "time"

// DomainEvent 领域事件
type DomainEvent[T any] struct {
	Type          string    `json:"type"`
	OccurredAt    time.Time `json:"timestamp"`
	Payload       T         `json:"data"`
	AggregateID   string    `json:"aggregate_id,omitempty"`
	AggregateName string    `json:"aggregate_name,omitempty"`
}

// NewDomainEvent 创建领域事件
func NewDomainEvent[T any](eventType string, aggregateID string, aggregateName string, data T) DomainEvent[T] {
	return DomainEvent[T]{
		Type:          eventType,
		OccurredAt:    time.Now(),
		Payload:       data,
		AggregateID:   aggregateID,
		AggregateName: aggregateName,
	}
}
