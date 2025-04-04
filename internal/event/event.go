package event

import (
	"log"
	"time"

	"github.com/dysodeng/app/internal/event/listener"
	"github.com/dysodeng/app/internal/pkg/event"
	"github.com/dysodeng/app/internal/pkg/mq"
	"github.com/dysodeng/mq/contract"
)

// Logged 事件调度类型
const (
	Registered event.Dispatcher = "registered" // 注册成功
	Logged     event.Dispatcher = "logged"     // 登录成功
)

var bus event.Bus

func init() {
	queueConsumer, err := mq.NewMessageQueueConsumer(event.DispatchQueueKey)
	if err != nil {
		log.Fatalf("event queue consumer init fail: %+v", err)
	}
	queueProducer, err := mq.NewMessageQueueProducer(&contract.Pool{
		MinConn:     2,
		MaxConn:     2,
		MaxIdleConn: 2,
		IdleTimeout: time.Hour,
	})
	if err != nil {
		log.Fatalf("event queue producer init fail: %+v", err)
	}

	bus = event.New(
		event.WithEventQueue("", queueConsumer, queueProducer),
	)

	// 注册事件监听器
	bus.Register(Registered, &listener.Registered{})
	bus.Register(Logged, &listener.Logged{})
}

// Dispatch 事件调度
func Dispatch(dispatcher event.Dispatcher, data map[string]interface{}, opts ...event.DispatchOption) {
	bus.Dispatch(dispatcher, data, opts...)
}

func WithQueue() event.DispatchOption {
	return event.WithQueue()
}
