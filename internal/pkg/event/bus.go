package event

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/dysodeng/app/internal/pkg/mq"

	"github.com/dysodeng/mq/contract"
)

// Dispatcher 事件调度类型
type Dispatcher string

// Bus 事件总线接口
type Bus interface {
	// Register 注册事件监听器
	Register(dispatcher Dispatcher, listeners ...Listener)
	// Dispatch 事件调度
	Dispatch(dispatcher Dispatcher, data map[string]interface{}, opts ...DispatchOption)
}

// bus 事件总线
type bus struct {
	events         sync.Map
	withQueue      bool
	queue          contract.Consumer
	queuePrefixKey string
	queueProducer  contract.Producer
}

// New 初始化事件总线
func New(opts ...Option) Bus {
	b := &bus{}

	for _, opt := range opts {
		opt(b)
	}

	if b.withQueue {
		queue, err := mq.NewMessageQueueConsumer(b.queuePrefixKey + eventDispatchQueueKey)
		if err != nil {
			log.Fatalf("%+v", err)
		}

		b.queue = queue

		b.queueProducer, err = mq.NewMessageQueueProducer(&contract.Pool{
			MinConn:     2,
			MaxConn:     2,
			MaxIdleConn: 2,
			IdleTimeout: time.Hour,
		})
		if err != nil {
			log.Fatalf("%+v", err)
		}

		go b.queueConsume()
	}

	return b
}

// Register 注册事件监听器
func (bus *bus) Register(dispatcher Dispatcher, listeners ...Listener) {
	e := &event{}

	if len(listeners) > 0 {
		for _, listener := range listeners {
			e.register(listener)
		}
	}

	bus.events.Store(dispatcher, e)
}

// Dispatch 事件调度
func (bus *bus) Dispatch(dispatcher Dispatcher, data map[string]interface{}, opts ...DispatchOption) {
	o := &dispatchOption{}
	for _, opt := range opts {
		opt(o)
	}

	if e, ok := bus.events.Load(dispatcher); ok {
		if bus.withQueue && o.withQueue {
			queueData, _ := json.Marshal(eventQueueData{Dispatcher: string(dispatcher), Data: data})
			_, err := bus.queueProducer.QueuePublish(mq.QueueKey(bus.queuePrefixKey+eventDispatchQueueKey), string(queueData))
			if err != nil {
				log.Printf("%+v", err)
			}
		} else {
			e.(*event).dispatch(data)
		}
	} else {
		log.Printf("事件调度器[%s]未注册", dispatcher)
	}
}
