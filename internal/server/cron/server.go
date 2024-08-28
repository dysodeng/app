package cron

import (
	"log"
	"time"

	"github.com/dysodeng/app/internal/config"
	cronJobIface "github.com/dysodeng/app/internal/cron"
	"github.com/dysodeng/app/internal/server"
	"github.com/robfig/cron/v3"
)

// cronJob 定时任务
type cronJob struct {
	// jobs 任务注册表
	jobs map[string]cronJobIface.JobInterface
	// schedule 任务调度器
	schedule *cron.Cron
}

func NewServer() server.Interface {
	jobServer := &cronJob{
		jobs: make(map[string]cronJobIface.JobInterface),
	}
	jobServer.register()
	return jobServer
}

// register 注册任务服务
func (cronJob *cronJob) register(jobs ...cronJobIface.JobInterface) {
	for _, job := range jobs {
		if _, ok := cronJob.jobs[job.JobKey()]; !ok {
			cronJob.jobs[job.JobKey()] = job
		}
	}
}

func (cronJob *cronJob) Serve() {
	if !config.Server.Cron.Enabled {
		return
	}
	log.Println("start cron server...")

	cronJob.schedule = cron.New(cron.WithLocation(time.Local))

	if len(cronJob.jobs) > 0 {
		for _, job := range cronJob.jobs {
			id, _ := cronJob.schedule.AddFunc(job.JobSpec(), job.Handle)
			log.Printf("register cron job: %s, job time: %s, entry_id: %v\n", job.JobKey(), job.JobSpec(), id)
		}
	}

	cronJob.schedule.Start()
}

func (cronJob *cronJob) Shutdown() {
	if !config.Server.Cron.Enabled {
		return
	}
	log.Println("shutdown cron server...")
	if cronJob.schedule != nil {
		cronJob.schedule.Stop()
	}
}
