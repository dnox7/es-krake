package hook

import (
	"context"
	"errors"
	"pech/es-krake/pkg/logging"
	"pech/es-krake/pkg/logging/gorm"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

func DefaultGormLogger(base logging.BaseLogger) logger.Interface {
	return NewGormLogger(base, *gorm.DefaultConfig())
}

type gormLogger struct {
	LogLevel logger.LogLevel
	conf     logger.Config
	logger   logging.BaseLogger
}

func NewGormLogger(base logging.BaseLogger, conf logger.Config) logger.Interface {
	return &gormLogger{
		LogLevel: conf.LogLevel,
		conf:     conf,
		logger:   base.WithField("service", "database"),
	}
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	log := *l
	log.LogLevel = level
	return &log
}

func (l *gormLogger) Info(ctx context.Context, fmt string, args ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.logger.WithContext(ctx).Infof(fmt, args...)
	}
}

func (l *gormLogger) Warn(ctx context.Context, fmt string, args ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.logger.WithContext(ctx).Warnf(fmt, args...)
	}
}

func (l *gormLogger) Error(ctx context.Context, fmt string, args ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.logger.WithContext(ctx).Errorf(fmt, args...)
	}
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fn func() (sql string, rowsAffected int64), err error) {
	sql, rows := fn()
	duration := time.Since(begin)
	logEntry := l.logger.
		WithContext(ctx).
		WithFields(logrus.Fields{
			"duration": duration.String(),
			"sql":      sql,
			"rows":     rows,
		})

	if err == nil {
		if duration.Milliseconds() > l.conf.SlowThreshold.Milliseconds() {
			logEntry.Warn("Performed SLOW SQL Query")
		} else {
			logEntry.Trace("Performed SQL Query")
		}
	} else {
		logEntry = logEntry.WithField("error", err)
		if errors.Is(err, logger.ErrRecordNotFound) {
			logEntry.Trace("Performed SQL Query")
		} else {
			logEntry.Error("SQL Query Failed")
		}
	}
}
