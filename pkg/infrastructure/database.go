package infrastructure

import (
	"database/sql"
	"errors"
	"os"
	customsource "pech/es-krake/pkg/infrastructure/custom-source"
	"pech/es-krake/pkg/logging/hook"
	wraperror "pech/es-krake/pkg/shared/wrap-error"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/golang-migrate/migrate/source"
	"github.com/sirupsen/logrus"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

func GetGormConfig() *gorm.Config {
	return &gorm.Config{
		DisableAutomaticPing: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
}

func NewDatabase(baseLogger *logrus.Logger) (
	db *gorm.DB,
	master *sql.DB,
	slave *sql.DB,
	err error,
) {
	gormConf := GetGormConfig()
	gormConf.Logger = hook.DefaultGormLogger(baseLogger).LogMode(logger.Info)
	gormConf.SkipDefaultTransaction = true

	connMaster, err := sql.Open("postgres", os.Getenv("POSTGRES_MASTER_CONNECTION_STRING"))
	if err != nil {
		return nil, nil, nil, err
	}

	connSlave, err := sql.Open("postgres", os.Getenv("POSTGRES_SLAVE_CONNECTION_STRING"))
	if err != nil {
		return nil, nil, nil, err
	}

	postgresConnMaxLifeTime, _ := strconv.Atoi(os.Getenv("POSTGRES_CONN_MAX_LIFTIME"))
	postgresConnMaxIdleTime, _ := strconv.Atoi(os.Getenv("POSTGRES_CONN_MAX_IDLE_TIME"))
	postgresMaxIdleConns, _ := strconv.Atoi(os.Getenv("POSTGRES_MAX_IDLE_CONNS"))
	postgresMaxOpenConns, _ := strconv.Atoi(os.Getenv("POSTGRES_MAX_OPEN_CONNS"))

	connMaster.SetConnMaxLifetime(time.Minute * time.Duration(postgresConnMaxLifeTime))
	connMaster.SetConnMaxIdleTime(time.Minute * time.Duration(postgresConnMaxIdleTime))
	connMaster.SetMaxIdleConns(postgresMaxIdleConns)
	connMaster.SetMaxOpenConns(postgresMaxOpenConns)

	connSlave.SetConnMaxLifetime(time.Minute * time.Duration(postgresConnMaxLifeTime))
	connSlave.SetConnMaxIdleTime(time.Minute * time.Duration(postgresConnMaxIdleTime))
	connSlave.SetMaxIdleConns(postgresMaxIdleConns)
	connSlave.SetMaxOpenConns(postgresMaxOpenConns)

	dirverMaster := postgresDriver.New(postgresDriver.Config{
		Conn: connMaster,
	})

	driverSlave := postgresDriver.New(postgresDriver.Config{
		Conn: connSlave,
	})

	db, err = gorm.Open(dirverMaster, gormConf)
	if err != nil {
		return nil, nil, nil, wraperror.WithTrace(err, nil, nil)
	}

	resolver := dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{driverSlave},
	})
	err = db.Use(resolver)
	if err != nil {
		return nil, nil, nil, err
	}

	err = Ping(db)
	if err != nil {
		return nil, nil, nil, err
	}

	return db, connMaster, connSlave, nil
}

func Ping(db *gorm.DB) error {
	postgresDB, err := db.DB()
	if err != nil {
		return err
	}

	return postgresDB.Ping()
}

func CloseDB(l *logrus.Logger, master *sql.DB, slave *sql.DB) {
	l.Info("Closing the master DB connection pool")
	if err := master.Close(); err != nil {
		l.Errorf("Error while closing the master DB connection pool: %v", err)
	}

	l.Info("Closing the slave DB connection pool")
	if err := slave.Close(); err != nil {
		l.Errorf("Error while closing the slave DB connection pool: %v", err)
	}
}

func StartLoggingPoolSize(pool *sql.DB, l *logrus.Entry) func() {
	stop := make(chan bool)
	go func() {
		previousOpened := 0
		for {
			time.Sleep(time.Second)
			select {
			case <-stop:
				logPoolSize(pool.Stats(), l)
			default:
				current := pool.Stats()
				if previousOpened != current.OpenConnections {
					logPoolSize(current, l)
					previousOpened = current.OpenConnections
				}
			}
		}
	}()

	return func() {
		stop <- true
	}
}

func logPoolSize(stats sql.DBStats, l *logrus.Entry) {
	l.WithFields(logrus.Fields{
		"inUse":  stats.InUse,
		"idle":   stats.Idle,
		"opened": stats.OpenConnections,
	}).Infof("Current number of opened connections in the pool")
}

// migrationLogger: this type is required to implement
// the Logger interface of golang-migrate
type migrationLogger struct {
	*logrus.Logger
}

func (l migrationLogger) Verbose() bool {
	return true
}

func (l migrationLogger) Printf(fmt string, v ...interface{}) {
	l.WithField("service", "database").Infof(fmt, v...)
}

func getMigrationPath(module string) string {
	return "file://" + os.Getenv("PE_MIGRATIONS_PATH") + "/" + module
}

func Migrate(masterDB *sql.DB, l *logrus.Logger, module string, migrationsTable string) error {
	postgres.DefaultMigrationsTable = migrationsTable
	driver, err := postgres.WithInstance(masterDB, &postgres.Config{})
	if err != nil {
		return err
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(getMigrationPath(module), "postgres", driver)
	if err != nil {
		return err
	}

	m.Log = migrationLogger{l}
	f := customsource.File{}
	fileSysMigrations, err := f.Open(os.Getenv("PE_MIGRATIONS_PATH") + "/" + module)
	if err != nil {
		return err
	}

	ver, _, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return err
	}

	lastestIndex := fileSysMigrations.GetLastestIndex()
	if ver > lastestIndex {
		err = m.Force(int(lastestIndex))
		if err != nil {
			return err
		}
	} else {
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			return err
		}
	}

	return nil
}

func CheckDatabaseVersion(masterDB *sql.DB, l *logrus.Logger, module string, migrationsTable string) error {
	postgres.DefaultMigrationsTable = migrationsTable
	driver, err := postgres.WithInstance(masterDB, &postgres.Config{})
	if err != nil {
		return err
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(getMigrationPath(module), "postgres", driver)
	if err != nil {
		return err
	}

	m.Log = migrationLogger{l}
	ver, dirty, err := m.Version()
	if err != nil {
		return err
	}

	if dirty {
		return migrate.ErrDirty{
			Version: int(ver),
		}
	}

	fileSysMigrations, err := source.Open(getMigrationPath(module))
	if err != nil {
		return err
	}

	_, err = fileSysMigrations.Next(ver)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}
