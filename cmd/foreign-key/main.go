package main

import (
	"ForeignKey/internal/config"
	"ForeignKey/internal/http-server/handlers/hello"
	mwLogger "ForeignKey/internal/http-server/middleware/logger"
	"ForeignKey/internal/logger"
	"ForeignKey/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	l "log"
	"log/slog"
	"net/http"
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

	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// handlers
	router.Get("/hello/{name}", hello.New(log))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Info("starting server", slog.String("addr", cfg.Address))
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}
