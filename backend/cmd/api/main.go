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

	"github.com/marlendd/anti-scam-trainer/internal/attempt"
	"github.com/marlendd/anti-scam-trainer/internal/auth"
	"github.com/marlendd/anti-scam-trainer/internal/evaluation"
	"github.com/marlendd/anti-scam-trainer/internal/platform/config"
	"github.com/marlendd/anti-scam-trainer/internal/platform/health"
	"github.com/marlendd/anti-scam-trainer/internal/platform/mailer"
	"github.com/marlendd/anti-scam-trainer/internal/platform/middleware"
	"github.com/marlendd/anti-scam-trainer/internal/platform/postgres"
	"github.com/marlendd/anti-scam-trainer/internal/progress"
	"github.com/marlendd/anti-scam-trainer/internal/scenario"
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

	// ---------- wiring attempts ----------
	scenarioRepository := scenario.NewPgRepository(db)
	scenarioCatalogService := scenario.NewCatalogService(&scenarioRepository)
	scenarioCatalogHandler := scenario.NewCatalogHandler(scenarioCatalogService, log)
	attemptRepository := attempt.NewPgRepository(db)
	attemptService := attempt.NewService(
		attemptRepository,
		attemptRepository,
		&scenarioRepository,
	)
	attemptHandler := attempt.NewHandler(attemptService, log)

	// ---------- evaluation ----------
	evaluator := evaluation.NewEvaluator()

	// ---------- progress ----------
	progressRepo := progress.NewPgRepository(db, log)
	progressService := progress.NewService(progressRepo, evaluator)
	progressHandler := progress.NewHandler(progressService, log)

	mux := http.NewServeMux()

	// health/ready
	mux.HandleFunc("GET /health", health.Health)
	mux.HandleFunc("GET /ready", health.Ready(db))

	// auth routes
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/auth/logout", authHandler.Logout)
	mux.HandleFunc("POST /api/v1/auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("POST /api/v1/auth/reset-password", authHandler.ResetPassword)

	handler := middleware.CORS(cfg.AllowedOrigins)(mux)

	// protected routes
	mux.Handle("GET /api/v1/users/me", requireAuth(http.HandlerFunc(authHandler.Me)))
	mux.Handle(
		"GET /api/v1/scenarios",
		requireAuth(http.HandlerFunc(scenarioCatalogHandler.List)),
	)
	mux.Handle(
		"POST /api/v1/scenarios/{scenarioID}/attempts",
		requireAuth(http.HandlerFunc(attemptHandler.Start)),
	)
	mux.Handle(
		"GET /api/v1/scenarios/{scenarioID}/attempts/active",
		requireAuth(http.HandlerFunc(attemptHandler.Resume)),
	)
	mux.Handle(
		"POST /api/v1/scenarios/{scenarioID}/attempts/restart",
		requireAuth(http.HandlerFunc(attemptHandler.Restart)),
	)
	mux.Handle(
		"POST /api/v1/attempts/{attemptID}/answers",
		requireAuth(http.HandlerFunc(attemptHandler.SubmitAnswer)),
	)

	// protected routes
	mux.Handle("GET /api/v1/users/me", requireAuth(http.HandlerFunc(authHandler.Me)))

	// progress not protected routes
	mux.HandleFunc("GET /api/v1/leaderboard", progressHandler.GetLeaderboard)
	// progress protected routes
	mux.Handle("GET /api/v1/profile/role-progress", requireAuth(http.HandlerFunc(progressHandler.GetMyRoleStats)))
	mux.Handle("GET /api/v1/profile/categories-progress", requireAuth(http.HandlerFunc(progressHandler.GetMyCategoryDashboard)))
	mux.Handle("GET /api/v1/profile/puzzle", requireAuth(http.HandlerFunc(progressHandler.GetMyPuzzleProgress)))
	mux.Handle("GET /api/v1/attempts/{id}/result", requireAuth(http.HandlerFunc(progressHandler.GetStatsOfAttempt)))
	mux.Handle("GET /api/v1/profile/rank-history", requireAuth(http.HandlerFunc(progressHandler.GetMyRankHistory)))
	addr := ":" + cfg.Port
	if cfg.Port == "" {
		addr = ":8080"
	}

	server := &http.Server{
		Addr:        addr,
		ReadTimeout: cfg.Timeout,
		Handler:     handler,
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
