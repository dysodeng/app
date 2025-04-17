package trace

import (
	"context"

	"github.com/dysodeng/app/internal/config"
	"github.com/dysodeng/app/internal/pkg/telemetry"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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

func Init() error {
	traceCtx = context.Background()
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
	tpOpts := []sdktrace.TracerProviderOption{
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
	}
	if config.Monitor.Tracer.OtlpEnabled {
		if config.Monitor.Tracer.OtlpEndpoint == "" {
			return errors.New("tracer otel endpoint is empty")
		}
		exp, err := otlptracehttp.New(
			traceCtx,
			otlptracehttp.WithEndpointURL(config.Monitor.Tracer.OtlpEndpoint),
		)
		if err != nil {
			return err
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
		telemetry.ServiceName(),
		trace.WithInstrumentationVersion(config.Monitor.ServiceVersion),
	)

	return nil
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

// ParseContextTraceId 从上下文获取 traceId
func ParseContextTraceId(ctx context.Context) string {
	var traceId string
	if ctx.Value("X-Trace-Id") != nil {
		traceId = ctx.Value("X-Trace-Id").(string)
	} else {
		span := trace.SpanFromContext(ctx)
		if span.SpanContext().HasTraceID() {
			traceId = span.SpanContext().TraceID().String()
		}
	}
	return traceId
}

func Gin(ctx *gin.Context) context.Context {
	return ctx.Request.Context()
}
