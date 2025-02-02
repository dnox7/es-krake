package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"pech/es-krake/config"
	"pech/es-krake/internal/infrastructure/db"
	"pech/es-krake/internal/infrastructure/db/migration"
	"pech/es-krake/pkg/log"
)

func main() {
	cfg := config.NewConfig()

	logger, ctx := log.Initialize(context.Background(), os.Stdout, cfg, []string{"request-id", "recurringID"})

	pg := db.NewOrGetSingleton(cfg)
	defer pg.Close()

	if err := pg.PingContext(ctx); err != nil {
		logger.Error("database ping failed", "detail", err)
		return
	}

	migrateFlag := flag.Bool("migrate", false, "Updates the database up to the latest migrations")
	flag.Parse()
	if *migrateFlag {
		err := migration.MigrateAll(cfg, pg.DB.DB, logger)
		if err != nil {
			logger.Error("cannot migrate database", "detail", err)
		}
		return
	}

	err := migration.CheckAll(cfg, pg.DB.DB, logger)
	if err != nil {
		logger.Error("The database is not up-to-date: %v", err)
		return
	}

	fmt.Println("OK!")
}
