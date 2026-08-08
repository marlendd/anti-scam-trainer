package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"os"
	"testing"
	"time"

	"github.com/marlendd/anti-scam-trainer/internal/platform/config"
	"github.com/stretchr/testify/require"
)

func TestRunIntegration_AuthFlow(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("skipping integration test; set RUN_INTEGRATION_TESTS=1 to run")
	}

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@127.0.0.1:5433/antiscam?sslmode=disable"
	}

	cfg := config.Config{
		LogLevel:        "DEBUG",
		Port:            "8089",
		Timeout:         2 * time.Second,
		DatabaseURL:     databaseURL,
		MigrationsPath:  "../../migrations",
		MaxOpenConns:    5,
		MaxIdleConns:    5,
		ConnMaxLifetime: 1 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
		SecureCookies:   false,
	}

	logger := mustMakeLogger(cfg.LogLevel)

	cleanupTestDatabase(t, databaseURL)

	errCh := make(chan error, 1)
	go func() {
		errCh <- run(&cfg, logger)
	}()

	time.Sleep(1 * time.Second)

	select {
	case err := <-errCh:
		t.Fatalf("Сервер завершился с ошибкой при старте: %v", err)
	default:
	}

	baseURL := "http://127.0.0.1:" + cfg.Port

	jar, err := cookiejar.New(nil)
	require.NoError(t, err)

	client := &http.Client{Jar: jar}

	email := fmt.Sprintf("test_%d@example.com", time.Now().UnixNano())
	password := "password123"

	t.Run("register", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"email":    email,
			"password": password,
		})

		res, err := client.Post(
			baseURL+"/api/v1/auth/register",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer closeBody(t, res)

		require.Equal(t, http.StatusCreated, res.StatusCode)
	})

	t.Run("login sets session cookie", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"email":    email,
			"password": password,
		})

		res, err := client.Post(
			baseURL+"/api/v1/auth/login",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer closeBody(t, res)

		require.Equal(t, http.StatusOK, res.StatusCode)

		var cookieFound bool
		for _, c := range res.Cookies() {
			if c.Name == "session_id" {
				cookieFound = true
				require.True(t, c.HttpOnly, "cookie должна быть HttpOnly")
			}
		}

		require.True(t, cookieFound, "ожидалась cookie session_id после логина")
	})

	t.Run("me returns current user with valid session", func(t *testing.T) {
		res, err := client.Get(baseURL + "/api/v1/users/me")
		require.NoError(t, err)
		defer closeBody(t, res)

		require.Equal(t, http.StatusOK, res.StatusCode)

		var payload struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		}

		require.NoError(t, json.NewDecoder(res.Body).Decode(&payload))
		require.Equal(t, email, payload.Email)
	})

	t.Run("logout clears session", func(t *testing.T) {
		res, err := client.Post(
			baseURL+"/api/v1/auth/logout",
			"application/json",
			nil,
		)
		require.NoError(t, err)
		defer closeBody(t, res)

		require.Equal(t, http.StatusNoContent, res.StatusCode)
	})

	t.Run("me returns 401 after logout", func(t *testing.T) {
		res, err := client.Get(baseURL + "/api/v1/users/me")
		require.NoError(t, err)
		defer closeBody(t, res)

		require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("login with wrong password fails", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"email":    email,
			"password": "wrong-password",
		})

		res, err := client.Post(
			baseURL+"/api/v1/auth/login",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer closeBody(t, res)

		require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("register with existing email fails", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"email":    email,
			"password": password,
		})

		res, err := client.Post(
			baseURL+"/api/v1/auth/register",
			"application/json",
			bytes.NewReader(body),
		)
		require.NoError(t, err)
		defer closeBody(t, res)

		require.Equal(t, http.StatusConflict, res.StatusCode)
	})
}

func cleanupTestDatabase(t *testing.T, databaseURL string) {
	t.Helper()

	db, err := sql.Open("postgres", databaseURL)
	require.NoError(t, err)

	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close test database: %v", err)
		}
	}()

	require.NoError(t, db.Ping())

	_, err = db.Exec(`
		TRUNCATE TABLE attempts, sessions CASCADE
	`)
	require.NoError(t, err)
}

func closeBody(t *testing.T, res *http.Response) {
	t.Helper()

	if err := res.Body.Close(); err != nil {
		slog.Error("error when close body", "error", err)
	}
}
