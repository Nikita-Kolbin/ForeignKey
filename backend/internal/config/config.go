package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `env:"ENV" env-required:"true"`
	StoragePath string `env:"STORAGE_PATH" env-required:"true"`
	StorageName string `env:"STORAGE_NAME" env-required:"true"`
	ImagesPath  string `env:"IMAGES_PATH" env-required:"true"`
	HTTPServer
	NotificationClientConfig NotificationClientConfig
}

type HTTPServer struct {
	Address     string        `env:"SERVER_ADDRESS" env-default:"localhost:8082"`
	Timeout     time.Duration `env:"SERVER_TIMEOUT" env-default:"10s"`
	IdleTimeout time.Duration `env:"SERVER_IDLE_TIMEOUT" env-default:"60s"`
}

type NotificationClientConfig struct {
	EmailURL    string `env:"NOTIFICATION_EMAIL_URL"`
	TelegramURL string `env:"NOTIFICATION_TELEGRAM_URL"`
}

func MustLoad(dotenvPath string) *Config {
	if os.Getenv("IN_DOCKER") == "" {
		err := godotenv.Load(dotenvPath)
		if err != nil {
			log.Fatalf("can't load dotenv: %s", err)
		}
	}

	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("can't read config: %s", err)
	}

	return &cfg
}
