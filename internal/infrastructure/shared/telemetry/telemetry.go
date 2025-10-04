package telemetry

import (
	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/log"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/metrics"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/trace"
)

type Monitor struct{}

func InitMonitor(cfg *config.Config) (*Monitor, error) {
	if err := log.Init(cfg); err != nil {
		return nil, err
	}
	if err := metrics.Init(cfg); err != nil {
		return nil, err
	}
	if err := trace.Init(cfg); err != nil {
		return nil, err
	}
	return &Monitor{}, nil
}
