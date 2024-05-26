package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	ImagesPath  string `yaml:"images_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
	Email       `yaml:"email"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8082"`
	Timeout     time.Duration `yaml:"timeout" env-default:"10s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Email struct {
	EmailAddress string `yaml:"email_address"`
	Password     string `yaml:"password"`
	SmtpHost     string `yaml:"smtp_host"`
	SmtpPort     string `yaml:"smtp_port"`
}

func MustLoad() *Config {
	//configPath := os.Getenv("CONFIG_PATH")
	//if configPath == "" {
	//	log.Fatal("CONFIG_PATH is not set")
	//}

	configPath := "./config/local.yaml"

	// check file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("connfig file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can't read config: %s", err)
	}

	return &cfg
}
