package main

import (
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/storage/postgres"
)

const (
	envProd  = "prod"
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	// init config
	cfg := config.MustLoad()
	// init logger
	log := setupLogger(cfg.Env)
	log.Info("starting service", slog.String("env", cfg.Env))
	// init storage
	strg, err := postgres.New(cfg.StorageConn)
	if err != nil {
		log.Error("failed to connect to storage", slog.String("error", err.Error()))
		os.Exit(1)
	}
	_ = strg
	log.Info("connected to storage and tables created")
	// init router

	// init server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	return log
}
