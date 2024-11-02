package cron

import (
	"log"
	"time"

	"github.com/dysodeng/app/internal/server/cron/job"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/server"
	"github.com/robfig/cron/v3"
)

// cronServer 定时任务服务
type cronServer struct {
	// jobs 任务注册表
	jobs map[string]job.JobInterface
	// schedule 任务调度器
	schedule *cron.Cron
}

func NewServer() server.Interface {
	jobServer := &cronServer{
		jobs: make(map[string]job.JobInterface),
	}
	return jobServer
}

// register 注册任务服务
func (cronJob *cronServer) register(jobs ...job.JobInterface) {
	for _, jobItem := range jobs {
		if _, ok := cronJob.jobs[jobItem.JobKey()]; !ok {
			cronJob.jobs[jobItem.JobKey()] = jobItem
		}
	}
}

func (cronJob *cronServer) Serve() {
	if !config.Server.Cron.Enabled {
		return
	}
	log.Println("start cron server...")

	cronJob.register()

	cronJob.schedule = cron.New(cron.WithLocation(time.Local))

	if len(cronJob.jobs) > 0 {
		for _, jobItem := range cronJob.jobs {
			id, _ := cronJob.schedule.AddFunc(jobItem.JobSpec(), jobItem.Handle)
			log.Printf("register cron job: %s, job time: %s, entry_id: %v\n", jobItem.JobKey(), jobItem.JobSpec(), id)
		}
	}

	cronJob.schedule.Start()
}

func (cronJob *cronServer) Shutdown() {
	if !config.Server.Cron.Enabled {
		return
	}
	log.Println("shutdown cron server...")
	if cronJob.schedule != nil {
		cronJob.schedule.Stop()
	}
}
