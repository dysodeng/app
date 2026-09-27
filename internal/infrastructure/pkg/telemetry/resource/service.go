package resource

import "github.com/dysodeng/app/internal/infrastructure/config"

func ServiceName() string {
	return serviceName(config.Current())
}

func serviceName(cfg *config.Config) string {
	name := cfg.App.Name
	if cfg.Monitor.ServiceName != "" {
		name = cfg.Monitor.ServiceName
	}
	return name
}
