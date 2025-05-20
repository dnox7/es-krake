package log

import (
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/dpe27/es-krake/config"
	"golang.org/x/net/context"
)

type loggerCtxKey struct{}

const messageKey = "message"

var (
	keys         []string
	logMapCtxKey = loggerCtxKey{}
)

// Initialize: initializes the logger with default handler
func Initialize(w io.Writer, cfg *config.Config, keyInput []string) {
	keys = append(keys, keyInput...)
	level := slog.LevelInfo
	if cfg.Log.Level == "DEBUG" {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
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

	slog.SetDefault(slog.New(&handler{
		Handler: slog.NewJSONHandler(w, opts),
	}))
}

// AddLogValToCtx adds a key-val pair to the context in sync.Map for thread safely
// this value automatically added to the log record with defaultHandler
func AddLogValToCtx(ctx context.Context, key string, val interface{}) context.Context {
	m, ok := ctx.Value(logMapCtxKey).(*sync.Map)
	if !ok {
		m = &sync.Map{}
	}
	m.Store(key, val)
	return context.WithValue(ctx, logMapCtxKey, m)
}

func Group(key string, args ...any) slog.Attr {
	return slog.Group(key, args...)
}
