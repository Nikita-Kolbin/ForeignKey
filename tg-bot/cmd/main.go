package main

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"log/slog"
	"net/http"
	_ "tg-bot/docs"
	"tg-bot/internal/api"
	"tg-bot/internal/client/emailClient"
	"tg-bot/internal/client/tgClient"
	"tg-bot/internal/config"
)

const dotenvPath = "config/.env"

// @title           Notification
// @version         1.0
// @host            localhost:8083
func main() {
	cfg := config.MustLoad(dotenvPath)

	email := emailClient.New(
		cfg.EmailConfig.Email,
		cfg.EmailConfig.Password,
		cfg.EmailConfig.SMTPHost,
		cfg.EmailConfig.SMTPPort,
	)

	tg := tgClient.New(
		cfg.TgConfig.Token,
	)

	router := chi.NewRouter()

	router.Post("/send-email", api.NewSendEmail(email))
	router.Post("/send-telegram", api.NewSendTelegram(tg))

	serverAddress := cfg.ServerConfig.Host + ":" + cfg.ServerConfig.Port
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(serverAddress+"/swagger/doc.json"),
	))

	srv := &http.Server{
		Addr:    serverAddress,
		Handler: router,
	}

	log.Println("start server")
	if err := srv.ListenAndServe(); err != nil {
		log.Println("failed to start server:", slog.String("err", err.Error()))
	}
}
