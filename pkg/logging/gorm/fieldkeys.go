package gorm

import (
	"time"

	"gorm.io/gorm/logger"
)

const (
	SourceKey  = "sourcefile"
	RowsKey    = "rows"
	ElaspedKey = "duration"
	SqlKey     = "sql"
)

func DefaultConfig() *logger.Config {
	return &logger.Config{
		SlowThreshold:             2000 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: true,
	}
}
