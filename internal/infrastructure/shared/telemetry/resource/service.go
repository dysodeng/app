package resource

import "github.com/dysodeng/app/internal/infrastructure/config"

func ServiceName() string {
	name := config.GlobalConfig.App.Name
	if config.GlobalConfig.Monitor.ServiceName != "" {
		name = config.GlobalConfig.Monitor.ServiceName
	}
	return name
}
