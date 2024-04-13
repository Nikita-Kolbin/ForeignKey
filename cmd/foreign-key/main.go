package main

import (
	"ForeignKey/internal/config"
	"ForeignKey/internal/logger"
	"ForeignKey/internal/storage/sqlite"
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

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", slog.String("error", err.Error()))
	}

	_ = storage

	// TODO: init router: chi, render

	// TODO: init server
}
