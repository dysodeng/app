package job

// Interface 任务接口
type Interface interface {
	// JobKey 任务key
	JobKey() string
	// JobSpec 任务执行时间规则
	JobSpec() string
	// Handle 任务执行方法
	Handle()
}
