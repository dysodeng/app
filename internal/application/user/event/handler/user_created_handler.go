package handler

import (
	"context"
	"log"

	"github.com/dysodeng/app/internal/application/user/service"
	domainEvent "github.com/dysodeng/app/internal/domain/shared/event"
	"github.com/dysodeng/app/internal/domain/user/event"
)

// UserCreatedHandler 用户创建成功事件处理器
type UserCreatedHandler struct {
	userService service.UserApplicationService
}

func NewUserCreatedHandler(userService service.UserApplicationService) *UserCreatedHandler {
	return &UserCreatedHandler{
		userService: userService,
	}
}

func (handler *UserCreatedHandler) Handle(ctx context.Context, event domainEvent.DomainEvent) error {
	// 这里处理事件相关业务
	log.Printf("%+v", event.EventData())
	return nil
}

func (handler *UserCreatedHandler) CanHandle(eventType string) bool {
	return eventType == event.UserCreatedEventType
}
