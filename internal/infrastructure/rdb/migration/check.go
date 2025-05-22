package migration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/dpe27/es-krake/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
)

func CheckAll(
	cfg *config.Config,
	db *sql.DB,
) error {
	for moduleName, mt := range migrationTables {
		if err := checkDatabaseVersion(cfg, db, moduleName, mt); err != nil {
			return fmt.Errorf("The migrations for module %v are not up-to-date: %w", moduleName, err)
		}
	}
	return nil
}

func checkDatabaseVersion(
	cfg *config.Config,
	db *sql.DB,
	moduleName string,
	migrationsTable string,
) error {
	conn, err := db.Conn(context.Background())
	if err != nil {
		return err
	}

	postgres.DefaultMigrationsTable = migrationsTable
	driver, err := postgres.WithConnection(context.Background(), conn, &postgres.Config{})
	if err != nil {
		return err
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(getMigrationsPath(cfg, moduleName), cfg.RDB.Name, driver)
	if err != nil {
		return err
	}

	m.Log = newMigrationLogger()
	ver, dirty, err := m.Version()
	if err != nil {
		return err
	}
	if dirty {
		return migrate.ErrDirty{
			Version: int(ver),
		}
	}

	fileSystemMigrations, err := source.Open(getMigrationsPath(cfg, moduleName))
	if err != nil {
		return err
	}

	_, err = fileSystemMigrations.Next(ver)
	if errors.Is(err, os.ErrNotExist) {
		// No higher version is available = the DB is up to date
		return nil
	}

	return conn.Close()
}
