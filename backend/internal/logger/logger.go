package logger

import (
	"errors"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func SetupLogger(env string) (*slog.Logger, error) {
	var log *slog.Logger
	var op *slog.HandlerOptions
	var th *slog.TextHandler
	var jh *slog.JSONHandler

	switch env {
	case envLocal:
		op = &slog.HandlerOptions{Level: slog.LevelDebug}
		th = slog.NewTextHandler(os.Stdout, op)
		log = slog.New(th)
	case envDev:
		op = &slog.HandlerOptions{Level: slog.LevelDebug}
		jh = slog.NewJSONHandler(os.Stdout, op)
		log = slog.New(jh)
	case envProd:
		op = &slog.HandlerOptions{Level: slog.LevelInfo}
		jh = slog.NewJSONHandler(os.Stdout, op)
		log = slog.New(jh)
	default:
		return nil, errors.New("incorrect env")
	}

	return log, nil
}
