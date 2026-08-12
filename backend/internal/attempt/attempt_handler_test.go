package attempt_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/marlendd/anti-scam-trainer/internal/attempt"
	"github.com/marlendd/anti-scam-trainer/internal/auth"
	"github.com/marlendd/anti-scam-trainer/internal/scenario"
)

func TestHandler_AttemptOperations(t *testing.T) {
	t.Parallel()

	startedAt := time.Date(2026, time.August, 8, 10, 0, 0, 0, time.UTC)
	updatedAt := startedAt.Add(time.Minute)
	currentNodeID := scenario.NodeID("node-start")
	currentAttempt := attempt.Attempt{
		ID:            "attempt-1",
		UserID:        "user-1",
		ScenarioID:    "scenario-1",
		Status:        attempt.StatusInProgress,
		CurrentNodeID: &currentNodeID,
		StartedAt:     startedAt,
		UpdatedAt:     updatedAt,
	}

	tests := []struct {
		name       string
		method     string
		path       string
		operation  string
		wantStatus int
		call       func(*attempt.Handler, http.ResponseWriter, *http.Request)
	}{
		{
			name:       "starts attempt",
			method:     http.MethodPost,
			path:       "/api/v1/scenarios/scenario-1/attempts",
			operation:  "start",
			wantStatus: http.StatusCreated,
			call:       (*attempt.Handler).Start,
		},
		{
			name:       "resumes attempt",
			method:     http.MethodGet,
			path:       "/api/v1/scenarios/scenario-1/attempts/active",
			operation:  "resume",
			wantStatus: http.StatusOK,
			call:       (*attempt.Handler).Resume,
		},
		{
			name:       "restarts attempt",
			method:     http.MethodPost,
			path:       "/api/v1/scenarios/scenario-1/attempts/restart",
			operation:  "restart",
			wantStatus: http.StatusCreated,
			call:       (*attempt.Handler).Restart,
		},
		{
			name:       "gets latest completed attempt",
			method:     http.MethodGet,
			path:       "/api/v1/scenarios/scenario-1/attempts/latest",
			operation:  "get_latest_completed",
			wantStatus: http.StatusOK,
			call:       (*attempt.Handler).GetLatestCompleted,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := &answerSubmitterStub{attemptResult: currentAttempt}
			handler := attempt.NewHandler(service, discardLogger())
			request := httptest.NewRequest(test.method, test.path, nil)
			request.SetPathValue("scenarioID", "scenario-1")
			request = request.WithContext(auth.ContextWithUserID(request.Context(), "user-1"))
			response := httptest.NewRecorder()

			test.call(handler, response, request)

			require.Equal(t, test.wantStatus, response.Code)
			require.Equal(t, test.operation, service.gotOperation)
			require.Equal(t, "user-1", service.gotUserID)
			require.Equal(t, scenario.ScenarioID("scenario-1"), service.gotScenarioID)

			var payload map[string]any
			require.NoError(t, json.NewDecoder(response.Body).Decode(&payload))
			require.Equal(t, "attempt-1", payload["id"])
			require.Equal(t, "scenario-1", payload["scenario_id"])
			require.Equal(t, "in_progress", payload["status"])
			require.Equal(t, "node-start", payload["current_node_id"])
			require.NotContains(t, payload, "UserID")
			require.NotContains(t, payload, "user_id")
		})
	}
}

func TestHandler_GetState(t *testing.T) {
	t.Parallel()

	currentNodeID := scenario.NodeID("node-start")
	answeredAt := time.Date(2026, time.August, 8, 10, 0, 30, 0, time.UTC)
	service := &answerSubmitterStub{
		stateResult: attempt.State{
			Attempt: attempt.Attempt{
				ID:            "attempt-1",
				ScenarioID:    "scenario-1",
				Status:        attempt.StatusInProgress,
				CurrentNodeID: &currentNodeID,
				StartedAt:     time.Date(2026, time.August, 8, 10, 0, 0, 0, time.UTC),
				UpdatedAt:     time.Date(2026, time.August, 8, 10, 1, 0, 0, time.UTC),
			},
			Scenario: attempt.ScenarioHeader{
				Title:       "Поддельная оплата",
				Description: "Проверка оплаты внутри приложения.",
				Role:        scenario.RoleSeller,
				Product: scenario.Product{
					Title: "iPhone 15 Pro 256 ГБ",
					Price: 84990,
				},
			},
			CurrentNode: &attempt.CurrentNode{
				ID:       currentNodeID,
				Author:   "seller",
				Text:     "Выберите действие",
				Messages: []scenario.Message{{Author: "seller", Text: "Выберите действие"}},
				Choices: []attempt.ChoiceOption{
					{ID: "choice-safe", Text: "Остаться на платформе"},
				},
			},
			History: []attempt.HistoryItem{
				{
					Node: attempt.HistoryNode{
						ID:     "node-previous",
						Author: "seller",
						Text:   "Предыдущее сообщение",
						Messages: []scenario.Message{
							{Author: "buyer", Text: "Первая реплика"},
							{Author: "seller", Text: "Предыдущее сообщение"},
						},
					},
					SelectedChoice: attempt.ChoiceOption{
						ID:   "choice-previous",
						Text: "Предыдущий ответ",
					},
					Consequence: "Последствие предыдущего ответа",
					AnsweredAt:  answeredAt,
				},
			},
		},
	}
	handler := attempt.NewHandler(service, discardLogger())
	request := httptest.NewRequest(http.MethodGet, "/api/v1/attempts/attempt-1", nil)
	request.SetPathValue("attemptID", "attempt-1")
	request = request.WithContext(auth.ContextWithUserID(request.Context(), "user-1"))
	response := httptest.NewRecorder()

	handler.GetState(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "get_state", service.gotOperation)
	require.Equal(t, "user-1", service.gotUserID)
	require.Equal(t, attempt.AttemptID("attempt-1"), service.gotAttemptID)

	var payload map[string]any
	require.NoError(t, json.NewDecoder(response.Body).Decode(&payload))
	require.Equal(t, "attempt-1", payload["id"])
	require.Equal(t, map[string]any{
		"title":       "Поддельная оплата",
		"description": "Проверка оплаты внутри приложения.",
		"role":        "seller",
		"product": map[string]any{
			"title": "iPhone 15 Pro 256 ГБ",
			"price": float64(84990),
		},
	}, payload["scenario"])
	currentNode := payload["current_node"].(map[string]any)
	require.Equal(t, "node-start", currentNode["id"])
	choices := currentNode["choices"].([]any)
	choice := choices[0].(map[string]any)
	require.Equal(t, map[string]any{
		"id":   "choice-safe",
		"text": "Остаться на платформе",
	}, choice)
	require.NotContains(t, choice, "score")
	require.NotContains(t, choice, "weight")
	require.NotContains(t, choice, "consequence")
	require.NotContains(t, choice, "explanation")
	require.NotContains(t, choice, "next_node_id")
	require.NotContains(t, choice, "ending_id")
	require.Equal(t, []any{
		map[string]any{
			"author": "seller",
			"text":   "Выберите действие",
		},
	}, currentNode["messages"])
	history := payload["history"].([]any)
	require.Len(t, history, 1)
	historyItem := history[0].(map[string]any)
	require.Equal(t, map[string]any{
		"id":     "node-previous",
		"author": "seller",
		"text":   "Предыдущее сообщение",
		"messages": []any{
			map[string]any{"author": "buyer", "text": "Первая реплика"},
			map[string]any{"author": "seller", "text": "Предыдущее сообщение"},
		},
	}, historyItem["node"])
	require.Equal(t, map[string]any{
		"id":   "choice-previous",
		"text": "Предыдущий ответ",
	}, historyItem["selected_choice"])
	require.Equal(t, "Последствие предыдущего ответа", historyItem["consequence"])
	require.Equal(t, answeredAt.Format(time.RFC3339), historyItem["answered_at"])
}

func TestHandler_GetStateRejectsInvalidRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		authenticated bool
		attemptID     string
		wantStatus    int
		wantError     string
	}{
		{
			name:       "missing authenticated user",
			attemptID:  "attempt-1",
			wantStatus: http.StatusUnauthorized,
			wantError:  "unauthorized",
		},
		{
			name:          "missing attempt id",
			authenticated: true,
			wantStatus:    http.StatusBadRequest,
			wantError:     "attempt_id is required",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := &answerSubmitterStub{}
			handler := attempt.NewHandler(service, discardLogger())
			request := httptest.NewRequest(http.MethodGet, "/api/v1/attempts/placeholder", nil)
			request.SetPathValue("attemptID", test.attemptID)
			if test.authenticated {
				request = request.WithContext(auth.ContextWithUserID(request.Context(), "user-1"))
			}
			response := httptest.NewRecorder()

			handler.GetState(response, request)

			require.Equal(t, test.wantStatus, response.Code)
			require.Equal(t, test.wantError, decodeError(t, response))
			require.Empty(t, service.gotOperation)
		})
	}
}

func TestHandler_GetStateMapsNotFound(t *testing.T) {
	t.Parallel()

	service := &answerSubmitterStub{stateErr: attempt.ErrAttemptNotFound}
	handler := attempt.NewHandler(service, discardLogger())
	request := httptest.NewRequest(http.MethodGet, "/api/v1/attempts/attempt-1", nil)
	request.SetPathValue("attemptID", "attempt-1")
	request = request.WithContext(auth.ContextWithUserID(request.Context(), "user-1"))
	response := httptest.NewRecorder()

	handler.GetState(response, request)

	require.Equal(t, http.StatusNotFound, response.Code)
	require.Equal(t, "attempt not found", decodeError(t, response))
}

func TestHandler_AttemptOperationRejectsInvalidRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		authenticated bool
		scenarioID    string
		wantStatus    int
		wantError     string
	}{
		{
			name:       "missing authenticated user",
			scenarioID: "scenario-1",
			wantStatus: http.StatusUnauthorized,
			wantError:  "unauthorized",
		},
		{
			name:          "missing scenario id",
			authenticated: true,
			wantStatus:    http.StatusBadRequest,
			wantError:     "scenario_id is required",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := &answerSubmitterStub{}
			handler := attempt.NewHandler(service, discardLogger())
			request := httptest.NewRequest(
				http.MethodPost,
				"/api/v1/scenarios/placeholder/attempts",
				bytes.NewReader(nil),
			)
			request.SetPathValue("scenarioID", test.scenarioID)
			if test.authenticated {
				request = request.WithContext(auth.ContextWithUserID(request.Context(), "user-1"))
			}
			response := httptest.NewRecorder()

			handler.Start(response, request)

			require.Equal(t, test.wantStatus, response.Code)
			require.Equal(t, test.wantError, decodeError(t, response))
			require.Empty(t, service.gotOperation)
		})
	}
}

func TestHandler_AttemptOperationMapsServiceErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantError  string
	}{
		{
			name:       "scenario not found",
			err:        scenario.ErrScenarioNotFound,
			wantStatus: http.StatusNotFound,
			wantError:  "scenario not found",
		},
		{
			name:       "active attempt not found",
			err:        attempt.ErrActiveAttemptNotFound,
			wantStatus: http.StatusNotFound,
			wantError:  "attempt not found",
		},
		{
			name:       "completed attempt not found",
			err:        attempt.ErrCompletedAttemptNotFound,
			wantStatus: http.StatusNotFound,
			wantError:  "attempt not found",
		},
		{
			name:       "inactive scenario",
			err:        scenario.ErrScenarioInactive,
			wantStatus: http.StatusConflict,
			wantError:  "scenario is inactive",
		},
		{
			name:       "active attempt exists",
			err:        errors.Join(errors.New("create failed"), attempt.ErrActiveAttemptExists),
			wantStatus: http.StatusConflict,
			wantError:  "active attempt already exists",
		},
		{
			name:       "unexpected error",
			err:        errors.New("database unavailable"),
			wantStatus: http.StatusInternalServerError,
			wantError:  "internal server error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := &answerSubmitterStub{attemptErr: test.err}
			handler := attempt.NewHandler(service, discardLogger())
			request := httptest.NewRequest(
				http.MethodPost,
				"/api/v1/scenarios/scenario-1/attempts",
				nil,
			)
			request.SetPathValue("scenarioID", "scenario-1")
			request = request.WithContext(auth.ContextWithUserID(request.Context(), "user-1"))
			response := httptest.NewRecorder()

			handler.Start(response, request)

			require.Equal(t, test.wantStatus, response.Code)
			require.Equal(t, test.wantError, decodeError(t, response))
			require.Equal(t, "start", service.gotOperation)
		})
	}
}
