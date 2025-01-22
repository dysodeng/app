package task

import (
	"log"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/pkg/mq"
	"github.com/dysodeng/app/internal/server"
	"github.com/dysodeng/app/internal/server/task/job"
	"github.com/dysodeng/mq/contract"
	"github.com/pkg/errors"
)

// taskServer MessageQueue消费者服务，支持即时消息和延迟消息的消费
type taskServer struct {
	jobs            map[string]job.Handler
	jobConsumerList []contract.Consumer
}

func NewServer() server.Server {
	ts := &taskServer{}
	return ts
}

func (server *taskServer) IsEnabled() bool {
	return config.Server.Task.Enabled
}

// register 注册任务
func (server *taskServer) register(jobs ...job.Handler) {
	if server.jobs == nil {
		server.jobs = make(map[string]job.Handler)
	}
	for _, jobItem := range jobs {
		if _, ok := server.jobs[jobItem.QueueKey()]; !ok {
			server.jobs[jobItem.QueueKey()] = jobItem
		}
	}
}

func (server *taskServer) Serve() {
	log.Println("start task server...")
	server.register(
		job.TaskTestTask{},
	)
	for _, jobHandler := range server.jobs {
		go func(jobHandler job.Handler) {
			consumer, err := mq.NewMessageQueueConsumer(jobHandler.QueueKey())
			if err != nil {
				log.Printf("%+v", errors.Wrap(err, "消息队列任务创建失败"))
			}

			if jobHandler.IsDelay() {
				err = consumer.DelayQueueConsume(jobHandler)
			} else {
				err = consumer.QueueConsume(jobHandler)
			}
			if err != nil {
				log.Printf("%+v", errors.Wrap(err, "消息队列消费者启动失败"))
			}

			server.jobConsumerList = append(server.jobConsumerList, consumer)
		}(jobHandler)
	}
}

func (server *taskServer) Shutdown() {
	log.Println("shutdown task server...")
	for _, consumer := range server.jobConsumerList {
		_ = consumer.Close()
	}
}
