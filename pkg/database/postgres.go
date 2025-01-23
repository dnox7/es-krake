package database

import (
	"fmt"
	"log/slog"
	"pech/es-krake/config"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type PostgreSQL struct {
	DB      *sqlx.DB
	Builder squirrel.StatementBuilderType
}

var (
	postgresSingleton *PostgreSQL
	once              sync.Once
)

func NewOrGetSingleton(cfg *config.Config) *PostgreSQL {
	once.Do(func() {
	})
	return nil
}

func initPostgres(cfg *config.Config) (*PostgreSQL, error) {
	pg := &PostgreSQL{}
	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	dsn, err := buildDSN(cfg)
	if err != nil {
		slog.Error("Error while building data source name", "detail", err.Error())
		return nil, err
	}

	connAttempts := cfg.RDB.ConnAttempts
	maxPoolSize := cfg.RDB.MaxOpenConns
	for connAttempts > 0 {
		xdb, err := sqlx.Open(cfg.RDB.Driver, dsn)
		if err == nil {
			xdb.DB.SetMaxOpenConns(maxPoolSize)
			xdb.DB.SetMaxIdleConns()
		}
	}
}

func buildDSN(cfg *config.Config) (string, error) {
	if cfg.RDB.Driver != "postgres" {
		return "", fmt.Errorf("Database driver is invalid")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.RDB.Host,
		cfg.RDB.Port,
		cfg.RDB.Username,
		cfg.RDB.Password,
		cfg.RDB.Name,
		cfg.RDB.SSLMode,
	)

	return dsn, nil
}
