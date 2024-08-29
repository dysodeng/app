package event

import "sync"

// event 事件
type event struct {
	listeners []Listener
}

// register 注册事件监听器
func (e *event) register(listener Listener) {
	e.listeners = append(e.listeners, listener)
}

// dispatch 事件调度
func (e *event) dispatch(data map[string]interface{}) {
	wg := sync.WaitGroup{}
	for _, listener := range e.listeners {
		wg.Add(1)
		go func(listener Listener) {
			defer wg.Done()
			listener.Handle(data)
		}(listener)
	}
	wg.Wait()
}

// Listener 事件监听器
type Listener interface {
	Handle(data map[string]interface{})
}
