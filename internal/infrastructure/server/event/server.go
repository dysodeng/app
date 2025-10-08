package event

import (
	"context"

	diEvent "github.com/dysodeng/app/internal/di/event"
	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/event"
)

// Server 事件消费者服务
type Server struct {
	cfg           *config.Config
	eventConsumer *event.ConsumerService
	registry      *diEvent.HandlerRegistry
}

func NewEventServer(
	cfg *config.Config,
	eventConsumer *event.ConsumerService,
	registry *diEvent.HandlerRegistry,
) *Server {
	return &Server{
		cfg:           cfg,
		eventConsumer: eventConsumer,
		registry:      registry,
	}
}

func (s *Server) IsEnabled() bool {
	return s.cfg.Server.Event.Enabled
}

func (s *Server) Addr() string {
	return ""
}

func (s *Server) Name() string {
	return "Event"
}

func (s *Server) Start() error {
	// 事件注册
	for _, handler := range s.registry.Handlers() {
		if err := s.eventConsumer.SubscribeHandler(handler); err != nil {
			return err
		}
	}
	ctx := context.Background()
	return s.eventConsumer.Start(ctx)
}

func (s *Server) Stop(_ context.Context) error {
	return s.eventConsumer.Stop()
}
