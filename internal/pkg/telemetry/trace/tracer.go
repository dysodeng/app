package trace

import (
	"context"

	"github.com/dysodeng/app/internal/config"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// provider 全局 TracerProvider
var provider *sdktrace.TracerProvider

// Tracer 全局 Tracer
var tracer trace.Tracer

var traceCtx context.Context

func init() {
	traceCtx = context.Background()
	tpOpts := []sdktrace.TracerProviderOption{
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName()),
			semconv.ServiceVersion(config.Monitor.Tracer.ServiceVersion),
			attribute.String("env", config.App.Env.String()),
		)),
	}
	if config.Monitor.Tracer.OtlpEnabled {
		if config.Monitor.Tracer.OtlpEndpoint == "" {
			panic("tracer otel endpoint is empty")
		}
		exp, err := otlptracehttp.New(
			traceCtx,
			otlptracehttp.WithEndpointURL(config.Monitor.Tracer.OtlpEndpoint),
		)
		if err != nil {
			panic(err)
		}
		tpOpts = append(tpOpts, sdktrace.WithBatcher(exp))
	}

	provider = sdktrace.NewTracerProvider(tpOpts...)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	tracer = NewTracer(
		serviceName(),
	)
}

func serviceName() string {
	name := config.App.Name
	if config.Monitor.Tracer.ServiceName != "" {
		name = config.Monitor.Tracer.ServiceName
	}
	return name
}

func Context() context.Context {
	return traceCtx
}

func TracerProvider() *sdktrace.TracerProvider {
	return provider
}

func NewTracer(traceName string, opts ...trace.TracerOption) trace.Tracer {
	return provider.Tracer(traceName, opts...)
}

func Tracer() trace.Tracer {
	return tracer
}

func Gin(ctx *gin.Context) context.Context {
	return ctx.Request.Context()
}
