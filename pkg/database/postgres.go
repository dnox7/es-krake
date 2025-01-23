package database

import (
	"database/sql"
	"time"
)

type PostgreSQL struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	DB *sql.DB
}
