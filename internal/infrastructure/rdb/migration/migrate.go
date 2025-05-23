package migration

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	"github.com/dpe27/es-krake/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type (
	migrationType int
	direction     string
	uintSlice     []uint
)

const (
	Down direction = "down"
	Up   direction = "up"

	migrateUpAll migrationType = iota + 1
	migrateDownAll
	migrateStep
)

var migrationTables = map[string]string{
	"product": "schema_migrations_product",
}

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

func getMigrationsPath(cfg *config.Config, moduleName string) string {
	return "file://" + cfg.RDB.MigrationsPath + "/" + moduleName
}

func MigrateUp(
	cfg *config.Config,
	db *sql.DB,
) error {
	keys := getSortedMigrationTableKeys()
	for _, module := range keys {
		if err := migrateModule(cfg, db, module, migrationTables[module], migrateUpAll, 0); err != nil {
			return fmt.Errorf("Failed to migrate module %v: %w", module, err)
		}
	}
	return nil
}

func MigrateDown(
	cfg *config.Config,
	db *sql.DB,
) error {
	keys := getSortedMigrationTableKeys()
	for _, module := range keys {
		if err := migrateModule(cfg, db, module, migrationTables[module], migrateDownAll, 0); err != nil {
			return fmt.Errorf("failed to migrate down module %v: %w", module, err)
		}
	}
	return nil
}

func MigrateStep(
	cfg *config.Config,
	db *sql.DB,
	module string,
	step int,
) error {
	table, ok := migrationTables[module]
	if !ok {
		return fmt.Errorf("migration table of module %s does not exist", module)
	}
	return migrateModule(cfg, db, module, table, migrateStep, step)
}

func migrateModule(
	cfg *config.Config,
	db *sql.DB,
	moduleName string,
	migrationTable string,
	migType migrationType,
	step int,
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

	m.Log = newMigrationLogger()
	f := file{}
	fileSystemMigrations, err := f.Open(cfg.RDB.MigrationsPath + "/" + moduleName)
	if err != nil {
		return err
	}

	ver, _, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return err
	}

	switch migType {
	case migrateUpAll:
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
	case migrateDownAll:
		err = m.Down()
		if err != nil && err != migrate.ErrNoChange {
			return err
		}
	case migrateStep:
		err = m.Steps(step)
		if err != nil && err != migrate.ErrNoChange {
			return err
		}
	}

	return conn.Close()
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
