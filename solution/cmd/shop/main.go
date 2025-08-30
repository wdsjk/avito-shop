package main

import (
	"log/slog"
	"os"

	"github.com/wdsjk/avito-shop/internal/config"
	"github.com/wdsjk/avito-shop/internal/employee"
	"github.com/wdsjk/avito-shop/internal/infra/storage"
	"github.com/wdsjk/avito-shop/internal/infra/storage/postgres"
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

	storage, err := storage.NewStorage(cfg)
	if err != nil {
		log.Error("failed to initialize storage", "error", err)
		os.Exit(1)
	}

	employeeRepo := postgres.NewEmployeeRepository(storage)
	employeeService := employee.NewEmployeeService(employeeRepo)

	// TODO: shop repo, service; transfer repo, service
	_ = employeeService
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
