package main

import (
	"context"
	"net/http"
	"os"

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
	if err != nil {
		log.Error(ctx, "failed to init Postgres", "error", err.Error())
		return
	}
	defer pg.Close()
	defer stopLoggingPoolSize()

	redisRepo, err := initializer.InitRedis(vault, cfg)
	if err != nil {
		log.Error(ctx, "failed to init Redis", "error", err.Error())
		return
	}
	defer redisRepo.Close(ctx)

	mongo, mongoCredLease, err := initializer.InitMongo(vault, cfg)
	if err != nil {
		log.Error(ctx, "failed to init MongoDB", "error", err)
		return
	}
	defer mongo.Close(ctx)

	esRepo, esCredLease, err := initializer.InitElasticSearch(vault, cfg)
	if err != nil {
		log.Error(ctx, "failed to init elasticsearch", "error", err)
		return
	}

	_, err = initializer.InitS3Repository(cfg)
	if err != nil {
		log.Error(ctx, err.Error())
		return
	}

	ses, err := initializer.InitSes(cfg)
	if err != nil {
		log.Error(ctx, "failed to init ses", "error", err)
		return
	}

	renewLeaseCtx, stopRenew := context.WithCancel(ctx)
	go func() {
		vault.PeriodicallyRenewLeases(
			renewLeaseCtx, authToken,
			rdbCredLease, pg.Reconn,
			mongoCredLease, mongo.Reconn,
			esCredLease, esRepo.Reconn,
		)
	}()
	defer stopRenew()

	router := infrastructure.NewGinRouter(cfg)
	server := &http.Server{
		Addr:    cfg.App.Port,
		Handler: router,
	}

	err = initializer.MountAPI(cfg, pg, mongo, redisRepo, ses, router)
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
