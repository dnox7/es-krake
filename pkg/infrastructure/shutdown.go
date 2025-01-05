package infrastructure

import (
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
)

func OnShutdown(logger *logrus.Logger, callback func() error) {
	go (func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)
		s := <-signals
		logger.Printf("Received signal '%v'. Shutting down...", s.String())

		if err := callback(); err != nil {
			logger.Fatalf("Error during graceful shutdown: %v", err)
		}
	})()
}
