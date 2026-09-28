package mq

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	mqv2 "github.com/goairix/mq/v2"
)

type preparedBackend interface {
	backend
	Prepare(context.Context, mqv2.Subscription) error
}

type connectionState struct {
	adapter         preparedBackend
	isClosed        func() bool
	closeConnection func() error
}

func (s *connectionState) close(ctx context.Context) error {
	adapterErr := s.adapter.Close(ctx)
	var connectionErr error
	if s.closeConnection != nil {
		connectionErr = s.closeConnection()
	}
	return errors.Join(adapterErr, connectionErr)
}

// reconnectingBackend rebuilds a RabbitMQ adapter after its connection closes.
// It never retries a publish because the previous outcome may be unknown.
type reconnectingBackend struct {
	mu      sync.Mutex
	factory func() (*connectionState, error)
	current *connectionState
	closed  bool
}

func newReconnectingBackend(factory func() (*connectionState, error)) (*reconnectingBackend, error) {
	state, err := factory()
	if err != nil {
		return nil, err
	}
	if state == nil || state.adapter == nil {
		return nil, fmt.Errorf("MQ connection factory returned no adapter")
	}
	return &reconnectingBackend{factory: factory, current: state}, nil
}

func (b *reconnectingBackend) currentState() (*connectionState, error) {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return nil, mqv2.ErrClosed
	}
	if b.current != nil && (b.current.isClosed == nil || !b.current.isClosed()) {
		state := b.current
		b.mu.Unlock()
		return state, nil
	}
	fresh, err := b.factory()
	if err != nil {
		b.mu.Unlock()
		return nil, fmt.Errorf("reconnect MQ: %w", err)
	}
	if fresh == nil || fresh.adapter == nil {
		b.mu.Unlock()
		return nil, fmt.Errorf("MQ connection factory returned no adapter")
	}
	previous := b.current
	b.current = fresh
	b.mu.Unlock()
	if previous != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := previous.close(ctx); err != nil {
			slog.Warn("关闭旧 MQ 连接失败", "error", err)
		}
		cancel()
	}
	return fresh, nil
}

func (b *reconnectingBackend) Publish(ctx context.Context, msg mqv2.Message) error {
	state, err := b.currentState()
	if err != nil {
		return err
	}
	return state.adapter.Publish(ctx, msg)
}

func (b *reconnectingBackend) Prepare(ctx context.Context, sub mqv2.Subscription) error {
	state, err := b.currentState()
	if err != nil {
		return err
	}
	return state.adapter.Prepare(ctx, sub)
}

func (b *reconnectingBackend) Run(ctx context.Context, sub mqv2.Subscription, handler mqv2.Handler) error {
	state, err := b.currentState()
	if err != nil {
		return err
	}
	return state.adapter.Run(ctx, sub, handler)
}

func (b *reconnectingBackend) Close(ctx context.Context) error {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return nil
	}
	b.closed = true
	state := b.current
	b.current = nil
	b.mu.Unlock()
	if state == nil {
		return nil
	}
	return state.close(ctx)
}
