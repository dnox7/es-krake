package hook

import (
	"context"
	"errors"
	"pech/es-krake/pkg/log"
	"pech/es-krake/pkg/log/gorm"
	"time"

	"gorm.io/gorm/logger"
)

func DefaultGormLogger() logger.Interface {
	return NewGormLogger(*gorm.DefaultConfig())
}

type gormLogger struct {
	LogLevel logger.LogLevel
	conf     logger.Config
	logger   *log.Logger
}

func NewGormLogger(conf logger.Config) logger.Interface {
	return &gormLogger{
		LogLevel: conf.LogLevel,
		conf:     conf,
		logger:   log.With("service", "database"),
	}
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	log := *l
	log.LogLevel = level
	return &log
}

func (l *gormLogger) Info(ctx context.Context, fmt string, args ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.logger.Info(ctx, fmt, args...)
	}
}

func (l *gormLogger) Warn(ctx context.Context, fmt string, args ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.logger.Warn(ctx, fmt, args...)
	}
}

func (l *gormLogger) Error(ctx context.Context, fmt string, args ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.logger.Error(ctx, fmt, args...)
	}
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fn func() (sql string, rowsAffected int64), err error) {
	sql, rows := fn()
	duration := time.Since(begin)
	logEntry := l.logger.
		With("duration", duration.String()).
		With("sql", sql).
		With("rows", rows)

	if err == nil {
		if duration.Milliseconds() > l.conf.SlowThreshold.Milliseconds() {
			logEntry.Warn(ctx, "Performed SLOW SQL Query")
		} else {
			logEntry.Info(ctx, "Performed SQL Query")
		}
	} else {
		logEntry = logEntry.With("error", err)
		if errors.Is(err, logger.ErrRecordNotFound) {
			logEntry.Info(ctx, "Performed SQL Query")
		} else {
			logEntry.Error(ctx, "SQL Query Failed")
		}
	}
}
