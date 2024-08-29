package event

import (
	"github.com/dysodeng/app/internal/event/listener"
	"github.com/dysodeng/app/internal/pkg/event"
)

// Logged 事件调度类型
const (
	Registered event.Dispatcher = "registered" // 注册成功
	Logged     event.Dispatcher = "logged"     // 登录成功
)

var bus event.Bus

func init() {
	bus = event.New(
		event.WithEventQueue(""),
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
