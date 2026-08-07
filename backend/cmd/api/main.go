package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/marlendd/anti-scam-trainer/internal/auth"
	"github.com/marlendd/anti-scam-trainer/internal/evaluation"
	"github.com/marlendd/anti-scam-trainer/internal/platform/config"
	"github.com/marlendd/anti-scam-trainer/internal/platform/mailer"
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

	// ---------- wiring mailer ----------
	m := mailer.New(mailer.Config{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUsername,
		Password: cfg.SMTPPassword,
		From:     cfg.SMTPFrom,
	})

	// ---------- wiring auth ----------
	userRepo := auth.NewPgUserRepository(db)
	sessionRepo := auth.NewPgSessionRepository(db)
	passwordResetRepo := auth.NewPgPasswordResetRepository(db)

	authService := auth.NewService(userRepo, sessionRepo, passwordResetRepo, m, cfg.AppBaseURL)
	authHandler := auth.NewHandler(authService, log, cfg.SecureCookies)
	requireAuth := auth.RequireAuth(authService, log)
	// evaluation
	evalRepo := evaluation.NewPgRepository(db)
	evalService := evaluation.NewService(evalRepo)
	evalHandler := evaluation.NewHandler(evalService, log)

	mux := http.NewServeMux()

	// health/ready
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// auth routes
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/logout", authHandler.Logout)
	mux.HandleFunc("POST /api/v1/auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("POST /api/v1/auth/reset-password", authHandler.ResetPassword)
	// eval routes
	mux.HandleFunc("GET  /api/v1/attempts/{id}/result", evalHandler.GetStatsOfAttempt)
	mux.HandleFunc("GET  /api/v1/profile/progress", evalHandler.GetGlobalStatsHandler)
	mux.HandleFunc("GET /api/v1/evaluation/categories", evalHandler.GetCategoryStats)

	// protected routes
	mux.Handle("GET /api/v1/users/me", requireAuth(http.HandlerFunc(authHandler.Me)))
	mux.Handle("GET /api/v1/profile/puzzle", requireAuth(http.HandlerFunc(evalHandler.GetMyPuzzleProgress)))
	addr := ":" + cfg.Port
	if cfg.Port == "" {
		addr = ":8080"
	}

	server := &http.Server{
		Addr:        addr,
		ReadTimeout: cfg.Timeout,
		Handler:     mux,
	}

	// ---------- graceful shutdown ----------
	serverErr := make(chan error, 1)
	go func() {
		log.Info("server run success", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- fmt.Errorf("server closed unexpectedly: %w", err)
			return
		}
		serverErr <- nil
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case sig := <-stop:
		log.Info("shutdown signal received", "signal", sig.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	log.Info("server stopped gracefully")
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
