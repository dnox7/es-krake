package gormlog

import (
	"context"
	"errors"
	"time"

	"github.com/dpe27/es-krake/pkg/log"
	"gorm.io/gorm/logger"
)

type gormLogKey string

const (
	rowsKey     gormLogKey = "rows"
	durationKey gormLogKey = "duration"
	sqlKey      gormLogKey = "sql"
)

type gormLogger struct {
	Level  logger.LogLevel
	cfg    logger.Config
	logger *log.Logger
}

// DefaultGormLogger creates a new Logger for GORM
func DefaultGormLogger() logger.Interface {
	return NewGormLogger(logger.Config{
		IgnoreRecordNotFoundError: true,
		SlowThreshold:             2000 * time.Millisecond,
		LogLevel:                  logger.Warn,
	})
}

// NewGormLogger creates a new Logger for GORM with detailed configurations.
func NewGormLogger(config logger.Config) logger.Interface {
	return &gormLogger{
		Level:  config.LogLevel,
		cfg:    config,
		logger: log.With("service", "postgres"),
	}
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	log := *l
	log.Level = level
	return &log
}

func (l *gormLogger) Info(ctx context.Context, fmt string, args ...interface{}) {
	if l.Level >= logger.Info {
		l.logger.Info(ctx, fmt, args...)
	}
}

func (l *gormLogger) Warn(ctx context.Context, fmt string, args ...interface{}) {
	if l.Level >= logger.Warn {
		l.logger.Warn(ctx, fmt, args...)
	}
}

func (l *gormLogger) Error(ctx context.Context, fmt string, args ...interface{}) {
	if l.Level >= logger.Error {
		l.logger.Error(ctx, fmt, args...)
	}
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time,
	fc func() (sql string, rowsAffected int64), err error,
) {
	sql, rows := fc()
	duration := time.Since(begin)
	logEntry := l.logger.
		With("duration", duration.String()).
		With("sql", sql).
		With("keys", rows)

	if err == nil {
		if duration.Milliseconds() > l.cfg.SlowThreshold.Milliseconds() {
			logEntry.Warn(ctx, "Performed SLOW SQL Query")
		} else {
			logEntry.Debug(ctx, "Performed SQL Query")
		}
	} else {
		logEntry = logEntry.With("error", err)
		if errors.Is(err, logger.ErrRecordNotFound) {
			logEntry.Debug(ctx, "Performed SQL Query")
		} else {
			logEntry.Error(ctx, "SQL Query failed")
		}
	}
}
