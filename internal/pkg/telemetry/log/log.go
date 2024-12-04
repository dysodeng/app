package log

import (
	"context"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/pkg/telemetry"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var otlpLoggerProvider *log.LoggerProvider

func init() {
	if err := Init(); err != nil {
		panic(err)
	}
}

func Init() error {
	ctx := context.Background()

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(telemetry.ServiceName()),
			semconv.ServiceVersion(config.Monitor.ServiceVersion),
			attribute.String("env", config.App.Env.String()),
		),
	)
	if err != nil {
		return err
	}

	var opts []otlploghttp.Option
	if config.Monitor.Log.OtlpEnabled {
		if config.Monitor.Log.OtlpEndpoint == "" {
			return errors.New("log otel endpoint is empty")
		}
		opts = append(
			opts,
			otlploghttp.WithEndpoint(config.Monitor.Log.OtlpEndpoint),
			otlploghttp.WithInsecure(),
		)
	}

	exporter, err := otlploghttp.New(ctx, opts...)
	if err != nil {
		return err
	}

	processor := log.NewBatchProcessor(exporter)
	otlpLoggerProvider = log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(processor),
	)

	global.SetLoggerProvider(otlpLoggerProvider)

	return nil
}

func Provider() *log.LoggerProvider {
	return otlpLoggerProvider
}
