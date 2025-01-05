package migrate

import (
	"database/sql"
	"fmt"
	"pech/es-krake/pkg/infrastructure"
	"sort"

	"github.com/sirupsen/logrus"
)

var migrationsTables = map[string]string{}

func getAllMigrationTableKeys() []string {
	j := 0
	keys := make([]string, len(migrationsTables))
	for k := range migrationsTables {
		keys[j] = k
		j++
	}
	sort.Strings(keys)
	return keys
}

func MigrateAll(masterDB *sql.DB, logger *logrus.Logger) error {
	keys := getAllMigrationTableKeys()
	for _, module := range keys {
		err := infrastructure.Migrate(masterDB, logger, module, migrationsTables[module])
		if err != nil {
			return fmt.Errorf("Failed to migrate module %v: %w", module, err)
		}
	}
	return nil
}

func CheckAll(masterDB *sql.DB, logger *logrus.Logger) error {
	for module, migrateTable := range migrationsTables {
		err := infrastructure.CheckDatabaseVersion(masterDB, logger, module, migrateTable)
		if err != nil {
			return fmt.Errorf("The migrations for module %v are not up-to-date: %w", module, err)
		}
	}
	return nil
}
