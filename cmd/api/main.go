package main

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb"
	"github.com/dpe27/es-krake/internal/infrastructure/rdb/migration"
	"github.com/dpe27/es-krake/internal/infrastructure/repository"
	vaultcli "github.com/dpe27/es-krake/internal/infrastructure/vault"
	"github.com/dpe27/es-krake/internal/initializer"
	"github.com/dpe27/es-krake/pkg/log"
)

func main() {
	cfg := config.NewConfig()

	ctx := context.Background()
	log.Initialize(os.Stdout, cfg, []string{"request-id", "recurringID"})

	vault, authToken, err := vaultcli.NewVaultAppRoleClient(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "unable to initialize vault connection", "address", cfg.Vault.Address, "error", err.Error())
	}

	rdbCred, rdbCredLease, err := vault.GetRdbCredentials(ctx)
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

	err = migration.CheckAll(cfg, pg.Conn())
	if err != nil {
		slog.Error("The database is not up-to-date: %v", "detail", err)
		return
	}

	var wg sync.WaitGroup
	renewLeaseCtx, stopRenew := context.WithCancel(ctx)
	wg.Add(1)
	go func() {
		vault.PeriodicallyRenewLeases(
			renewLeaseCtx,
			authToken,
			rdbCredLease,
			pg.RetryConn,
		)
		wg.Done()
	}()
	defer func() {
		stopRenew()
		wg.Wait()
	}()

	repositories := repository.NewRepositoriesContainer(pg)
	err = initializer.MountAll(repositories, pg)
	if err != nil {
		panic(err)
	}
}
