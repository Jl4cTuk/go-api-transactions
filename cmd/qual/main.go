package main

import (
	"infotex/internal/config"
	"infotex/internal/http-server/handlers/url/getbalance"
	"infotex/internal/http-server/handlers/url/getlast"
	"infotex/internal/http-server/handlers/url/send"
	"infotex/internal/lib/logger/handlers/slogpretty"
	"infotex/internal/lib/logger/sl"
	"infotex/internal/storage/postgresql"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {

	// init config
	cfg := config.MustLoad()

	// init logger
	log := setupLogger(cfg.Env)
	log.Info("starting api", slog.String("env", cfg.Env))

	// init storage
	storage, err := postgresql.New(cfg.DBServer)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	// init wallets
	err = storage.GenRandomWallet(10)
	if err != nil {
		log.Error("failed to create 10 wallets", sl.Err(err))
	}

	// init router
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/api/send", send.New(log, storage))
	router.Get("/api/transactions", getlast.New(log, storage))
	router.Get("/api/wallet/{address}/balance", getbalance.New(log, storage))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// start server
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = setupPrettySlog()
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

// setupPrettySlog initializes a pretty logger for development environment.
// For debugging purposes, it outputs logs with colors and additional context.
func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
