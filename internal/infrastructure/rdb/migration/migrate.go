package migration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var migrationTables = map[string]string{
	"product": "schema_migrations_product",
}

type (
	direction string
	uintSlice []uint
)

const (
	Down direction = "down"
	Up   direction = "up"
)

func getSortedMigrationTableKeys() []string {
	i := 0
	keys := make([]string, len(migrationTables))
	for k := range migrationTables {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func MigrateAll(
	cfg *config.Config,
	db *sql.DB,
) error {
	keys := getSortedMigrationTableKeys()
	for _, module := range keys {
		if err := migrateSingleModule(cfg, db, module, migrationTables[module]); err != nil {
			return fmt.Errorf("Failed to migrate module %v: %w", module, err)
		}
	}
	return nil
}

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

func migrateSingleModule(
	cfg *config.Config,
	db *sql.DB,
	moduleName string,
	migrationTable string,
) error {
	conn, err := db.Conn(context.Background())
	if err != nil {
		return err
	}

	postgres.DefaultMigrationsTable = migrationTable
	driver, err := postgres.WithConnection(context.Background(), conn, &postgres.Config{})
	if err != nil {
		return err
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(getMigrationsPath(cfg, moduleName), cfg.RDB.Name, driver)
	if err != nil {
		return err
	}

	m.Log = migrationLogger{log.With()}
	f := file{}
	fileSystemMigrations, err := f.Open(cfg.RDB.MigrationsPath + "/" + moduleName)
	if err != nil {
		return err
	}

	ver, _, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return err
	}

	latestIndex := fileSystemMigrations.GetLastestIndex()
	if ver > latestIndex {
		err = m.Force(int(latestIndex))
		if err != nil {
			return err
		}
	} else {
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			return err
		}
	}

	return conn.Close()
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

	m.Log = migrationLogger{log.With()}
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

func getMigrationsPath(cfg *config.Config, moduleName string) string {
	return "file://" + cfg.RDB.MigrationsPath + "/" + moduleName
}

type migration struct {
	version    uint
	identifier string
	direction  direction
	raw        string
}

type migrations struct {
	index      uintSlice
	migrations map[uint]map[direction]*migration
}

func newMigrations() *migrations {
	return &migrations{
		index:      make(uintSlice, 0),
		migrations: make(map[uint]map[direction]*migration),
	}
}

func (m *migrations) Append(ele *migration) (ok bool) {
	if ele == nil {
		return false
	}

	if m.migrations[ele.version] == nil {
		m.migrations[ele.version] = make(map[direction]*migration)
	}

	if _, dup := m.migrations[ele.version][ele.direction]; dup {
		return false
	}

	m.migrations[ele.version][ele.direction] = ele
	m.buildIndex()

	return true
}

func (m *migrations) buildIndex() {
	m.index = make(uintSlice, 0, len(m.migrations))
	for ver := range m.migrations {
		m.index = append(m.index, ver)
	}
	sort.Slice(m.index, func(i, j int) bool {
		return m.index[i] < m.index[j]
	})
}
