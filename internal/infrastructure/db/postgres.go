package db

import (
	"context"
	"fmt"
	"log/slog"
	"pech/es-krake/config"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type PostgreSQL struct {
	DB      *sqlx.DB
	Builder squirrel.StatementBuilderType

	maxOpenConns int
	maxIdleConns int
	maxLifeTime  time.Duration
	maxIdleTime  time.Duration
	connTimeout  time.Duration
	connAttempts int
}

var (
	postgresInstance *PostgreSQL
	once             sync.Once
)

func NewOrGetSingleton(cfg *config.Config) *PostgreSQL {
	once.Do(func() {
		pg, err := initPostgres(cfg)
		if err != nil {
			panic(err)
		}
		postgresInstance = pg
	})
	return postgresInstance
}

func (pg *PostgreSQL) Close() error {
	return pg.DB.Close()
}

func (pg *PostgreSQL) Ping() error {
	return pg.DB.Ping()
}

func (pg *PostgreSQL) PingContext(ctx context.Context) error {
	return pg.DB.PingContext(ctx)
}

func initPostgres(cfg *config.Config) (*PostgreSQL, error) {
	pg := &PostgreSQL{
		maxOpenConns: cfg.RDB.MaxOpenConns,
		maxIdleConns: cfg.RDB.MaxIdleConns,
		maxLifeTime:  time.Duration(cfg.RDB.MaxIdleTime) * time.Millisecond,
		maxIdleTime:  time.Duration(cfg.RDB.MaxIdleTime) * time.Millisecond,
		connTimeout:  time.Duration(cfg.RDB.ConnTimeout) * time.Millisecond,
		connAttempts: cfg.RDB.ConnAttempts,
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	dsn, err := buildDSN(cfg)
	if err != nil {
		slog.Error("Error while building data source name", "detail", err.Error())
		return nil, err
	}

	for pg.connAttempts > 0 {
		xdb, err := sqlx.Open(cfg.RDB.Driver, dsn)
		if err == nil {
			xdb.DB.SetMaxOpenConns(pg.maxOpenConns)
			xdb.DB.SetMaxIdleConns(pg.maxIdleConns)
			xdb.DB.SetConnMaxLifetime(pg.maxLifeTime)
			xdb.DB.SetConnMaxIdleTime(pg.maxIdleTime)

			pg.DB = xdb
			break
		}

		slog.Warn("PostgreSQL is trying to connect", "attempts left", pg.connAttempts)
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	if pg.DB == nil {
		return nil, fmt.Errorf("PostgreSQL (initDB): connection failed")
	}

	return pg, nil
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
