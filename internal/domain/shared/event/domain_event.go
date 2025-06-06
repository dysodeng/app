package event

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent 领域事件接口
type DomainEvent interface {
	// EventID 返回事件的唯一标识符
	EventID() uuid.UUID
	// EventType 返回事件的类型名称
	EventType() string
	// AggregateID 返回触发该事件的聚合根ID
	AggregateID() uint64
	// OccurredOn 返回事件发生的时间戳
	OccurredOn() time.Time
	// EventData 返回事件的具体数据载荷
	EventData() map[string]interface{}
}

// BaseDomainEvent 基础领域事件
type BaseDomainEvent struct {
	id          uuid.UUID // 事件ID
	eventType   string    // 事件类型
	aggregateID uint64    // 触发事件的聚合根ID
	occurredOn  time.Time // 事件发生时间
}

func NewBaseDomainEvent(eventType string, aggregateID uint64) BaseDomainEvent {
	return BaseDomainEvent{
		id:          uuid.New(),
		eventType:   eventType,
		aggregateID: aggregateID,
		occurredOn:  time.Now(),
	}
}

func (e BaseDomainEvent) EventID() uuid.UUID    { return e.id }
func (e BaseDomainEvent) EventType() string     { return e.eventType }
func (e BaseDomainEvent) AggregateID() uint64   { return e.aggregateID }
func (e BaseDomainEvent) OccurredOn() time.Time { return e.occurredOn }
