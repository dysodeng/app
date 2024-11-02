package task

import (
	"log"

	"github.com/dysodeng/app/internal/pkg/mq"
	"github.com/dysodeng/app/internal/server/task/job"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/server"
	"github.com/dysodeng/mq/contract"
	"github.com/pkg/errors"
)

// taskServer MessageQueue消费者服务，支持即时消息和延迟消息的消费
type taskServer struct {
	jobs            map[string]job.Interface
	jobConsumerList []contract.Consumer
}

func NewServer() server.Interface {
	ts := &taskServer{}
	return ts
}

// register 注册任务
func (server *taskServer) register(jobs ...job.Interface) {
	if server.jobs == nil {
		server.jobs = make(map[string]job.Interface)
	}
	for _, jobItem := range jobs {
		if _, ok := server.jobs[jobItem.QueueKey()]; !ok {
			server.jobs[jobItem.QueueKey()] = jobItem
		}
	}
}

func (server *taskServer) Serve() {
	if !config.Server.Task.Enabled {
		return
	}
	log.Println("start task server...")

	server.register(
		job.TaskTestTask{},
	)

	for _, jobItem := range server.jobs {
		go func(jobItem job.Interface) {
			consumer, err := mq.NewMessageQueueConsumer(jobItem.QueueKey())
			if err != nil {
				log.Printf("%+v", errors.Wrap(err, "消息队列任务创建失败"))
			}

			if jobItem.IsDelay() {
				err = consumer.DelayQueueConsume(jobItem)
			} else {
				err = consumer.QueueConsume(jobItem)
			}
			if err != nil {
				log.Printf("%+v", errors.Wrap(err, "消息队列消费者启动失败"))
			}

			server.jobConsumerList = append(server.jobConsumerList, consumer)
		}(jobItem)
	}
}

func (server *taskServer) Shutdown() {
	if !config.Server.Task.Enabled {
		return
	}
	log.Println("shutdown task server...")
	for _, consumer := range server.jobConsumerList {
		_ = consumer.Close()
	}
}
