package db

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"go.uber.org/zap"

	"github.com/dysodeng/app/internal/pkg/logger"
	gormLogger "gorm.io/gorm/logger"
)

// GormLogger gorm日志
type GormLogger struct {
	SlowThreshold time.Duration
	LogLevel      gormLogger.LogLevel
}

func NewGormLogger() gormLogger.Interface {
	return &GormLogger{
		SlowThreshold: 300 * time.Millisecond,
		LogLevel:      gormLogger.Info,
	}
}

func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return &GormLogger{
		SlowThreshold: 300 * time.Millisecond,
		LogLevel:      level,
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
	duration := time.Since(begin).Milliseconds()
	sql, rows := fc()
	_, file, line, _ := runtime.Caller(3)

	traceFields := l.trace(ctx)
	dbLogger := logger.WithOptions(zap.WithCaller(false))

	if err != nil {
		// 错误日志
		fields := []zap.Field{
			zap.Any("error", err),
			zap.Any("file", fmt.Sprintf("%s:%d", file, line)),
			zap.Any("rows", rows),
			zap.Any("duration", fmt.Sprintf("%dms", duration)),
		}
		for _, field := range traceFields {
			fields = append(fields, field)
		}
		dbLogger.Error(
			fmt.Sprintf("SQL ERROR: sql( %s )", sql),
			fields...,
		)
	} else {
		// 慢查询日志
		if duration > l.SlowThreshold.Milliseconds() {
			fields := []zap.Field{
				zap.Any("file", fmt.Sprintf("%s:%d", file, line)),
				zap.Any("rows", rows),
				zap.Any("duration", fmt.Sprintf("%dms", duration)),
			}
			for _, field := range traceFields {
				fields = append(fields, field)
			}
			dbLogger.Warn(
				fmt.Sprintf("SQL SLOW: sql( %s )", sql),
				fields...,
			)
		} else {
			fields := []zap.Field{
				zap.Any("file", fmt.Sprintf("%s:%d", file, line)),
				zap.Any("rows", rows),
				zap.Any("duration", fmt.Sprintf("%dms", duration)),
			}
			for _, field := range traceFields {
				fields = append(fields, field)
			}
			dbLogger.Warn(
				fmt.Sprintf("SQL DEBUG: sql( %s )", sql),
				fields...,
			)
		}
	}
}

func (l *GormLogger) trace(ctx context.Context) []zap.Field {
	var fields []zap.Field
	if ctx.Value("traceId") != nil {
		fields = append(fields, zap.Any("traceId", ctx.Value("traceId")))
	}
	if ctx.Value("spanId") != nil {
		fields = append(fields, zap.Any("spanId", ctx.Value("spanId")))
	}
	if ctx.Value("spanName") != nil {
		fields = append(fields, zap.Any("spanName", ctx.Value("spanName")))
	}
	if ctx.Value("parentSpanId") != nil {
		fields = append(fields, zap.Any("parentSpanId", ctx.Value("parentSpanId")))
	}
	return fields
}
