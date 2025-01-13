package main

import (
	"context"
	"net/http"
	"os"
	"pech/es-krake/pkg/infrastructure"
	"pech/es-krake/pkg/log"
)

func main() {
	ctx := context.Background()
	log.Initialize(ctx, os.Stdout, false, []string{"request-id", "recurringID"})

	db, master, slave, err := infrastructure.NewDatabase()
	if err != nil {
		log.Fatal(ctx, "Failed to connect to PostgreSQL", err)
		return
	}
	defer infrastructure.CloseDB(master, slave)

	stopMasterLogger := infrastructure.StartLoggingPoolSize(master, "master")
	defer stopMasterLogger()

	stopSlaveLogger := infrastructure.StartLoggingPoolSize(slave, "slave")
	defer stopSlaveLogger()

	router := infrastructure.NewServer(db)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	infrastructure.OnShutdown(func() error {
		return server.Shutdown(context.Background())
	})

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(ctx, "An error happened while starting the HTTP server: ", err)
	}
}
