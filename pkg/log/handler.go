package log

import (
	"log/slog"
	"sync"

	"golang.org/x/net/context"
)

type customHandler struct {
	slog.Handler
}

func (t *customHandler) Handle(ctx context.Context, r slog.Record) error {
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

func (t *customHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &customHandler{
		Handler: t.Handler.WithAttrs(attrs),
	}
}
