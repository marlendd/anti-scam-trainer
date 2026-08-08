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
	"github.com/marlendd/anti-scam-trainer/internal/platform/postgres"
	"github.com/marlendd/anti-scam-trainer/internal/scenario"
	"github.com/marlendd/anti-scam-trainer/internal/testfixture"
	"github.com/stretchr/testify/require"
)

type attemptAPIResponse struct {
	ID            string  `json:"id"`
	ScenarioID    string  `json:"scenario_id"`
	Status        string  `json:"status"`
	CurrentNodeID *string `json:"current_node_id"`
}

type submitAnswerAPIResponse struct {
	AttemptID  string  `json:"attempt_id"`
	NodeID     string  `json:"node_id"`
	ChoiceID   string  `json:"choice_id"`
	NextNodeID *string `json:"next_node_id"`
	EndingID   *string `json:"ending_id"`
	Completed  bool    `json:"completed"`
	Score      *int    `json:"score"`
}

func TestRunIntegration_APIFlow(t *testing.T) {
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
	testDB, err := postgres.NewDB(cfg)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, testDB.Close())
	})

	scenarioID := insertAPITestScenario(t, testDB)
	t.Cleanup(func() {
		_, cleanupErr := testDB.Exec(`DELETE FROM attempts WHERE scenario_id = $1`, scenarioID)
		require.NoError(t, cleanupErr)
		_, cleanupErr = testDB.Exec(`DELETE FROM scenario_versions WHERE id = $1`, scenarioID)
		require.NoError(t, cleanupErr)
	})

	jar, err := cookiejar.New(nil)
	require.NoError(t, err)

	client := &http.Client{Jar: jar}

	email := fmt.Sprintf("test_%d@example.com", time.Now().UnixNano())
	password := "password123"

	t.Run("register", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"name":     "Test User",
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

	t.Run("completes scenario through HTTP API", func(t *testing.T) {
		startResponse, err := client.Post(
			baseURL+"/api/v1/scenarios/"+scenarioID+"/attempts",
			"application/json",
			nil,
		)
		require.NoError(t, err)
		defer closeBody(t, startResponse)
		require.Equal(t, http.StatusCreated, startResponse.StatusCode)

		var started attemptAPIResponse
		require.NoError(t, json.NewDecoder(startResponse.Body).Decode(&started))
		require.Equal(t, scenarioID, started.ScenarioID)
		require.Equal(t, "in_progress", started.Status)
		require.NotNil(t, started.CurrentNodeID)
		require.Equal(t, string(testfixture.StartNodeID), *started.CurrentNodeID)

		first := submitAnswerThroughAPI(
			t,
			client,
			baseURL,
			testDB,
			started.ID,
			string(testfixture.StartNodeID),
			string(testfixture.StartChoiceID),
		)
		require.False(t, first.Completed)
		require.NotNil(t, first.NextNodeID)
		require.Equal(t, string(testfixture.MiddleNodeID), *first.NextNodeID)

		second := submitAnswerThroughAPI(
			t,
			client,
			baseURL,
			testDB,
			started.ID,
			string(testfixture.MiddleNodeID),
			"middle-choice-1",
		)
		require.False(t, second.Completed)
		require.NotNil(t, second.NextNodeID)
		require.Equal(t, string(testfixture.FinalNodeID), *second.NextNodeID)

		final := submitAnswerThroughAPI(
			t,
			client,
			baseURL,
			testDB,
			started.ID,
			string(testfixture.FinalNodeID),
			string(testfixture.FinalChoiceID),
		)
		require.True(t, final.Completed)
		require.Nil(t, final.NextNodeID)
		require.NotNil(t, final.EndingID)
		require.Equal(t, string(testfixture.SafeEndingID), *final.EndingID)
		require.NotNil(t, final.Score)
		require.Equal(t, 100, *final.Score)

		var status string
		var score int
		var answerCount int
		require.NoError(t, testDB.QueryRow(`
			SELECT status,
			       score,
			       (SELECT count(*) FROM answers WHERE attempt_id = attempts.id)
			FROM attempts
			WHERE id = $1
		`, started.ID).Scan(&status, &score, &answerCount))
		require.Equal(t, "completed", status)
		require.Equal(t, 100, score)
		require.Equal(t, 3, answerCount)
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
			"name":     "Test User",
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

func insertAPITestScenario(t *testing.T, db interface {
	QueryRow(query string, args ...any) *sql.Row
}) string {
	t.Helper()

	fixture := testfixture.ValidScenario()
	content, err := json.Marshal(scenario.Content{
		StartNodeID: fixture.StartNodeID,
		Nodes:       fixture.Nodes,
		Endings:     fixture.Endings,
	})
	require.NoError(t, err)

	var scenarioID string
	require.NoError(t, db.QueryRow(`
		INSERT INTO scenario_versions (
			logical_id,
			version,
			role,
			title,
			description,
			is_active,
			content
		)
		VALUES (gen_random_uuid(), 1, $1, $2, $3, TRUE, $4::jsonb)
		RETURNING id
	`, fixture.Role, fixture.Title, fixture.Description, content).Scan(&scenarioID))

	return scenarioID
}

func submitAnswerThroughAPI(
	t *testing.T,
	client *http.Client,
	baseURL string,
	db interface {
		QueryRow(query string, args ...any) *sql.Row
	},
	attemptID string,
	nodeID string,
	choiceID string,
) submitAnswerAPIResponse {
	t.Helper()

	var idempotencyKey string
	require.NoError(t, db.QueryRow(`SELECT gen_random_uuid()`).Scan(&idempotencyKey))

	body, err := json.Marshal(map[string]string{
		"node_id":         nodeID,
		"choice_id":       choiceID,
		"idempotency_key": idempotencyKey,
	})
	require.NoError(t, err)

	response, err := client.Post(
		baseURL+"/api/v1/attempts/"+attemptID+"/answers",
		"application/json",
		bytes.NewReader(body),
	)
	require.NoError(t, err)
	defer closeBody(t, response)
	require.Equal(t, http.StatusOK, response.StatusCode)

	var result submitAnswerAPIResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&result))
	require.Equal(t, attemptID, result.AttemptID)
	require.Equal(t, nodeID, result.NodeID)
	require.Equal(t, choiceID, result.ChoiceID)

	return result
func cleanupTestDatabase(t *testing.T, databaseURL string) {
	t.Helper()

	db, err := sql.Open("pgx", databaseURL)
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
