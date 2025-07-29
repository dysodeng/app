package task

import (
	"context"
	"log"
	"time"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/pkg/logger"
	"github.com/dysodeng/app/internal/pkg/mq"
	"github.com/dysodeng/app/internal/server"
	"github.com/dysodeng/app/internal/server/task/job"
	"github.com/dysodeng/mq/contract"
	"github.com/pkg/errors"
)

// taskServer MessageQueue消费者服务，支持即时消息和延迟消息的消费
type taskServer struct {
	jobs map[string]job.Handler
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
		if _, ok := server.jobs[jobItem.TopicKey()]; !ok {
			server.jobs[jobItem.TopicKey()] = jobItem
		}
	}
}

func (server *taskServer) Serve() {
	log.Println("start task server...")
	server.register(
		job.TaskTestTask{},
	)

	ctx := context.Background()

	for _, jobHandler := range server.jobs {
		go func(jobHandler job.Handler) {
			// 创建中间件链
			middlewareChain := contract.NewMiddlewareChain(
				contract.LoggingMiddleware(logger.ZapLogger()),
				contract.TimeoutMiddleware(30*time.Second),
				contract.RetryMiddleware(3, time.Second),
			)

			handler := middlewareChain.Apply(jobHandler.Handle)

			err := mq.Consumer().Subscribe(ctx, jobHandler.TopicKey(), handler)
			if err != nil {
				log.Printf("%+v", errors.Wrap(err, "消息队列消费者启动失败"))
			}
		}(jobHandler)
	}
}

func (server *taskServer) Shutdown() {
	log.Println("shutdown task server...")
	_ = mq.Consumer().Close()
	_ = mq.Instance().Close()
}
