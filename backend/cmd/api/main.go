package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/marlendd/anti-scam-trainer/internal/platform/config"
	"github.com/marlendd/anti-scam-trainer/internal/platform/postgres"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := mustMakeLogger(cfg.LogLevel)

	if err := run(&cfg, log); err != nil {
		log.Error("server stopped with error", "error", err)
		os.Exit(1)
	}
}

func run(cfg *config.Config, log *slog.Logger) error {
	db, err := postgres.NewDB(*cfg)
	if err != nil {
		return fmt.Errorf("failed to init postgres: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error("failed to close postgres connection", "error", err)
		}
	}()

	log.Info("database connection established successfully")

	if err := postgres.RunMigrations(cfg.DatabaseURL, cfg.MigrationsPath); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	
	mux := http.NewServeMux()

	mux.HandleFunc("GET /example", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintf(w, "Hello world!"); err != nil {
			log.Error("error when write bytes", "error", err)
		}
	})

	addr := ":" + cfg.Port
	if cfg.Port == "" {
		addr = ":8080"
	}

	server := http.Server{
		Addr:        addr,
		ReadTimeout: cfg.Timeout,
		Handler:     mux,
	}

	log.Info("server run success", "addr", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server closed unexpectedly: %w", err)
		}
	}

	return nil
}

func mustMakeLogger(logLevel string) *slog.Logger {
	var level slog.Level
	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "ERROR":
		level = slog.LevelError
	default:
		panic("unknown log level: " + logLevel)
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})

	return slog.New(handler)
}
