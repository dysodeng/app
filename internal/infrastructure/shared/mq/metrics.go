package mq

import (
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

type metricsObserver struct {
	meter  metric.Meter
	logger *zap.Logger
}

func (o *metricsObserver) GetMeter() metric.Meter {
	return o.meter
}

func (o *metricsObserver) GetLogger() *zap.Logger {
	return o.logger
}
