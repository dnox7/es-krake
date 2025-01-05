package logging

import (
	"context"

	"github.com/sirupsen/logrus"
)

const ctxKeyKey = "__ctxKey_key"

type ContextFieldKeys struct {
	keys *[]string
}

func NewCtxKeys(keys ...string) *ContextFieldKeys {
	return &ContextFieldKeys{
		keys: &keys,
	}
}

func (c *ContextFieldKeys) WithKeys(keys ...string) *ContextFieldKeys {
	if c == nil {
		return NewCtxKeys(keys...)
	}

	s := make([]string, 0, len(*c.keys)+len(keys))
	for _, v := range *c.keys {
		s = append(s, v)
	}

	for _, v := range keys {
		s = append(s, v)
	}

	return &ContextFieldKeys{
		keys: &s,
	}
}

func (c *ContextFieldKeys) WithContextKeys(keys *ContextFieldKeys) *ContextFieldKeys {
	if c == nil {
		return keys
	}
	return c.WithKeys(*keys.keys...)
}

func (c *ContextFieldKeys) Dup() *ContextFieldKeys {
	if c == nil {
		return nil
	}
	return &*c
}

func (c *ContextFieldKeys) GetValues(ctx context.Context) logrus.Fields {
	if c == nil || ctx == nil {
		return logrus.Fields{}
	}

	fields := make(logrus.Fields, len(*c.keys))
	for _, v := range *c.keys {
		cv := ctx.Value(v)
		if cv == nil {
			continue
		}

		fields[v] = ctx.Value(v)
	}
	return fields
}

func (c *ContextFieldKeys) SaveTo(ctx context.Context) context.Context {
	if c == nil || ctx == nil {
		return ctx
	}
	return context.WithValue(ctx, ctxKeyKey, c)
}

func GetContextFieldKeys(ctx context.Context, defaultValue *ContextFieldKeys) *ContextFieldKeys {
	if v, ok := ctx.Value(ctxKeyKey).(*ContextFieldKeys); ok {
		return v
	}
	return defaultValue
}
