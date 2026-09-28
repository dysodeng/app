package mq

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dysodeng/app/internal/infrastructure/config"
	mqv2 "github.com/goairix/mq/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
)

func TestDisabledQueueUsesMemoryAdapter(t *testing.T) {
	client, err := Init(&config.Config{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Close(context.Background()) })

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sub := mqv2.Subscription{Topic: "test.event", Name: QueuePrefix}
	if err := client.Prepare(ctx, sub); err != nil {
		t.Fatal(err)
	}
	received := make(chan mqv2.Message, 1)
	done := make(chan error, 1)
	go func() {
		done <- client.Run(ctx, sub, func(_ context.Context, msg mqv2.Message) error {
			received <- msg
			return nil
		})
	}()
	msg, err := mqv2.NewMessage(sub.Topic, []byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Publish(ctx, msg); err != nil {
		t.Fatal(err)
	}
	select {
	case got := <-received:
		if got.ID != msg.ID || string(got.Payload) != "hello" {
			t.Fatalf("unexpected message: %+v", got)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for message")
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Run returned %v, want context cancellation", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Run did not stop after cancellation")
	}
}

func TestEnabledQueueRejectsUnknownDriver(t *testing.T) {
	_, err := Init(&config.Config{MessageQueue: config.MessageQueue{Enabled: true, Driver: "unknown"}})
	if err == nil {
		t.Fatal("unknown MQ driver should fail initialization")
	}
}

func TestClientPropagatesTraceContext(t *testing.T) {
	client, err := Init(&config.Config{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Close(context.Background()) })
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{1, 2, 3},
		SpanID:     trace.SpanID{4, 5, 6},
		TraceFlags: trace.FlagsSampled,
	})
	ctx := trace.ContextWithSpanContext(context.Background(), spanContext)
	sub := mqv2.Subscription{Topic: "trace.event", Name: QueuePrefix}
	received := make(chan trace.SpanContext, 1)
	runCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- client.Run(runCtx, sub, func(ctx context.Context, _ mqv2.Message) error {
			received <- trace.SpanContextFromContext(ctx)
			return nil
		})
	}()
	msg, err := mqv2.NewMessage(sub.Topic, []byte("traced"))
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Publish(ctx, msg); err != nil {
		t.Fatal(err)
	}
	select {
	case got := <-received:
		if got.TraceID() != spanContext.TraceID() {
			t.Fatalf("trace ID = %s, want %s", got.TraceID(), spanContext.TraceID())
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for traced message")
	}
	cancel()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("trace consumer did not stop")
	}
}

type reconnectTestBackend struct {
	dead  atomic.Bool
	fail  bool
	calls atomic.Int32
}

func (b *reconnectTestBackend) Publish(context.Context, mqv2.Message) error {
	b.calls.Add(1)
	if b.fail {
		b.dead.Store(true)
		return errors.New("connection lost")
	}
	return nil
}

func (b *reconnectTestBackend) Run(ctx context.Context, _ mqv2.Subscription, _ mqv2.Handler) error {
	<-ctx.Done()
	return ctx.Err()
}

func (b *reconnectTestBackend) Prepare(context.Context, mqv2.Subscription) error { return nil }

func (b *reconnectTestBackend) Close(context.Context) error { return nil }

func TestReconnectDoesNotRepublishUnknownOutcome(t *testing.T) {
	first := &reconnectTestBackend{fail: true}
	second := &reconnectTestBackend{}
	attempts := 0
	queue, err := newReconnectingBackend(func() (*connectionState, error) {
		attempts++
		backend := second
		if attempts == 1 {
			backend = first
		}
		return &connectionState{adapter: backend, isClosed: backend.dead.Load}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = queue.Close(context.Background()) })
	msg, err := mqv2.NewMessage("test.event", []byte("payload"))
	if err != nil {
		t.Fatal(err)
	}
	if err := queue.Publish(context.Background(), msg); err == nil {
		t.Fatal("first publish should report uncertain result")
	}
	if first.calls.Load() != 1 || second.calls.Load() != 0 {
		t.Fatal("failed publish was automatically retried")
	}
	if err := queue.Publish(context.Background(), msg); err != nil {
		t.Fatal(err)
	}
	if attempts != 2 || second.calls.Load() != 1 {
		t.Fatalf("reconnect attempts=%d, second publishes=%d", attempts, second.calls.Load())
	}
}

func TestRabbitURLPreservesCredentialsAndVirtualHost(t *testing.T) {
	cfg := &config.Config{}
	cfg.MessageQueue.Amqp.Host = "rabbit.example"
	cfg.MessageQueue.Amqp.Port = "5672"
	cfg.MessageQueue.Amqp.Username = "user@company"
	cfg.MessageQueue.Amqp.Password = "p@ss:/word"
	cfg.MessageQueue.Amqp.Vhost = "/team/prod"
	parsed, err := amqp.ParseURI(rabbitURL(cfg))
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Host != "rabbit.example" || parsed.Port != 5672 || parsed.Username != "user@company" || parsed.Password != "p@ss:/word" || parsed.Vhost != "/team/prod" {
		t.Fatalf("RabbitMQ URL lost configuration: %+v", parsed)
	}
}

func TestRedisMQConnectionModes(t *testing.T) {
	tests := []struct {
		name    string
		config  config.RedisItem
		wantErr bool
		check   func(*testing.T, redis.UniversalClient)
	}{
		{
			name: "standalone defaults",
			check: func(t *testing.T, client redis.UniversalClient) {
				t.Helper()
				standalone, ok := client.(*redis.Client)
				if !ok || standalone.Options().Addr != "127.0.0.1:6379" {
					t.Fatalf("unexpected standalone client: %T", client)
				}
			},
		},
		{
			name:   "cluster",
			config: config.RedisItem{Mode: "cluster", Cluster: config.ClusterConfig{Addrs: []string{"redis-1:6379", "redis-2:6379"}}},
			check: func(t *testing.T, client redis.UniversalClient) {
				t.Helper()
				cluster, ok := client.(*redis.ClusterClient)
				if !ok || len(cluster.Options().Addrs) != 2 {
					t.Fatalf("unexpected cluster client: %T", client)
				}
			},
		},
		{name: "cluster missing addresses", config: config.RedisItem{Mode: "cluster"}, wantErr: true},
		{name: "sentinel missing master", config: config.RedisItem{Mode: "sentinel"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := newRedisClient(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected configuration error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = client.Close() })
			tt.check(t, client)
		})
	}
}
