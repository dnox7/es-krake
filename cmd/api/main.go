package main

import (
	"context"
	"log/slog"
	"os"
	"pech/es-krake/config"
	"pech/es-krake/internal/infrastructure/rdb"
	"pech/es-krake/internal/infrastructure/rdb/migration"
	"pech/es-krake/internal/infrastructure/repository"
	"pech/es-krake/internal/initializer"
	"pech/es-krake/pkg/log"
)

func main() {
	cfg := config.NewConfig()

	ctx := context.Background()
	log.Initialize(os.Stdout, cfg, []string{"request-id", "recurringID"})

	pg := rdb.NewOrGetSingleton(cfg)
	defer pg.Close()

	poolLogger := pg.StartLoggingPoolSize()
	defer poolLogger()

	if err := pg.Ping(ctx); err != nil {
		slog.Error("database ping failed", "detail", err)
		return
	}

	err := migration.CheckAll(cfg, pg.Conn())
	if err != nil {
		slog.Error("The database is not up-to-date: %v", "detail", err)
		return
	}

	repositories := repository.NewRepositoriesContainer(pg)
	err = initializer.MountAll(repositories, pg)
}
