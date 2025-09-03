package main

import (
	"log/slog"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/wdsjk/avito-shop/internal/config"
	"github.com/wdsjk/avito-shop/internal/employee"
	"github.com/wdsjk/avito-shop/internal/infra/storage"
	"github.com/wdsjk/avito-shop/internal/infra/storage/postgres"
	"github.com/wdsjk/avito-shop/internal/infra/transport/http/handlers"
	mwlogger "github.com/wdsjk/avito-shop/internal/infra/transport/http/middleware/logger"
	"github.com/wdsjk/avito-shop/internal/infra/transport/http/server"
	"github.com/wdsjk/avito-shop/internal/shop"
	"github.com/wdsjk/avito-shop/internal/transfer"
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

	shop := shop.NewShop()

	employeeRepo := postgres.NewEmployeeRepository(storage)
	employeeService := employee.NewEmployeeService(employeeRepo)

	transferRepo := postgres.NewTransferRepository(storage)
	transferService := transfer.NewTransferService(transferRepo)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(mwlogger.New(log)) // middleware with our logger
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat) // strong coherence with chi, might want to refactor in future

	infoHandler := handlers.NewInfoHandler(employeeService, transferService)
	coinHandler := handlers.NewCoinHandler(employeeService, transferService)

	r.HandleFunc("/api/info", infoHandler.Handle)     // GET
	r.HandleFunc("/api/sendCoin", coinHandler.Handle) // POST
	r.HandleFunc("/api/buy/{item}", nil)              // GET
	r.HandleFunc("/api/auth", nil)                    // POST

	server := server.NewServer(cfg, r)
	err = server.Start(log)
	if err != nil {
		log.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	_ = employeeService
	_ = transferService
	_ = shop
	_ = server
	_ = r
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
