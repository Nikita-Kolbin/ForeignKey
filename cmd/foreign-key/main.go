package main

import (
	img "ForeignKey/internal/http-server/handlers/image"
	"ForeignKey/internal/http-server/handlers/website"
	"ForeignKey/internal/storage/image"
	"log/slog"
	"net/http"
	"os"

	_ "ForeignKey/docs"
	"ForeignKey/internal/config"
	"ForeignKey/internal/http-server/handlers/admin"
	mwLogger "ForeignKey/internal/http-server/middleware/logger"
	"ForeignKey/internal/logger"
	"ForeignKey/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	l "log"
)

// @title           ForeignKey
// @version         1.0

// @host      localhost:8082
// @BasePath  /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// config
	cfg := config.MustLoad()

	// logger
	log, err := logger.SetupLogger(cfg.Env)
	if err != nil {
		l.Fatalf("can't setup logger: %s", err)
	}
	log.Info("starting service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// storage
	storage, err := sqlite.New(cfg.StoragePath)
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

	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// TODO: указать статусы для всех ответов
	// handlers
	router.Post("/api/admin/sign-up", admin.NewSignUp(storage, log))
	router.Post("/api/admin/sign-in", admin.NewSignIn(storage, log))

	router.Post("/api/image/upload", img.NewUpload(imageSaver, storage, log))
	router.Get("/api/image/download/{id}", img.NewDownload(imageSaver, storage, log))

	router.Post("/api/website/create", website.NewCreate(storage, log))
	router.Get("/api/website/aliases", website.NewGetAliases(storage, log))

	// swagger
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8082/swagger/doc.json"),
	))

	// server
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Info("starting server", slog.String("addr", cfg.Address))
	if err = srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}
