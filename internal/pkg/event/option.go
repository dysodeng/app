package event

import "github.com/dysodeng/mq/contract"

type Option func(*bus)

func WithEventQueue(queuePrefixKey string, queue contract.Consumer, queueProducer contract.Producer) Option {
	return func(b *bus) {
		b.withQueue = true
		b.queuePrefixKey = queuePrefixKey
		b.queue = queue
		b.queueProducer = queueProducer
	}
}

// DispatchQueueKey 事件调度队列key
const DispatchQueueKey = "event_dispatch"

// dispatchOption 事件调度选项
type dispatchOption struct {
	withQueue bool // 启用事件队列
}

type DispatchOption func(*dispatchOption)

func WithQueue() DispatchOption {
	return func(option *dispatchOption) {
		option.withQueue = true
	}
}
