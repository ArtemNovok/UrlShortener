package main

import (
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/redirect"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/handlers/slogpretty"
	"url-shortener/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envProd  = "prod"
	envLocal = "local"
	envDev   = "dev"
)

type Storage struct {
	urlSaver   save.URLSaver
	urlGetter  redirect.URLGetter
	urlDeleter delete.URLDeleter
}

func main() {
	// init config
	cfg := config.MustLoad()
	// init logger
	log := setupLogger(cfg.Env)
	log.Info("starting service", slog.String("env", cfg.Env))
	// init storage
	strg, err := postgres.New(cfg.StorageConn, *log)
	if err != nil {
		log.Error("failed to connect to storage", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.Info("connected to storage and tables created")
	// init router
	storage := Storage{
		urlSaver:   strg,
		urlGetter:  strg,
		urlDeleter: strg,
	}
	router := setupRouter(log, storage, *cfg)

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.TimeOut,
		WriteTimeout: cfg.HttpServer.TimeOut,
		IdleTimeout:  cfg.HttpServer.IdleTimeOut,
	}
	log.Info("starting server on address", slog.String("address", cfg.HttpServer.Address))
	if err := srv.ListenAndServe(); err != nil {
		log.Error("Fatal error", slog.String("error", err.Error()))
	}
	log.Info("server stopped")
	// init server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = SetupPrettyLogger()
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	return log
}

func setupRouter(log *slog.Logger, saver Storage, cfg config.Config) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Get("/{alias}", redirect.New(log, saver.urlGetter))
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HttpServer.Email: cfg.HttpServer.Password,
		}))
		r.Delete("/", delete.New(log, saver.urlDeleter))
		r.Post("/", save.New(log, saver.urlSaver))
	})
	return router
}

func SetupPrettyLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
