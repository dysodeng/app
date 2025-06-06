package event

import sharedEvent "github.com/dysodeng/app/internal/domain/shared/event"

const UserCreatedEventType = "user.created"

// UserCreatedEvent 用户创建成功事件
type UserCreatedEvent struct {
	sharedEvent.BaseDomainEvent
	UserID    uint64 `json:"user_id"`
	Telephone string `json:"telephone"`
}

func NewUserCreatedEvent(userId uint64, telephone string) *UserCreatedEvent {
	return &UserCreatedEvent{
		BaseDomainEvent: sharedEvent.NewBaseDomainEvent(UserCreatedEventType, userId),
		UserID:          userId,
		Telephone:       telephone,
	}
}

func (e *UserCreatedEvent) EventData() map[string]interface{} {
	return map[string]interface{}{
		"user_id":   e.UserID,
		"telephone": e.Telephone,
	}
}
