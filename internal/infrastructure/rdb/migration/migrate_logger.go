package migration

import (
	"context"
	"fmt"

	"github.com/dpe27/es-krake/pkg/log"
	"github.com/golang-migrate/migrate/v4"
)

// this type is required to implement the Logger interface of golang-migrate
type migrationLogger struct {
	logger *log.Logger
}

func newMigrationLogger() migrate.Logger {
	return migrationLogger{
		logger: log.With("object", "migration"),
	}
}

func (l migrationLogger) Verbose() bool {
	return true
}

func (l migrationLogger) Printf(format string, v ...interface{}) {
	l.logger.Info(context.Background(), fmt.Sprintf(format, v...))
}
