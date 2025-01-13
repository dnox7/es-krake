package log

import (
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"golang.org/x/net/context"
)

type Logger struct {
	logger *slog.Logger
}

type defaultHandler struct {
	slog.Handler
}

var (
	keys []string
	_    slog.Handler = &defaultHandler{}
)

func (t *defaultHandler) Handle(ctx context.Context, r slog.Record) error {
	if v, ok := ctx.Value(logMapCtxKey).(*sync.Map); ok {
		v.Range(func(key, value any) bool {
			if key, ok := key.(string); ok {
				r.AddAttrs(slog.Any(key, value))
			}
			return true
		})
	}
	for _, key := range keys {
		if ctx.Value(key) != nil {
			r.AddAttrs(slog.Any(key, ctx.Value(key)))
		}
	}
	return t.Handler.Handle(ctx, r)
}

func (t *defaultHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &defaultHandler{
		Handler: t.Handler.WithAttrs(attrs),
	}
}

type loggerCtxKey struct{}

var logMapCtxKey loggerCtxKey = loggerCtxKey{}

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

const messageKey = "message"

// Initialize: initializes the logger with default handler
// if is deDebug is false, it sets the log level to [LevelDebug] otherwise [LevelInfo]
// it uses JSONHandler for logging and automatically adds the key-value pair to the log record
// using the CtxWithValue function.
// The fields in the keys are used to retrieve their values from the context and write them to the logger
func Initialize(ctx context.Context, w io.Writer, isDebug bool, keyInput []string) context.Context {
	keys = append(keys, keyInput...)
	level := slog.LevelInfo
	if isDebug {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(&defaultHandler{
		Handler: slog.NewJSONHandler(w, &slog.HandlerOptions{
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
		}),
	}))
	return context.WithValue(ctx, logMapCtxKey, &sync.Map{})
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
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

func Info(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, args...)
}

func Debug(ctx context.Context, msg string, args ...any) {
	slog.DebugContext(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}

func Fatal(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
	os.Exit(1)
}

func With(args ...any) *Logger {
	return &Logger{
		logger: slog.With(args...),
	}
}

func Group(key string, args ...any) slog.Attr {
	return slog.Group(key, args...)
}
