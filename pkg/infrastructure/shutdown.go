package infrastructure

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
)

func OnShutdown(callback func() error) {
	go (func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)
		s := <-signals
		slog.InfoContext(context.Background(), "Received signal '%v'. Shutting down...", "s", s.String())

		if err := callback(); err != nil {
			slog.ErrorContext(context.Background(), "Error during graceful shutdown", "detail", err)
		}
	})()
}
