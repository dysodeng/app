package telemetry

import "github.com/dysodeng/app/internal/config"

func ServiceName() string {
	name := config.App.Name
	if config.Monitor.ServiceName != "" {
		name = config.Monitor.ServiceName
	}
	return name
}
