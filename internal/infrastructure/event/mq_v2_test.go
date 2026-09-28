package event

import (
	"context"
	"encoding/json"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dysodeng/app/internal/infrastructure/config"
	appmq "github.com/dysodeng/app/internal/infrastructure/pkg/mq"
	mqv2 "github.com/goairix/mq/v2"
	"go.uber.org/zap"
)

type testEventHandler struct {
	received chan BaseDomainEvent[json.RawMessage]
}

func (h *testEventHandler) InterestedEventTypes() []string { return []string{"file.uploaded"} }

func (h *testEventHandler) Handle(_ context.Context, event any) error {
	h.received <- event.(BaseDomainEvent[json.RawMessage])
	return nil
}

func TestMQV2PublishesAndConsumesDomainEvent(t *testing.T) {
	client, err := appmq.Init(&config.Config{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Close(context.Background()) })
	bus := NewMQEventBus(client)
	consumer := NewEventConsumerService(client, zap.NewNop())
	handler := &testEventHandler{received: make(chan BaseDomainEvent[json.RawMessage], 1)}
	if err := consumer.SubscribeHandler(handler); err != nil {
		t.Fatal(err)
	}
	if err := consumer.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = consumer.Stop(context.Background()) })
	if err := bus.PublishEvent(context.Background(), NewDomainEvent("file.uploaded", "file-1", "file", map[string]string{"name": "a.txt"})); err != nil {
		t.Fatal(err)
	}
	select {
	case got := <-handler.received:
		if got.Type != "file.uploaded" || got.AggID != "file-1" || got.AggName != "file" || string(got.Data) != `{"name":"a.txt"}` {
			t.Fatalf("unexpected event: %+v", got)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for domain event")
	}
	if err := consumer.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
}

type flakySubscriber struct {
	message mqv2.Message
	calls   atomic.Int32
}

func (s *flakySubscriber) Prepare(context.Context, mqv2.Subscription) error { return nil }

func (s *flakySubscriber) Run(ctx context.Context, _ mqv2.Subscription, handler mqv2.Handler) error {
	if s.calls.Add(1) == 1 {
		return errors.New("temporary broker failure")
	}
	if err := handler(ctx, s.message); err != nil {
		return err
	}
	<-ctx.Done()
	return ctx.Err()
}

func TestConsumerRestartsAfterTransportFailure(t *testing.T) {
	payload, err := json.Marshal(NewDomainEvent("file.uploaded", "file-2", "file", map[string]string{"name": "b.txt"}))
	if err != nil {
		t.Fatal(err)
	}
	msg, err := mqv2.NewMessage("file.uploaded", payload)
	if err != nil {
		t.Fatal(err)
	}
	transport := &flakySubscriber{message: msg}
	consumer := NewEventConsumerService(transport, zap.NewNop())
	handler := &testEventHandler{received: make(chan BaseDomainEvent[json.RawMessage], 1)}
	if err := consumer.SubscribeHandler(handler); err != nil {
		t.Fatal(err)
	}
	if err := consumer.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = consumer.Stop(context.Background()) })
	select {
	case got := <-handler.received:
		if got.AggID != "file-2" || transport.calls.Load() < 2 {
			t.Fatalf("message was not delivered after reconnect: %+v", got)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("consumer did not retry after transport failure")
	}
}
