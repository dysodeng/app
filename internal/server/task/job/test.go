package job

import (
	"log"

	"github.com/dysodeng/mq/message"
)

type TaskTestTask struct{}

func (task TaskTestTask) QueueKey() string {
	return "task_test"
}

func (task TaskTestTask) IsDelay() bool {
	return false
}

func (task TaskTestTask) Handle(message message.Message) error {
	log.Println("test task")
	return nil
}
