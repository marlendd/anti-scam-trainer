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

type attemptStateAPIResponse struct {
	attemptAPIResponse
	CurrentNode *struct {
		ID      string `json:"id"`
		Author  string `json:"author"`
		Text    string `json:"text"`
		Choices []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"choices"`
	} `json:"current_node"`
}

type submitAnswerAPIResponse struct {
	AttemptID        string  `json:"attempt_id"`
	NodeID           string  `json:"node_id"`
	ChoiceID         string  `json:"choice_id"`
	NextNodeID       *string `json:"next_node_id"`
	EndingID         *string `json:"ending_id"`
	Completed        bool    `json:"completed"`
	Score            *int    `json:"score"`
	RewardFragmentID *string `json:"reward_fragment_id"`
}

type scenarioCatalogAPIItem struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type scenarioCatalogAPIResponse struct {
	Items []scenarioCatalogAPIItem `json:"items"`
}

const buyerGPUSeedScenarioID = "45d4cc8c-f604-4a7c-b8c5-f2464717b71f"

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
		SeedsPath:       "../../seeds",
		MaxOpenConns:    5,
		MaxIdleConns:    5,
		ConnMaxLifetime: 1 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
		SecureCookies:   false,
	}

	logger := mustMakeLogger(cfg.LogLevel)

	require.NoError(t, postgres.RunMigrations(databaseURL, cfg.MigrationsPath))
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

	t.Run("catalog exposes two seeded scenarios per role", func(t *testing.T) {
		buyerCatalog := getScenarioCatalog(t, client, baseURL, scenario.RoleBuyer)
		require.Len(t, buyerCatalog.Items, 2)

		sellerCatalog := getScenarioCatalog(t, client, baseURL, scenario.RoleSeller)
		require.Len(t, sellerCatalog.Items, 2)
	})

	t.Run("completes loaded seed through HTTP API", func(t *testing.T) {
		catalogItem := getScenarioCatalogItem(t, client, baseURL, scenario.RoleBuyer, buyerGPUSeedScenarioID)
		require.Equal(t, "not_started", catalogItem.Status)

		startResponse, err := client.Post(
			baseURL+"/api/v1/scenarios/"+buyerGPUSeedScenarioID+"/attempts",
			"application/json",
			nil,
		)
		require.NoError(t, err)
		defer closeBody(t, startResponse)
		require.Equal(t, http.StatusCreated, startResponse.StatusCode)

		var started attemptAPIResponse
		require.NoError(t, json.NewDecoder(startResponse.Body).Decode(&started))
		require.Equal(t, buyerGPUSeedScenarioID, started.ScenarioID)
		require.Equal(t, "in_progress", started.Status)
		require.NotNil(t, started.CurrentNodeID)
		require.Equal(t, "n1_scarcity_pressure", *started.CurrentNodeID)

		steps := []struct {
			nodeID     string
			choiceID   string
			nextNodeID string
		}{
			{
				nodeID:     "n1_scarcity_pressure",
				choiceID:   "n1_take_time_to_review",
				nextNodeID: "n2_platform_issue",
			},
			{
				nodeID:     "n2_platform_issue",
				choiceID:   "n2p_contact_support_independently",
				nextNodeID: "n3_official_verification",
			},
			{
				nodeID:     "n3_official_verification",
				choiceID:   "n3v_refuse_and_report",
				nextNodeID: "n4_protected",
			},
		}

		for _, step := range steps {
			result := submitAnswerThroughAPI(
				t,
				client,
				baseURL,
				testDB,
				started.ID,
				step.nodeID,
				step.choiceID,
			)
			require.False(t, result.Completed)
			require.NotNil(t, result.NextNodeID)
			require.Equal(t, step.nextNodeID, *result.NextNodeID)
		}

		final := submitAnswerThroughAPI(
			t,
			client,
			baseURL,
			testDB,
			started.ID,
			"n4_protected",
			"n4p_block_and_report",
		)
		require.True(t, final.Completed)
		require.Nil(t, final.NextNodeID)
		require.NotNil(t, final.EndingID)
		require.Equal(t, "ending_safe", *final.EndingID)
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
		require.Equal(t, 4, answerCount)

		catalogItem = getScenarioCatalogItem(
			t,
			client,
			baseURL,
			scenario.RoleBuyer,
			buyerGPUSeedScenarioID,
		)
		require.Equal(t, "completed", catalogItem.Status)
	})

	scenarioID := insertAPITestScenario(t, testDB)
	t.Cleanup(func() {
		_, cleanupErr := testDB.Exec(`DELETE FROM user_inventory WHERE scenario_id = $1`, scenarioID)
		require.NoError(t, cleanupErr)
		_, cleanupErr = testDB.Exec(`DELETE FROM attempts WHERE scenario_id = $1`, scenarioID)
		require.NoError(t, cleanupErr)
		_, cleanupErr = testDB.Exec(`DELETE FROM scenario_versions WHERE id = $1`, scenarioID)
		require.NoError(t, cleanupErr)
	})

	t.Run("completes scenario through HTTP API", func(t *testing.T) {
		catalogItem := getScenarioCatalogItem(t, client, baseURL, scenario.RoleBuyer, scenarioID)
		require.Equal(t, "not_started", catalogItem.Status)

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

		state := getAttemptStateThroughAPI(t, client, baseURL, started.ID)
		require.Equal(t, started.ID, state.ID)
		require.NotNil(t, state.CurrentNode)
		require.Equal(t, string(testfixture.StartNodeID), state.CurrentNode.ID)
		require.NotEmpty(t, state.CurrentNode.Choices)
		require.Equal(t, string(testfixture.StartChoiceID), state.CurrentNode.Choices[0].ID)

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

		state = getAttemptStateThroughAPI(t, client, baseURL, started.ID)
		require.NotNil(t, state.CurrentNode)
		require.Equal(t, string(testfixture.MiddleNodeID), state.CurrentNode.ID)

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
		require.NotNil(t, final.RewardFragmentID)
		require.Equal(t, string(testfixture.RewardFragmentID), *final.RewardFragmentID)

		var status string
		var score int
		var answerCount int
		var fragmentCount int
		require.NoError(t, testDB.QueryRow(`
			SELECT status,
			       score,
			       (SELECT count(*) FROM answers WHERE attempt_id = attempts.id),
			       (SELECT count(*)
			        FROM user_inventory
			        WHERE user_id = attempts.user_id
			          AND scenario_id = attempts.scenario_id)
			FROM attempts
			WHERE id = $1
		`, started.ID).Scan(&status, &score, &answerCount, &fragmentCount))
		require.Equal(t, "completed", status)
		require.Equal(t, 100, score)
		require.Equal(t, 3, answerCount)
		require.Equal(t, 1, fragmentCount)

		catalogItem = getScenarioCatalogItem(t, client, baseURL, scenario.RoleBuyer, scenarioID)
		require.Equal(t, "completed", catalogItem.Status)

		state = getAttemptStateThroughAPI(t, client, baseURL, started.ID)
		require.Equal(t, "completed", state.Status)
		require.Nil(t, state.CurrentNode)
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

func getAttemptStateThroughAPI(
	t *testing.T,
	client *http.Client,
	baseURL string,
	attemptID string,
) attemptStateAPIResponse {
	t.Helper()

	response, err := client.Get(baseURL + "/api/v1/attempts/" + attemptID)
	require.NoError(t, err)
	defer closeBody(t, response)
	require.Equal(t, http.StatusOK, response.StatusCode)

	var state attemptStateAPIResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&state))

	return state
}

func getScenarioCatalogItem(
	t *testing.T,
	client *http.Client,
	baseURL string,
	role scenario.Role,
	scenarioID string,
) scenarioCatalogAPIItem {
	t.Helper()

	catalog := getScenarioCatalog(t, client, baseURL, role)
	for _, item := range catalog.Items {
		if item.ID == scenarioID {
			return item
		}
	}

	t.Fatalf("scenario %s not found in catalog", scenarioID)
	return scenarioCatalogAPIItem{}
}

func getScenarioCatalog(
	t *testing.T,
	client *http.Client,
	baseURL string,
	role scenario.Role,
) scenarioCatalogAPIResponse {
	t.Helper()

	response, err := client.Get(baseURL + "/api/v1/scenarios?role=" + string(role))
	require.NoError(t, err)
	defer closeBody(t, response)
	require.Equal(t, http.StatusOK, response.StatusCode)

	var catalog scenarioCatalogAPIResponse
	require.NoError(t, json.NewDecoder(response.Body).Decode(&catalog))
	return catalog
}

func insertAPITestScenario(t *testing.T, db interface {
	QueryRow(query string, args ...any) *sql.Row
}) string {
	t.Helper()

	fixture := testfixture.ValidScenario()
	content, err := json.Marshal(scenario.Content{
		StartNodeID:         fixture.StartNodeID,
		SuccessfulEndingIDs: fixture.SuccessfulEndingIDs,
		Nodes:               fixture.Nodes,
		Endings:             fixture.Endings,
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
			reward_fragment_id,
			content
		)
		VALUES (gen_random_uuid(), 1, $1, $2, $3, TRUE, $4, $5::jsonb)
		RETURNING id
	`, fixture.Role, fixture.Title, fixture.Description, fixture.RewardFragmentID, content).Scan(&scenarioID))

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
}

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
