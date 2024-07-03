package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	TgConfig     TgConfig
	EmailConfig  EmailConfig
	ServerConfig ServerConfig
}

type TgConfig struct {
	Token string `env:"TG_TOKEN"`
}

type EmailConfig struct {
	Email    string `env:"EMAIL"`
	Password string `env:"EMAIL_PASSWORD"`
	SMTPHost string `env:"SMTP_HOST"`
	SMTPPort string `env:"SMTP_PORT"`
}

type ServerConfig struct {
	Host string `env:"SERVER_HOST"`
	Port string `env:"SERVER_PORT"`
}

func MustLoad(dotenvPath string) Config {
	if os.Getenv("IN_DOCKER") == "" {
		err := godotenv.Load(dotenvPath)
		if err != nil {
			panic(err)
		}
	}

	cfg := Config{}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}

	return cfg
}
