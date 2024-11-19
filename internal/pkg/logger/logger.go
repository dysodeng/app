package logger

import (
	"context"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	_logger *zap.Logger
}

type Field struct {
	Key   string
	Value interface{}
}

// ErrorField 错误堆栈
func ErrorField(err error) Field {
	return Field{Key: "error_field", Value: err}
}

var _logger *logger

func init() {
	newZapLogger()
	_logger = &logger{_logger: _zapLogger}
}

func (l *logger) log(ctx context.Context, level zapcore.Level, message string, fields ...Field) {
	fields = append(
		fields,
		l.trace(ctx)...,
	)

	zipFields := make([]zap.Field, 0, len(fields))
	for _, field := range fields {
		if field.Key == "error_field" {
			zipFields = append(zipFields, zap.Error(field.Value.(error)))
		} else {
			zipFields = append(zipFields, zap.Any(field.Key, field.Value))
		}
	}

	check := l._logger.Check(level, message)
	check.Write(zipFields...)
}

func (l *logger) trace(ctx context.Context) []Field {
	span := trace.SpanFromContext(ctx)
	var fields []Field
	if span.SpanContext().HasTraceID() {
		fields = append(fields, Field{Key: "traceId", Value: span.SpanContext().TraceID().String()})
	}
	if span.SpanContext().HasSpanID() {
		fields = append(fields, Field{Key: "spanId", Value: span.SpanContext().SpanID().String()})
	}
	if spanIns, ok := span.(sdktrace.ReadWriteSpan); ok {
		fields = append(fields, Field{Key: "spanName", Value: spanIns.Name()})
		if spanIns.Parent().HasSpanID() {
			fields = append(fields, Field{Key: "parentSpanId", Value: spanIns.Parent().SpanID().String()})
		}
	}
	return fields
}

func Debug(ctx context.Context, message string, fields ...Field) {
	_logger.log(ctx, zapcore.DebugLevel, message, fields...)
}

func Info(ctx context.Context, message string, fields ...Field) {
	_logger.log(ctx, zapcore.InfoLevel, message, fields...)
}

func Warn(ctx context.Context, message string, fields ...Field) {
	_logger.log(ctx, zapcore.WarnLevel, message, fields...)
}

func Error(ctx context.Context, message string, fields ...Field) {
	_logger.log(ctx, zapcore.ErrorLevel, message, fields...)
}

func Fatal(ctx context.Context, message string, fields ...Field) {
	_logger.log(ctx, zapcore.FatalLevel, message, fields...)
}

func Panic(ctx context.Context, message string, fields ...Field) {
	_logger.log(ctx, zapcore.PanicLevel, message, fields...)
}

func WithOptions(opts ...zap.Option) *zap.Logger {
	return _logger._logger.WithOptions(opts...)
}
