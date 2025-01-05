package main

import (
	"context"
	"net/http"
	"pech/es-krake/pkg/infrastructure"
	"pech/es-krake/pkg/logging"
	"pech/es-krake/pkg/logging/hook"
)

func main() {
	logger := infrastructure.NewLogger()

	ctxKey := logging.NewCtxKeys("request-id", "recurringID")
	hook := hook.NewHookWithFallback(ctxKey)

	logger.AddHook(hook)

	db, master, slave, err := infrastructure.NewDatabase(logger)
	if err != nil {
		logger.Fatalln("Failed to connect to PostgreSQL", err)
		return
	}
	defer infrastructure.CloseDB(logger, master, slave)

	dbLogEntry := logger.WithField("service", "database")
	stopMasterLogger := infrastructure.StartLoggingPoolSize(master, dbLogEntry.WithField("pool", "master"))
	defer stopMasterLogger()

	stopSlaveLogger := infrastructure.StartLoggingPoolSize(slave, dbLogEntry.WithField("pool", "slave"))
	defer stopSlaveLogger()

	router := infrastructure.NewServer(db, logger)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	infrastructure.OnShutdown(logger, func() error {
		return server.Shutdown(context.Background())
	})

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Fatalln("An error happened while starting the HTTP server: ", err)
	}
}
