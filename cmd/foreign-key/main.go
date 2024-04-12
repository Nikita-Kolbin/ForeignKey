package main

import (
	"ForeignKey/internal/config"
)

func main() {
	cfg := config.MustLoad()

	_ = cfg

	// TODO: init logger: slog

	// TODO: init storage: postgres

	// TODO: init router: chi, render

	// TODO: init server
}
