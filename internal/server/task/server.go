package task

import (
	"log"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/server"
	"github.com/dysodeng/app/internal/task"
	"github.com/dysodeng/app/pkg/mq"
	"github.com/dysodeng/mq/contract"
	"github.com/pkg/errors"
)

type taskServer struct {
	jobs            map[string]task.JobInterface
	jobConsumerList []contract.Consumer
}

func NewServer() server.Interface {
	ts := &taskServer{}
	ts.register()
	return ts
}

// register 注册任务
func (server *taskServer) register(jobs ...task.JobInterface) {
	if server.jobs == nil {
		server.jobs = make(map[string]task.JobInterface)
	}
	for _, job := range jobs {
		if _, ok := server.jobs[job.QueueKey()]; !ok {
			server.jobs[job.QueueKey()] = job
		}
	}
}

func (server *taskServer) Serve() {
	if !config.Server.Task.Enabled {
		return
	}
	log.Println("start task server...")
	for _, jobItem := range server.jobs {
		go func(jobItem task.JobInterface) {
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
