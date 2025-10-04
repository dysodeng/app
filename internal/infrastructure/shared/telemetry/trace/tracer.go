package trace

import (
	"context"

	"github.com/dysodeng/app/internal/infrastructure/shared/telemetry/resource"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/dysodeng/app/internal/infrastructure/config"
)

// provider 全局 TracerProvider
var provider *sdktrace.TracerProvider

// Tracer 全局 Tracer
var tracer trace.Tracer

var traceCtx context.Context

func Init(cfg *config.Config) error {
	traceCtx = context.Background()
	tpOpts := []sdktrace.TracerProviderOption{
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.Resource()),
	}
	if cfg.Monitor.Tracer.OtlpEnabled {
		if cfg.Monitor.Tracer.OtlpEndpoint == "" {
			return errors.New("tracer otel endpoint is empty")
		}
		exp, err := otlptracehttp.New(
			traceCtx,
			otlptracehttp.WithEndpointURL(cfg.Monitor.Tracer.OtlpEndpoint),
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
		resource.ServiceName(),
		trace.WithInstrumentationVersion(cfg.Monitor.ServiceVersion),
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
	if tracer == nil {
		// 返回一个 no-op tracer 而不是 nil
		return otel.Tracer("fallback")
	}
	return tracer
}

func ContextWithSpan(ctx context.Context, span trace.Span) context.Context {
	return trace.ContextWithSpan(ctx, span)
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
