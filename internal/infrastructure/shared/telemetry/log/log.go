package log

import (
	"context"

	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/resource"
	"github.com/pkg/errors"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"

	"github.com/dysodeng/app/internal/infrastructure/config"
)

var otlpLoggerProvider *log.LoggerProvider

func Init(cfg *config.Config) error {
	ctx := context.Background()

	var opts []otlploghttp.Option
	if cfg.Monitor.Log.OtlpEnabled {
		if cfg.Monitor.Log.OtlpEndpoint == "" {
			return errors.New("log otel endpoint is empty")
		}
		opts = append(
			opts,
			otlploghttp.WithEndpoint(cfg.Monitor.Log.OtlpEndpoint),
			otlploghttp.WithInsecure(),
		)
	}

	exporter, err := otlploghttp.New(ctx, opts...)
	if err != nil {
		return err
	}

	processor := log.NewBatchProcessor(exporter)
	otlpLoggerProvider = log.NewLoggerProvider(
		log.WithResource(resource.Resource()),
		log.WithProcessor(processor),
	)

	global.SetLoggerProvider(otlpLoggerProvider)

	return nil
}

func Provider() *log.LoggerProvider {
	return otlpLoggerProvider
}
