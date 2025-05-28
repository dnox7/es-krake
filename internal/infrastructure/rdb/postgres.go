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
	mu   sync.RWMutex
	db   *gorm.DB
	conn *sql.DB

	logger  *log.Logger
	gormCfg *gorm.Config
	params  pgParams
}

type pgParams struct {
	driver       string
	host         string
	port         string
	name         string
	sslmode      string
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

func NewOrGetSingleton(cfg *config.Config, cred *config.RdbCredentials) *PostgreSQL {
	once.Do(func() {
		pg, err := initPostgres(cfg, cred)
		if err != nil {
			panic(err)
		}
		postgresInstance = pg
	})
	return postgresInstance
}

func initPostgres(cfg *config.Config, cred *config.RdbCredentials) (*PostgreSQL, error) {
	pg := &PostgreSQL{
		logger: log.With("service", "postgres"),
		params: pgParams{
			driver:       cfg.RDB.Driver,
			host:         cfg.RDB.Host,
			port:         cfg.RDB.Port,
			name:         cfg.RDB.Name,
			sslmode:      cfg.RDB.SSLMode,
			maxOpenConns: cfg.RDB.MaxOpenConns,
			maxIdleConns: cfg.RDB.MaxIdleConns,
			maxLifeTime:  time.Duration(cfg.RDB.MaxIdleTime) * time.Millisecond,
			maxIdleTime:  time.Duration(cfg.RDB.MaxIdleTime) * time.Millisecond,
			connTimeout:  time.Duration(cfg.RDB.ConnTimeout) * time.Millisecond,
			connAttempts: cfg.RDB.ConnAttempts,
		},
	}

	pg.setGormConfig()
	pg.RetryConn(cred)

	if pg.db == nil {
		return nil, fmt.Errorf("PostgreSQL (initDB): connection failed")
	}

	return pg, nil
}

func (pg *PostgreSQL) GormDB() *gorm.DB {
	pg.mu.RLock()
	defer pg.mu.RUnlock()
	return pg.db
}

func (pg *PostgreSQL) Conn() *sql.DB {
	pg.mu.RLock()
	defer pg.mu.RUnlock()
	return pg.conn
}

func (pg *PostgreSQL) updateDB(newDB *gorm.DB, newConn *sql.DB) {
	pg.mu.Lock()
	defer pg.mu.Unlock()
	pg.db = newDB
	pg.conn = newConn
}

func (pg *PostgreSQL) RetryConn(cred *config.RdbCredentials) {
	connAttempts := pg.params.connAttempts
	for connAttempts > 0 {
		err := pg.connect(cred)
		pg.logger.Warn(
			context.Background(),
			"PostgreSQL is trying to connect",
			"error", err.Error(),
			"attempts left", pg.params.connAttempts,
		)
		time.Sleep(pg.params.connTimeout)
		connAttempts--
	}
}

func (pg *PostgreSQL) connect(
	cred *config.RdbCredentials,
) error {
	dsn, err := pg.buildDSN(cred)
	if err != nil {
		pg.logger.Error(context.Background(), "Error while building data source name", "error", err.Error())
		return err
	}

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		pg.logger.Error(context.Background(), "failed to open sql conn", "detail", err.Error())
		return err
	}

	conn.SetMaxOpenConns(pg.params.maxOpenConns)
	conn.SetMaxIdleConns(pg.params.maxIdleConns)
	conn.SetConnMaxLifetime(pg.params.maxLifeTime)
	conn.SetConnMaxIdleTime(pg.params.maxIdleTime)

	driver := pgDriver.New(pgDriver.Config{
		WithoutQuotingCheck: true,
		Conn:                conn,
	})

	db, err := gorm.Open(driver, pg.gormCfg)
	if err != nil {
		err = wraperror.WithTrace(err, nil, nil)
		slog.Error(err.Error())
	}

	pg.updateDB(db, conn)
	return nil
}

func (pg *PostgreSQL) Ping(ctx context.Context) error {
	pgDB, err := pg.GormDB().DB()
	if err != nil {
		return nil
	}
	return pgDB.PingContext(ctx)
}

func (pg *PostgreSQL) Close() {
	pg.logger.Info(context.Background(), "Closing the DB connaction pool")
	if err := pg.Conn().Close(); err != nil {
		pg.logger.Error(context.Background(), "Error while closing the DB connactio pool")
	}
}

func (pg *PostgreSQL) StartLoggingPoolSize() func() {
	stop := make(chan bool)
	go func() {
		previousOpened := 0
		for {
			time.Sleep(time.Second)
			select {
			case <-stop:
				pg.logPoolSize(pg.Conn().Stats())
				return
			default:
				curr := pg.Conn().Stats()
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

func (pg *PostgreSQL) buildDSN(cred *config.RdbCredentials) (string, error) {
	if pg.params.driver != "postgres" {
		return "", fmt.Errorf("Database driver is invalid")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pg.params.host,
		pg.params.port,
		cred.Username,
		cred.Password,
		pg.params.name,
		pg.params.sslmode,
	)

	return dsn, nil
}

func (pg *PostgreSQL) setGormConfig() {
	gormLogger := gormlog.DefaultGormLogger().LogMode(logger.Info)
	pg.gormCfg = &gorm.Config{
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true,
		Logger:                 gormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
}
