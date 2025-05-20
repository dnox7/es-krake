package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb/migration"
	"github.com/dpe27/es-krake/internal/infrastructure/repository"
	"github.com/dpe27/es-krake/internal/initializer"
	"github.com/dpe27/es-krake/pkg/log"
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
	if err != nil {
		panic(err)
	}
}
