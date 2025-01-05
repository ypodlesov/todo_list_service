package main

import (
	"log/slog"
	"os"
	"todo_list_service/pkg/config"
)

func main() {
	cfg := config.MustLoad()
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	log.Info(
		"starting todo_list service",
		slog.String("env", cfg.Env),
		slog.String("address", cfg.HTTPServer.Address),
	)

}
