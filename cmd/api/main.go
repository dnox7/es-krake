package main

import (
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/dpe27/es-krake/config"
	"github.com/dpe27/es-krake/internal/infrastructure"
	"github.com/dpe27/es-krake/internal/initializer"
	"github.com/dpe27/es-krake/pkg/log"
)

func main() {
	cfg := config.NewConfig()

	ctx := context.Background()
	log.Initialize(os.Stdout, cfg, []string{"request-id", "recurringID"})

	vault, authToken, err := initializer.InitVault(ctx, cfg)
	if err != nil {
		log.Error(ctx, "failed to init Vault", "error", err.Error())
		return
	}

	pg, rdbCredLease, stopLoggingPoolSize, err := initializer.InitPostgres(vault, cfg)
	defer pg.Close()
	defer stopLoggingPoolSize()
	if err != nil {
		log.Error(ctx, "failed to init Postgres", "error", err.Error())
		return
	}

	redisRepo, err := initializer.InitRedis(vault, cfg)
	defer redisRepo.Close(ctx)
	if err != nil {
		log.Error(ctx, "failed to init Redis", "error", err.Error())
		return
	}

	var wg sync.WaitGroup
	renewLeaseCtx, stopRenew := context.WithCancel(ctx)
	wg.Add(1)
	go func() {
		vault.PeriodicallyRenewLeases(
			renewLeaseCtx, authToken,
			rdbCredLease, pg.RetryConn,
		)
		wg.Done()
	}()
	defer func() {
		stopRenew()
		wg.Wait()
	}()

	router := infrastructure.NewGinRouter(cfg)
	server := &http.Server{
		Addr:    cfg.App.Port,
		Handler: router,
	}

	err = initializer.MountAll(pg, router, cfg)
	if err != nil {
		log.Fatal(ctx, "failed to mount dependencies", "error", err.Error())
	}

	infrastructure.OnShutdown(func() error {
		return server.Shutdown(ctx)
	})

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(ctx, "an error happened while starting the HTTP server", "error", err.Error())
	}
}
