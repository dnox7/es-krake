package infrastructure

import (
	"context"
	"os"
	"os/signal"
	"pech/es-krake/pkg/log"
)

func OnShutdown(callback func() error) {
	go (func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)
		s := <-signals
		log.Info(context.Background(), "Received signal '%v'. Shutting down...", s.String())

		if err := callback(); err != nil {
			log.Error(context.Background(), "Error during graceful shutdown: %v", err)
		}
	})()
}
