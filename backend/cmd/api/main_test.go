package main

import (
	"io"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRunIntegration(t *testing.T) {
	cfg := Config{
		LogLevel: "DEBUG",
		Address:  "127.0.0.1:8089",
		Timeout:  2 * time.Second,
	}

	logger := mustMakeLogger(cfg.LogLevel)

	errCh := make(chan error, 1)
	go func() {
		errCh <- run(&cfg, logger)
	}()

	time.Sleep(100 * time.Millisecond)

	select {
	case err := <-errCh:
		require.NoError(t, err, "Сервер завершился с ошибкой при старте")
	default:
	}

	res, err := http.Get("http://" + cfg.Address + "/example")

	require.NoError(t, err, "HTTP-запрос должен выполниться без ошибок")
	defer func() {
		err := res.Body.Close()
		if err != nil {
			slog.Error("error when close body")
		}
	}()

	require.Equal(t, http.StatusOK, res.StatusCode, "Ожидался статус 200 OK")

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err, "Тело ответа должно читаться без ошибок")
	require.Equal(t, "Hello world!", string(body), "Тело ответа не совпадает с ожидаемым")
}
