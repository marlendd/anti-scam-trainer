package main

import (
	"io"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/marlendd/anti-scam-trainer/internal/platform/config"
	"github.com/stretchr/testify/require"
)

func TestRunIntegration(t *testing.T) {
	cfg := config.Config{
		LogLevel:        "DEBUG",
		Port:            "8089",
		Timeout:         2 * time.Second,
		DatabaseURL:     "postgres://postgres:postgres@127.0.0.1:5433/antiscam_test?sslmode=disable",
		MigrationsPath:  "../../migrations",
		MaxOpenConns:    5,
		MaxIdleConns:    5,
		ConnMaxLifetime: 1 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}

	logger := mustMakeLogger(cfg.LogLevel)

	errCh := make(chan error, 1)
	go func() {
		errCh <- run(&cfg, logger)
	}()

	// Даем серверу и БД время на запуск
	time.Sleep(1 * time.Second)

	select {
	case err := <-errCh:
		t.Fatalf("Сервер завершился с ошибкой при старте: %v", err)
	default:
	}

	url := "http://127.0.0.1:" + cfg.Port + "/example"
	res, err := http.Get(url)
	require.NoError(t, err, "HTTP-запрос должен выполниться без ошибок")
	defer func() {
		if err := res.Body.Close(); err != nil {
			slog.Error("error when close body", "error", err)
		}
	}()

	require.Equal(t, http.StatusOK, res.StatusCode, "Ожидался статус 200 OK")

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err, "Тело ответа должно читаться без ошибок")
	require.Equal(t, "Hello world!", string(body), "Тело ответа не совпадает с ожидаемым")
}
