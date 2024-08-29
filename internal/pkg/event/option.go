package event

type Option func(*bus)

func WithEventQueue(queuePrefixKey string) Option {
	return func(b *bus) {
		b.withQueue = true
		b.queuePrefixKey = queuePrefixKey
	}
}

// eventDispatchQueueKey 事件调度队列key
const eventDispatchQueueKey = "event_dispatch"

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
