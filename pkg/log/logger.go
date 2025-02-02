package log

import (
	"io"
	"log/slog"
	"pech/es-krake/config"
	"sync"
	"time"

	"golang.org/x/net/context"
)

type loggerCtxKey struct{}

const messageKey = "message"

type Logger struct {
	logger *slog.Logger
}

var (
	keys         []string
	logMapCtxKey = loggerCtxKey{}
)

// CtxWithValue: adds a key-val pair to the context in sync.Map for thread safely
// this value automatically added to the log record with defaultHandler
func CtxWithValue(ctx context.Context, key string, val interface{}) context.Context {
	m, ok := ctx.Value(logMapCtxKey).(*sync.Map)
	if !ok {
		m = &sync.Map{}
	}
	m.Store(key, val)
	return context.WithValue(ctx, logMapCtxKey, m)
}

func With(args ...any) *Logger {
	return &Logger{
		logger: slog.With(args...),
	}
}

func Group(key string, args ...any) slog.Attr {
	return slog.Group(key, args...)
}

// Initialize: initializes the logger with default handler
// if is deDebug is false, it sets the log level to [LevelDebug] otherwise [LevelInfo]
// it uses JSONHandler for logging and automatically adds the key-value pair to the log record
// using the CtxWithValue function.
// The fields in the keys are used to retrieve their values from the context and write them to the logger
func Initialize(ctx context.Context, w io.Writer, cfg *config.Config, keyInput []string) (*Logger, context.Context) {
	keys = append(keys, keyInput...)
	level := slog.LevelInfo
	if cfg.Log.Level == "DEBUG" {
		level = slog.LevelDebug
	}

	opt := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if v, ok := a.Value.Any().(time.Duration); ok {
				a.Value = slog.StringValue(v.String())
			}
			if a.Key != slog.MessageKey {
				return a
			}
			a.Key = messageKey
			return a
		},
	}

	h := &customHandler{}
	if cfg.Log.Format == "TEXT" {
		h.Handler = slog.NewTextHandler(w, opt)
	} else {
		h.Handler = slog.NewJSONHandler(w, opt)
	}

	l := slog.New(h)
	slog.SetDefault(l)

	return &Logger{l}, context.WithValue(ctx, logMapCtxKey, &sync.Map{})
}

func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		logger: l.logger.With(args...),
	}
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	var newLogger *Logger = l.With()
	if v, ok := ctx.Value(logMapCtxKey).(*sync.Map); ok {
		v.Range(func(key, value any) bool {
			if key, ok := key.(string); ok {
				newLogger = newLogger.With(key, ctx.Value(key))
			}
			return true
		})
	}
	for _, key := range keys {
		if ctx.Value(key) != nil {
			newLogger = newLogger.With(key, ctx.Value(key))
		}
	}
	return newLogger
}
