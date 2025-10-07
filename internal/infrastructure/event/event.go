package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Event 泛型事件接口
type Event[T any] interface {
	// EventType 返回事件类型
	EventType() string
	// OccurredAt 返回事件发生时间
	OccurredAt() time.Time
	// Payload 返回事件数据
	Payload() T
}

// BaseEvent 基础事件实现
type BaseEvent[T any] struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Data      T         `json:"data"`
}

func (e BaseEvent[T]) EventType() string {
	return e.Type
}

func (e BaseEvent[T]) OccurredAt() time.Time {
	return e.Timestamp
}

func (e BaseEvent[T]) Payload() T {
	return e.Data
}

// NewEvent 创建新事件
func NewEvent[T any](eventType string, data T) Event[T] {
	return BaseEvent[T]{
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
	}
}

// DomainEvent 领域事件接口
type DomainEvent[T any] interface {
	Event[T]
	// AggregateID 返回聚合根ID
	AggregateID() string
	// AggregateName 返回聚合根名称
	AggregateName() string
}

// BaseDomainEvent 基础领域事件实现
type BaseDomainEvent[T any] struct {
	BaseEvent[T]
	AggID   string `json:"aggregate_id"`
	AggName string `json:"aggregate_name"`
}

func (e BaseDomainEvent[T]) AggregateID() string {
	return e.AggID
}

func (e BaseDomainEvent[T]) AggregateName() string {
	return e.AggName
}

// NewDomainEvent 创建领域事件
func NewDomainEvent[T any](eventType string, aggregateID string, aggregateName string, data T) DomainEvent[T] {
	return BaseDomainEvent[T]{
		BaseEvent: BaseEvent[T]{
			Type:      eventType,
			Timestamp: time.Now(),
			Data:      data,
		},
		AggID:   aggregateID,
		AggName: aggregateName,
	}
}

// Handler 事件处理器接口
type Handler interface {
	// HandleTypedEvent 处理类型化事件
	HandleTypedEvent(ctx context.Context, event any) error
	// InterestedEventTypes 返回感兴趣的事件类型列表
	InterestedEventTypes() []string
}

// Bus 事件总线接口
type Bus interface {
	// PublishEvent 发布事件
	PublishEvent(ctx context.Context, event any) error
	// SubscribeHandler 订阅事件处理器
	SubscribeHandler(handler any) error
}

// DomainEventHandler 领域事件基础处理器
type DomainEventHandler[T any] struct{}

// ParseDomainEvent 解析领域事件
func (h *DomainEventHandler[T]) ParseDomainEvent(ctx context.Context, event any) (BaseDomainEvent[T], error) {
	domainEventRaw, ok := event.(BaseDomainEvent[json.RawMessage])
	var zero BaseDomainEvent[T]
	if !ok {
		return zero, fmt.Errorf("expected event.BaseDomainEvent[json.RawMessage], got %T", event)
	}

	// 解析JSON数据为类型T结构
	var payload T
	if err := json.Unmarshal(domainEventRaw.Data, &payload); err != nil {
		return zero, fmt.Errorf("failed to unmarshal %T event: %w", payload, err)
	}

	// 创建强类型的领域事件
	domainEvent := BaseDomainEvent[T]{
		BaseEvent: BaseEvent[T]{
			Type:      domainEventRaw.Type,
			Timestamp: domainEventRaw.Timestamp,
			Data:      payload,
		},
		AggID:   domainEventRaw.AggID,
		AggName: domainEventRaw.AggName,
	}

	return domainEvent, nil
}
