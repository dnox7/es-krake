package hook

import (
	"pech/es-krake/pkg/logging"

	"github.com/sirupsen/logrus"
)

type ContextValuesAddingHook struct {
	fallback *logging.ContextFieldKeys
}

func NewHook() *ContextValuesAddingHook {
	return NewHookWithFallback(nil)
}

// NewHookWithFallback: returns a new logrus hook with fallback.
// fallback will be used when ContextFieldKeys is not stored in Context.
// To store ContextFieldKeys in Context, use ContextFieldKeys.SaveTo(Context)
func NewHookWithFallback(fallback *logging.ContextFieldKeys) *ContextValuesAddingHook {
	return &ContextValuesAddingHook{
		fallback: fallback,
	}
}

func (h *ContextValuesAddingHook) Fire(entry *logrus.Entry) (err error) {
	ctx := entry.Context
	if ctx == nil {
		return
	}

	ctxKeys := logging.GetContextFieldKeys(ctx, h.fallback)
	if ctxKeys == nil {
		return
	}

	fields := ctxKeys.GetValues(ctx)
	for k, v := range fields {
		entry.Data[k] = v
	}
	return
}

func (h *ContextValuesAddingHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
