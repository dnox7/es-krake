package main

import (
	"context"
	"flag"
	"os"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb/migration"
	"github.com/dpe27/es-krake/internal/initializer"
	"github.com/dpe27/es-krake/pkg/log"
)

func main() {
	cfg := config.NewConfig()

	ctx := context.Background()
	log.Initialize(os.Stdout, cfg, []string{})

	vault, _, err := initializer.InitVault(ctx, cfg)
	if err != nil {
		log.Error(ctx, "failed to init Vault", "error", err.Error())
		return
	}

	rdbCred, _, err := vault.GetRdbCredentials(ctx)
	if err != nil {
		log.Fatal(ctx, "unable to retrieve database credentials from vault", "error", err.Error())
	}

	pg := rdb.NewOrGetSingleton(cfg, rdbCred)
	defer pg.Close()

	loggingPoolSizeCtx, stopLogging := context.WithCancel(ctx)
	pg.LoggingPoolSize(loggingPoolSizeCtx)
	defer stopLogging()

	if err := pg.Ping(ctx); err != nil {
		log.Error(ctx, "database ping failed", "detail", err)
		return
	}

	migrateType := flag.String("type", "up", "Migration type: up, down, step (required)")
	step := flag.Int("step", 0, "Number of steps for 'step' action")
	module := flag.String("module", "", "Name of the module to migrate (required)")

	flag.Parse()
	switch *migrateType {
	case "step":
		err = migration.MigrateStep(cfg, pg.Conn(), *module, *step)
	case "down":
		err = migration.MigrateDown(cfg, pg.Conn())
	default:
		err = migration.MigrateUp(cfg, pg.Conn())
	}

	if err != nil {
		log.Fatal(ctx, "cannot migrate database", "detail", err)
	}
}
