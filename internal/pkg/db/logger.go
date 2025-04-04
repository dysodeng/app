package db

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/dysodeng/app/internal/pkg/logger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	gormLogger "gorm.io/gorm/logger"
)

const defaultSlowThresholdTime = 300 * time.Millisecond

// GormLogger gorm日志
type GormLogger struct {
	SlowThreshold time.Duration
	LogLevel      gormLogger.LogLevel
	_zapLogger    *zap.Logger
}

func NewGormLogger() gormLogger.Interface {
	return &GormLogger{
		SlowThreshold: defaultSlowThresholdTime,
		LogLevel:      gormLogger.Warn,
		_zapLogger:    logger.WithOptions(zap.WithCaller(true)),
	}
}

func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return &GormLogger{
		SlowThreshold: defaultSlowThresholdTime,
		LogLevel:      level,
		_zapLogger:    logger.WithOptions(zap.WithCaller(true)),
	}
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Info {
		var fields []logger.Field
		if len(data)%2 == 0 {
			for i, datum := range data {
				index := i + 1
				if index%2 != 0 {
					fields = append(fields, logger.Field{Key: datum.(string), Value: data[index]})
				}
			}
		}
		logger.Info(ctx, msg, fields...)
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Warn {
		var fields []logger.Field
		if len(data)%2 == 0 {
			for i, datum := range data {
				index := i + 1
				if index%2 != 0 {
					fields = append(fields, logger.Field{Key: datum.(string), Value: data[index]})
				}
			}
		}
		logger.Warn(ctx, msg, fields...)
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Error {
		var fields []logger.Field
		if len(data)%2 == 0 {
			for i, datum := range data {
				index := i + 1
				if index%2 != 0 {
					fields = append(fields, logger.Field{Key: datum.(string), Value: data[index]})
				}
			}
		}
		logger.Error(ctx, msg, fields...)
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}
	duration := time.Since(begin).Milliseconds()
	sql, rows := fc()
	_, file, line, _ := runtime.Caller(3)

	traceFields := l.trace(ctx)

	if err != nil {
		// 错误日志
		fields := []zap.Field{
			zap.Any("error", err),
			zap.Any("stack_file", fmt.Sprintf("%s:%d", file, line)),
			zap.Any("rows", rows),
			zap.Any("duration", fmt.Sprintf("%dms", duration)),
		}
		fields = append(fields, traceFields...)
		l._zapLogger.Error(
			fmt.Sprintf("SQL ERROR: %s", sql),
			fields...,
		)
	} else {
		// 慢查询日志
		if duration > l.SlowThreshold.Milliseconds() {
			fields := []zap.Field{
				zap.Any("stack_file", fmt.Sprintf("%s:%d", file, line)),
				zap.Any("rows", rows),
				zap.Any("duration", fmt.Sprintf("%dms", duration)),
			}
			fields = append(fields, traceFields...)
			l._zapLogger.Warn(
				fmt.Sprintf("SQL SLOW: %s", sql),
				fields...,
			)
		} else {
			if l.LogLevel == gormLogger.Info {
				fields := []zap.Field{
					zap.Any("stack_file", fmt.Sprintf("%s:%d", file, line)),
					zap.Any("rows", rows),
					zap.Any("duration", fmt.Sprintf("%dms", duration)),
				}
				fields = append(fields, traceFields...)
				l._zapLogger.Debug(
					fmt.Sprintf("SQL DEBUG: %s", sql),
					fields...,
				)
			}
		}
	}
}

func (l *GormLogger) trace(ctx context.Context) []zap.Field {
	span := trace.SpanFromContext(ctx)
	var fields []zap.Field
	if span.SpanContext().HasTraceID() {
		fields = append(fields, zap.Any("trace_id", span.SpanContext().TraceID().String()))
	}
	if span.SpanContext().HasSpanID() {
		fields = append(fields, zap.Any("span_id", span.SpanContext().SpanID().String()))
	}
	if spanIns, ok := span.(sdktrace.ReadWriteSpan); ok {
		fields = append(fields, zap.Any("span_name", spanIns.Name()))
		if spanIns.Parent().HasSpanID() {
			fields = append(fields, zap.Any("parent_span_id", spanIns.Parent().SpanID().String()))
		}
	}
	return fields
}
