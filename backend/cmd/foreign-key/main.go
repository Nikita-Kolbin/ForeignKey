package main

import (
	"ForeignKey/internal/image"
	"ForeignKey/internal/notification_client"
	"ForeignKey/internal/router"
	"log/slog"
	"net/http"
	"os"

	_ "ForeignKey/docs"
	"ForeignKey/internal/config"
	"ForeignKey/internal/logger"
	"ForeignKey/internal/storage/sqlite"

	httpSwagger "github.com/swaggo/http-swagger"
	l "log"
)

const dotenvPath = "./backend/config/local.env"

// @title           ForeignKey
// @version         1.0

// @host      localhost:8082
// @BasePath  /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// config
	cfg := config.MustLoad(dotenvPath)

	// logger
	log, err := logger.SetupLogger(cfg.Env)
	if err != nil {
		l.Fatalf("can't setup logger: %s", err)
	}
	log.Info("starting service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// storage
	storage, err := sqlite.New(cfg.StoragePath, cfg.StorageName)
	if err != nil {
		log.Error("failed to initialize storage", slog.String("error", err.Error()))
		os.Exit(0)
	}

	// image saver
	imageSaver, err := image.New(cfg.ImagesPath)
	if err != nil {
		log.Error("failed to initialize image saver", slog.String("error", err.Error()))
		os.Exit(0)
	}

	// email
	nc := notification_client.New(cfg.NotificationClientConfig)

	// http router
	r := router.New(storage, imageSaver, nc, log)

	// swagger
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(cfg.Address+"/swagger/doc.json"),
	))

	// server
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Info("starting server", slog.String("addr", cfg.Address))
	if err = srv.ListenAndServe(); err != nil {
		log.Error("failed to start server:", slog.String("err", err.Error()))
	}

	log.Error("server stopped")
}
