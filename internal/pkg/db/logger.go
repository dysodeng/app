package db

import (
	"context"
	"fmt"
	"time"

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
	if err != nil {
		// 错误日志
		logger.Error(
			ctx,
			fmt.Sprintf("SQL ERROR: sql( %s )", sql),
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "rows", Value: rows},
			logger.Field{Key: "duration", Value: fmt.Sprintf("%dms", duration)},
		)
	} else {
		// 慢查询日志
		if duration > l.SlowThreshold.Milliseconds() {
			logger.Warn(
				ctx,
				fmt.Sprintf("SQL SLOW: sql( %s )", sql),
				logger.Field{Key: "rows", Value: rows},
				logger.Field{Key: "duration", Value: fmt.Sprintf("%dms", duration)},
			)
		} else {
			logger.Debug(
				ctx,
				fmt.Sprintf("SQL DEBUG: sql( %s )", sql),
				logger.Field{Key: "rows", Value: rows},
				logger.Field{Key: "duration", Value: fmt.Sprintf("%dms", duration)},
			)
		}
	}
}
