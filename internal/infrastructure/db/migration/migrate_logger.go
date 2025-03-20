package migration

import (
	"context"
	"fmt"
	"pech/es-krake/pkg/log"
)

// this type is required to implement the Logger interface of golang-migrate
type migrationLogger struct {
	*log.Logger
}

func (l migrationLogger) Verbose() bool {
	return true
}

func (l migrationLogger) Printf(format string, v ...interface{}) {
	l.With("service", "database").Info(context.Background(), fmt.Sprintf(format, v...))
}
