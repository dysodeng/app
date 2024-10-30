package logger

import (
	"context"

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
	var fields []Field
	if ctx.Value("traceId") != nil {
		fields = append(fields, Field{Key: "traceId", Value: ctx.Value("traceId")})
	}
	if ctx.Value("spanId") != nil {
		fields = append(fields, Field{Key: "spanId", Value: ctx.Value("spanId")})
	}
	if ctx.Value("spanName") != nil {
		fields = append(fields, Field{Key: "spanName", Value: ctx.Value("spanName")})
	}
	if ctx.Value("parentSpanId") != nil {
		fields = append(fields, Field{Key: "parentSpanId", Value: ctx.Value("parentSpanId")})
	}
	if ctx.Value("parentSpanName") != nil {
		fields = append(fields, Field{Key: "parentSpanName", Value: ctx.Value("parentSpanName")})
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
