package event

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/dysodeng/mq/contract"
	"github.com/dysodeng/mq/message"
	"go.uber.org/zap"
)

// handlerWrapper 处理器包装器
type handlerWrapper struct {
	handler Handler
	logger  *zap.Logger
}

// createEvent 根据事件数据创建事件对象
func (w *handlerWrapper) createEvent(eventType string, eventData json.RawMessage, timestamp time.Time, aggregateID, aggregateName string) any {
	if aggregateID != "" && aggregateName != "" {
		// 创建领域事件
		return BaseDomainEvent[json.RawMessage]{
			BaseEvent: BaseEvent[json.RawMessage]{
				Type:      eventType,
				Timestamp: timestamp,
				Data:      eventData,
			},
			AggID:   aggregateID,
			AggName: aggregateName,
		}
	}

	// 创建普通事件
	return BaseEvent[json.RawMessage]{
		Type:      eventType,
		Timestamp: timestamp,
		Data:      eventData,
	}
}

// handle 处理事件
func (w *handlerWrapper) handle(ctx context.Context, eventType string, eventData json.RawMessage, timestamp time.Time, aggregateID, aggregateName string) error {
	// 记录开始处理事件
	w.logger.Debug("Handler starting to process event",
		zap.String("eventType", eventType),
		zap.String("handlerType", fmt.Sprintf("%T", w.handler)),
		zap.String("aggregateID", aggregateID),
		zap.String("aggregateName", aggregateName),
	)

	// 创建事件对象
	event := w.createEvent(eventType, eventData, timestamp, aggregateID, aggregateName)
	if event == nil {
		err := fmt.Errorf("failed to create event object for type: %s", eventType)
		w.logger.Error("Event creation failed", zap.Error(err))
		return err
	}

	// 处理事件
	if err := w.handler.HandleTypedEvent(ctx, event); err != nil {
		w.logger.Error("Handler failed to process event",
			zap.String("eventType", eventType),
			zap.String("handlerType", fmt.Sprintf("%T", w.handler)),
			zap.Error(err),
		)
		return fmt.Errorf("handler %T failed to process event %s: %w", w.handler, eventType, err)
	}

	w.logger.Debug("Handler successfully processed event",
		zap.String("eventType", eventType),
		zap.String("handlerType", fmt.Sprintf("%T", w.handler)),
	)
	return nil
}

// ConsumerService 事件消费者服务
type ConsumerService struct {
	consumer contract.Consumer
	handlers map[string][]*handlerWrapper // 存储处理器包装器
	logger   *zap.Logger
	mu       sync.RWMutex
}

// NewEventConsumerService 创建事件消费者服务
func NewEventConsumerService(consumer contract.Consumer, logger *zap.Logger) *ConsumerService {
	return &ConsumerService{
		consumer: consumer,
		handlers: make(map[string][]*handlerWrapper),
		logger:   logger,
	}
}

// SubscribeHandler 订阅事件处理器
func (s *ConsumerService) SubscribeHandler(handler any) error {
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 类型断言为EventHandler
	eventHandler, ok := handler.(Handler)
	if !ok {
		return fmt.Errorf("handler %T does not implement EventHandler interface", handler)
	}

	// 验证感兴趣的事件类型
	interestedTypes := eventHandler.InterestedEventTypes()
	if len(interestedTypes) == 0 {
		return fmt.Errorf("handler %T must be interested in at least one event type", handler)
	}

	// 创建处理器包装器
	wrapper := &handlerWrapper{
		handler: eventHandler,
		logger:  s.logger,
	}

	// 为每个感兴趣的事件类型注册包装器
	registeredCount := 0
	for _, eventType := range interestedTypes {
		if eventType == "" {
			s.logger.Warn("Skipping empty event type for handler",
				zap.String("handlerType", fmt.Sprintf("%T", handler)),
			)
			continue
		}

		s.handlers[eventType] = append(s.handlers[eventType], wrapper)
		registeredCount++

		s.logger.Info("Registered handler for event type",
			zap.String("eventType", eventType),
			zap.String("handlerType", fmt.Sprintf("%T", handler)),
			zap.Int("totalHandlersForType", len(s.handlers[eventType])),
		)
	}

	if registeredCount == 0 {
		return fmt.Errorf("no valid event types found for handler %T", handler)
	}

	s.logger.Info("Handler registration completed",
		zap.String("handlerType", fmt.Sprintf("%T", handler)),
		zap.Int("registeredEventTypes", registeredCount),
	)

	return nil
}

// Start 启动事件消费服务
func (s *ConsumerService) Start(ctx context.Context) error {
	s.logger.Info("Starting event consumer service")

	// 获取所有需要订阅的事件类型
	eventTypes := make([]string, 0)
	s.mu.RLock()
	for eventType := range s.handlers {
		eventTypes = append(eventTypes, eventType)
	}
	s.mu.RUnlock()

	// 如果没有处理器，则不需要启动消费
	if len(eventTypes) == 0 {
		s.logger.Warn("No event handlers registered, skipping consumer start")
		return nil
	}

	// 为每个事件类型单独订阅事件处理器
	for _, eventType := range eventTypes {
		if err := s.consumer.Subscribe(ctx, eventType, s.handleMessage); err != nil {
			s.logger.Error("Failed to subscribe to event type",
				zap.String("eventType", eventType),
				zap.Error(err),
			)
			return fmt.Errorf("failed to subscribe to event type %s: %w", eventType, err)
		}
		s.logger.Info("Subscribed to event type", zap.String("eventType", eventType))
	}

	s.logger.Info("Event consumer service started successfully",
		zap.Int("totalEventTypes", len(eventTypes)),
		zap.Strings("eventTypes", eventTypes),
	)
	return nil
}

// Stop 停止事件消费服务
func (s *ConsumerService) Stop() error {
	s.logger.Info("Stopping event consumer service")

	if err := s.consumer.Close(); err != nil {
		s.logger.Error("Failed to stop consumer", zap.Error(err))
		return fmt.Errorf("failed to stop consumer: %w", err)
	}

	s.logger.Info("Event consumer service stopped successfully")
	return nil
}

// handleMessage 处理消息
func (s *ConsumerService) handleMessage(ctx context.Context, msg *message.Message) error {
	eventType := msg.Topic

	// 获取注册的处理器
	handlers := s.getHandlers(eventType)
	if len(handlers) == 0 {
		s.logger.Warn("No handlers registered for event type", zap.String("eventType", eventType))
		return nil
	}

	// 解析事件数据
	data, err := s.parseEventData(msg.Payload)
	if err != nil {
		s.logger.Error("Failed to parse event data",
			zap.String("eventType", eventType),
			zap.Error(err),
		)
		return fmt.Errorf("failed to parse event data: %w", err)
	}

	// 处理事件
	return s.processEvent(ctx, eventType, data, handlers)
}

// getHandlers 获取指定事件类型的处理器
func (s *ConsumerService) getHandlers(eventType string) []*handlerWrapper {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.handlers[eventType]
}

// parseEventData 解析事件数据
func (s *ConsumerService) parseEventData(payload []byte) (*eventData, error) {
	var data eventData
	if err := json.Unmarshal(payload, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// eventData 事件数据结构
type eventData struct {
	Type          string          `json:"type"`
	Timestamp     time.Time       `json:"timestamp"`
	Data          json.RawMessage `json:"data"`
	AggregateID   string          `json:"aggregate_id,omitempty"`
	AggregateName string          `json:"aggregate_name,omitempty"`
}

// processEvent 处理事件
func (s *ConsumerService) processEvent(ctx context.Context, eventType string, data *eventData, handlers []*handlerWrapper) error {
	var errs []error

	for i, wrapper := range handlers {
		s.logger.Debug("Processing event with handler",
			zap.String("eventType", eventType),
			zap.Int("handlerIndex", i),
		)

		if err := wrapper.handle(ctx, eventType, data.Data, data.Timestamp, data.AggregateID, data.AggregateName); err != nil {
			s.logger.Error("Handler failed to process event",
				zap.String("eventType", eventType),
				zap.Int("handlerIndex", i),
				zap.Error(err),
			)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to process event %s with %d errors: %v", eventType, len(errs), errs)
	}

	s.logger.Debug("Event processed successfully",
		zap.String("eventType", eventType),
		zap.Int("handlerCount", len(handlers)),
	)
	return nil
}

// PublishEvent 发布事件
func (s *ConsumerService) PublishEvent(ctx context.Context, event any) error {
	// ConsumerService不应该负责发布事件
	// 这里只是为了实现TypedEventBus接口
	return fmt.Errorf("ConsumerService does not support publishing events")
}
