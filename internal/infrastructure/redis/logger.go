package redis

import (
	"context"
	"fmt"

	"github.com/dpe27/es-krake/pkg/log"
)

type redisLogger struct {
	logger *log.Logger
}

func (l *redisLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	l.logger.Error(ctx, fmt.Sprintf(format, v...))
}
