package mq

import (
	"context"
	"errors"
	"fmt"

	"github.com/goairix/mq/adapters/memory/v2"
	rabbitadapter "github.com/goairix/mq/adapters/rabbitmq/v2"
	redisadapter "github.com/goairix/mq/adapters/redis/v2"
	otelmq "github.com/goairix/mq/observability/otel/v2"
	mqv2 "github.com/goairix/mq/v2"

	"github.com/dysodeng/app/internal/infrastructure/config"
)

const QueuePrefix = "app-service"

type backend interface {
	mqv2.Publisher
	mqv2.Subscriber
}

// Client owns an MQ v2 adapter and its broker connection.
type Client struct {
	backend         backend
	publisher       mqv2.Publisher
	subscriber      mqv2.Subscriber
	closeConnection func() error
}

func newClient(backend backend, closeConnection func() error) (*Client, error) {
	instrumentation, err := otelmq.New(otelmq.Options{})
	if err != nil {
		_ = backend.Close(context.Background())
		if closeConnection != nil {
			_ = closeConnection()
		}
		return nil, fmt.Errorf("initialize MQ telemetry: %w", err)
	}
	return &Client{
		backend:         backend,
		publisher:       instrumentation.Publisher(backend),
		subscriber:      instrumentation.Subscriber(backend),
		closeConnection: closeConnection,
	}, nil
}

// Init creates the configured MQ v2 adapter.
func Init(cfg *config.Config) (*Client, error) {
	if cfg == nil {
		return nil, errors.New("MQ configuration is nil")
	}
	if !cfg.MessageQueue.Enabled || cfg.MessageQueue.Driver == "memory" {
		broker, err := memory.New(10000)
		if err != nil {
			return nil, fmt.Errorf("create memory MQ: %w", err)
		}
		return newClient(broker, nil)
	}

	switch cfg.MessageQueue.Driver {
	case "redis":
		redisClient, err := newRedisClient(cfg.Redis.MQ)
		if err != nil {
			return nil, err
		}
		adapter, err := redisadapter.New(redisClient, redisadapter.Options{Prefix: QueuePrefix})
		if err != nil {
			_ = redisClient.Close()
			return nil, fmt.Errorf("create Redis MQ adapter: %w", err)
		}
		return newClient(adapter, redisClient.Close)
	case "amqp":
		adapter, err := newReconnectingBackend(func() (*connectionState, error) {
			conn, err := newRabbitConnection(cfg)
			if err != nil {
				return nil, err
			}
			queue, err := rabbitadapter.New(conn, rabbitadapter.Options{Prefix: QueuePrefix})
			if err != nil {
				_ = conn.Close()
				return nil, fmt.Errorf("create RabbitMQ adapter: %w", err)
			}
			return &connectionState{
				adapter:         queue,
				isClosed:        conn.IsClosed,
				closeConnection: conn.Close,
			}, nil
		})
		if err != nil {
			return nil, err
		}
		return newClient(adapter, nil)
	default:
		return nil, fmt.Errorf("unsupported MQ driver %q", cfg.MessageQueue.Driver)
	}
}

// Publish sends a message and waits for broker confirmation.
func (c *Client) Publish(ctx context.Context, message mqv2.Message) error {
	return c.publisher.Publish(ctx, message)
}

// Prepare creates a durable subscription when the adapter supports it.
func (c *Client) Prepare(ctx context.Context, sub mqv2.Subscription) error {
	if err := sub.Validate(); err != nil {
		return err
	}
	if preparer, ok := c.backend.(interface {
		Prepare(context.Context, mqv2.Subscription) error
	}); ok {
		return preparer.Prepare(ctx, sub)
	}
	return nil
}

// Run consumes a subscription until its context is canceled or the adapter fails.
func (c *Client) Run(ctx context.Context, sub mqv2.Subscription, handler mqv2.Handler) error {
	return c.subscriber.Run(ctx, sub, handler)
}

// Close shuts down the adapter before closing its broker connection.
func (c *Client) Close(ctx context.Context) error {
	adapterErr := c.backend.Close(ctx)
	var connectionErr error
	if c.closeConnection != nil {
		connectionErr = c.closeConnection()
	}
	return errors.Join(adapterErr, connectionErr)
}
