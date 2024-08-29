package event

import (
	"encoding/json"
	"log"

	"github.com/dysodeng/mq/message"
)

// eventQueueData 事件队列数据
type eventQueueData struct {
	Dispatcher string                 `json:"dispatcher"`
	Data       map[string]interface{} `json:"data"`
}

// eventQueueHandler 事件队列消息处理器
type eventQueueHandler struct {
	bus *bus
}

func (handler *eventQueueHandler) Handle(message message.Message) error {
	body := message.Body()
	var data eventQueueData

	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		return err
	}

	if e, ok := handler.bus.events.Load(Dispatcher(data.Dispatcher)); ok {
		e.(*event).dispatch(data.Data)
	}

	return nil
}

// queueConsume 事件队列消费者
func (bus *bus) queueConsume() {
	if bus.withQueue {
		err := bus.queue.QueueConsume(&eventQueueHandler{bus: bus})
		if err != nil {
			log.Printf("%+v", err)
		}
	}
}
