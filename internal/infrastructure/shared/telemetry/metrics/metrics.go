package metrics

import (
	"context"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"github.com/dysodeng/app/internal/infrastructure/config"
	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/resource"
)

var meter metric.Meter

// Init 初始化指标 meterProvider
func Init(cfg *config.Config) error {
	mpOpts := []sdkmetric.Option{
		sdkmetric.WithResource(resource.Resource()),
	}
	if cfg.Monitor.Metrics.OtlpEnabled {
		if cfg.Monitor.Metrics.OtlpEndpoint == "" {
			panic(errors.New("missing metrics endpoint"))
		}
		exp, err := otlpmetrichttp.New(
			context.Background(),
			otlpmetrichttp.WithEndpointURL(cfg.Monitor.Metrics.OtlpEndpoint),
		)
		if err != nil {
			panic(err)
		}

		var readerOpts []sdkmetric.PeriodicReaderOption
		interval := cfg.Monitor.Metrics.OtlpInterval
		if interval > 0 {
			readerOpts = append(readerOpts, sdkmetric.WithInterval(interval))
		}
		mpOpts = append(mpOpts, sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp, readerOpts...)))
	}

	meterProvider := sdkmetric.NewMeterProvider(mpOpts...)
	otel.SetMeterProvider(meterProvider)

	meter = otel.Meter(cfg.App.Name)

	return nil
}

func Meter() metric.Meter {
	return meter
}
