package metrics

import (
	"context"

	"github.com/dysodeng/app/internal/config"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var meter metric.Meter

// MetricProviderInit 初始化指标 meterProvider
func MetricProviderInit() error {
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName()),
			semconv.ServiceVersion(config.Monitor.ServiceVersion),
			attribute.String("env", config.App.Env.String()),
		),
	)
	if err != nil {
		return err
	}

	mpOpts := []sdkmetric.Option{
		sdkmetric.WithResource(res),
	}
	if config.Monitor.Metrics.OtlpEnabled {
		if config.Monitor.Metrics.OtlpEndpoint == "" {
			return errors.New("metrics otel endpoint is empty")
		}
		exp, err := otlpmetrichttp.New(
			context.Background(),
			otlpmetrichttp.WithEndpointURL(config.Monitor.Tracer.OtlpEndpoint),
		)
		if err != nil {
			return err
		}
		mpOpts = append(mpOpts, sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp)))
	}

	meterProvider := sdkmetric.NewMeterProvider(mpOpts...)
	otel.SetMeterProvider(meterProvider)

	meter = otel.Meter(config.App.Name)

	return nil
}

func serviceName() string {
	name := config.App.Name
	if config.Monitor.ServiceName != "" {
		name = config.Monitor.ServiceName
	}
	return name
}

func Meter() metric.Meter {
	return meter
}
