package rdb

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/pkg/log"
	gormlog "github.com/dpe27/es-krake/pkg/log/gorm"
	"github.com/dpe27/es-krake/pkg/wraperror"
	pgDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type PostgreSQL struct {
	DB     *gorm.DB
	conn   *sql.DB
	logger *log.Logger

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

func initPostgres(cfg *config.Config) (*PostgreSQL, error) {
	pg := &PostgreSQL{
		logger:       log.With("service", "postgres"),
		maxOpenConns: cfg.RDB.MaxOpenConns,
		maxIdleConns: cfg.RDB.MaxIdleConns,
		maxLifeTime:  time.Duration(cfg.RDB.MaxIdleTime) * time.Millisecond,
		maxIdleTime:  time.Duration(cfg.RDB.MaxIdleTime) * time.Millisecond,
		connTimeout:  time.Duration(cfg.RDB.ConnTimeout) * time.Millisecond,
		connAttempts: cfg.RDB.ConnAttempts,
	}

	dsn, err := pg.buildDSN(cfg)
	if err != nil {
		slog.Error("Error while building data source name", "detail", err.Error())
		return nil, err
	}

	gormCfg := pg.getGormConfig()

	for pg.connAttempts > 0 {
		err = pg.retryConn(dsn, cfg.RDB.Driver, gormCfg)
		if err == nil {
			break
		}

		pg.logger.Warn(
			context.Background(),
			"PostgreSQL is trying to connect",
			"error", err.Error(),
			"attempts left", pg.connAttempts,
		)
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	if pg.DB == nil {
		return nil, fmt.Errorf("PostgreSQL (initDB): connection failed")
	}

	return pg, nil
}

func (pg *PostgreSQL) Ping(ctx context.Context) error {
	pgDB, err := pg.DB.DB()
	if err != nil {
		return nil
	}
	return pgDB.PingContext(ctx)
}

func (pg *PostgreSQL) Close() {
	pg.logger.Info(context.Background(), "Closing the DB connaction pool")
	if err := pg.conn.Close(); err != nil {
		pg.logger.Error(context.Background(), "Error while closing the DB connactio pool")
	}
}

func (pg *PostgreSQL) Conn() *sql.DB {
	return pg.conn
}

func (pg *PostgreSQL) StartLoggingPoolSize() func() {
	stop := make(chan bool)
	go func() {
		previousOpened := 0
		for {
			time.Sleep(time.Second)
			select {
			case <-stop:
				pg.logPoolSize(pg.conn.Stats())
				return
			default:
				curr := pg.conn.Stats()
				if previousOpened != curr.OpenConnections {
					pg.logPoolSize(curr)
					previousOpened = curr.OpenConnections
				}
			}
		}
	}()

	return func() {
		stop <- true
	}
}

func (pg *PostgreSQL) logPoolSize(stats sql.DBStats) {
	pg.logger.With("inUse", stats.InUse).
		With("idle", stats.Idle).
		With("opened", stats.OpenConnections).
		Info(context.Background(), "Current  number of opened connections in the pool")
}

func (pg *PostgreSQL) retryConn(dsn string, driverName string, gormCfg *gorm.Config) error {
	conn, err := sql.Open(driverName, dsn)
	if err != nil {
		slog.Error("failed to open sql conn", "detail", err.Error())
		return err
	}

	conn.SetMaxOpenConns(pg.maxOpenConns)
	conn.SetMaxIdleConns(pg.maxIdleConns)
	conn.SetConnMaxLifetime(pg.maxLifeTime)
	conn.SetConnMaxIdleTime(pg.maxIdleTime)

	driver := pgDriver.New(pgDriver.Config{
		WithoutQuotingCheck: true,
		Conn:                conn,
	})

	db, err := gorm.Open(driver, gormCfg)
	if err != nil {
		err = wraperror.WithTrace(err, nil, nil)
		slog.Error(err.Error())
	}

	pg.conn = conn
	pg.DB = db
	return nil
}

func (pg *PostgreSQL) buildDSN(cfg *config.Config) (string, error) {
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

func (pg *PostgreSQL) getGormConfig() *gorm.Config {
	gormLogger := gormlog.DefaultGormLogger().LogMode(logger.Info)
	return &gorm.Config{
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true,
		Logger:                 gormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
}
