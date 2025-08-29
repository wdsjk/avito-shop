package main

import (
	"log/slog"
	"os"

	"github.com/wdsjk/avito-shop/internal/config"
	"github.com/wdsjk/avito-shop/internal/storage/postgres"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting avito-shop service", "env", cfg.Env)
	log.Debug("debug mode is enabled")

	storage, err := postgres.New(
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbName,
		cfg.DbHost,
		cfg.DbPort,
	)
	if err != nil {
		log.Error("failed to initialize storage", "error", err)
		os.Exit(1)
	}

	_ = storage
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
