package main

import (
	"ForeignKey/internal/config"
	"ForeignKey/internal/logger"
	l "log"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()

	log, err := logger.SetupLogger(cfg.Env)
	if err != nil {
		l.Fatalf("can't setup logger: %s", err)
	}

	log.Info("starting service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// TODO: init storage: postgres

	// TODO: init router: chi, render

	// TODO: init server
}
