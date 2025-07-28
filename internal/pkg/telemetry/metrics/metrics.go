package metrics

import (
	"context"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/pkg/telemetry"
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

// 初始化指标 meterProvider
func init() {
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
		panic(err)
	}

	mpOpts := []sdkmetric.Option{
		sdkmetric.WithResource(res),
	}
	if config.Monitor.Metrics.OtlpEnabled {
		if config.Monitor.Metrics.OtlpEndpoint == "" {
			panic(errors.New("missing metrics endpoint"))
		}
		exp, err := otlpmetrichttp.New(
			context.Background(),
			otlpmetrichttp.WithEndpointURL(config.Monitor.Tracer.OtlpEndpoint),
		)
		if err != nil {
			panic(err)
		}
		mpOpts = append(mpOpts, sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp)))
	}

	meterProvider := sdkmetric.NewMeterProvider(mpOpts...)
	otel.SetMeterProvider(meterProvider)

	meter = otel.Meter(config.App.Name)
}

func Init() error {
	return nil
}

func Meter() metric.Meter {
	return meter
}
