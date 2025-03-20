package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"pech/es-krake/config"
	"pech/es-krake/internal/infrastructure/db"
	"pech/es-krake/internal/infrastructure/db/migration"
	"pech/es-krake/internal/infrastructure/repository"
	"pech/es-krake/internal/initializer"
	"pech/es-krake/pkg/log"
)

func main() {
	cfg := config.NewConfig()

	ctx := context.Background()
	log.Initialize(context.Background(), os.Stdout, cfg, []string{"request-id", "recurringID"})

	pg := db.NewOrGetSingleton(cfg)
	defer pg.Close()

	if err := pg.PingContext(ctx); err != nil {
		slog.Error("database ping failed", "detail", err)
		return
	}

	migrateFlag := flag.Bool("migrate", false, "Updates the database up to the latest migration")
	flag.Parse()
	if *migrateFlag {
		err := migration.MigrateAll(cfg, pg.DB.DB)
		if err != nil {
			slog.Error("cannot migrate database", "detail", err)
		}
		return
	}

	err := migration.CheckAll(cfg, pg.DB.DB)
	if err != nil {
		slog.Error("The database is not up-to-date: %v", "detail", err)
		return
	}

	repositories := repository.NewRepositoriesContainer(pg)
	err = initializer.MountAll(repositories, pg)

}
