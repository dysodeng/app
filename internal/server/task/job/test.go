package job

import (
	"context"
	"log"
	"time"

	"github.com/dysodeng/app/internal/pkg/helper"

	"github.com/dysodeng/mq/message"
)

type TaskTestTask struct{}

func (task TaskTestTask) TopicKey() string {
	return "task_test"
}

func (task TaskTestTask) Handle(ctx context.Context, message *message.Message) error {
	log.Println("test task")
	log.Printf(
		`{"id": "%s", "topic": "%s", "payload": "%s", "headers": "%v", "delay": "%v", "create_at": "%s"}`,
		message.ID,
		message.Topic,
		helper.BytesToString(message.Payload),
		message.Headers,
		message.Delay,
		message.CreateAt.Format(time.DateTime),
	)
	return nil
}
