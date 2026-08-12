package attempt_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/marlendd/anti-scam-trainer/internal/attempt"
	"github.com/marlendd/anti-scam-trainer/internal/auth"
	"github.com/marlendd/anti-scam-trainer/internal/engine"
	"github.com/marlendd/anti-scam-trainer/internal/scenario"
)

type answerSubmitterStub struct {
	result        attempt.SubmitAnswerResult
	err           error
	gotInput      attempt.SubmitAnswerInput
	callCount     int
	attemptResult attempt.Attempt
	attemptErr    error
	gotOperation  string
	gotUserID     string
	gotScenarioID scenario.ScenarioID
	stateResult   attempt.State
	stateErr      error
	gotAttemptID  attempt.AttemptID
}

func (stub *answerSubmitterStub) GetState(
	_ context.Context,
	userID string,
	attemptID attempt.AttemptID,
) (attempt.State, error) {
	stub.gotOperation = "get_state"
	stub.gotUserID = userID
	stub.gotAttemptID = attemptID
	return stub.stateResult, stub.stateErr
}

func (stub *answerSubmitterStub) Start(
	_ context.Context,
	userID string,
	scenarioID scenario.ScenarioID,
) (attempt.Attempt, error) {
	stub.gotOperation = "start"
	stub.gotUserID = userID
	stub.gotScenarioID = scenarioID
	return stub.attemptResult, stub.attemptErr
}

func (stub *answerSubmitterStub) Resume(
	_ context.Context,
	userID string,
	scenarioID scenario.ScenarioID,
) (attempt.Attempt, error) {
	stub.gotOperation = "resume"
	stub.gotUserID = userID
	stub.gotScenarioID = scenarioID
	return stub.attemptResult, stub.attemptErr
}

func (stub *answerSubmitterStub) GetLatestCompleted(
	_ context.Context,
	userID string,
	scenarioID scenario.ScenarioID,
) (attempt.Attempt, error) {
	stub.gotOperation = "get_latest_completed"
	stub.gotUserID = userID
	stub.gotScenarioID = scenarioID
	return stub.attemptResult, stub.attemptErr
}

func (stub *answerSubmitterStub) Restart(
	_ context.Context,
	userID string,
	scenarioID scenario.ScenarioID,
) (attempt.Attempt, error) {
	stub.gotOperation = "restart"
	stub.gotUserID = userID
	stub.gotScenarioID = scenarioID
	return stub.attemptResult, stub.attemptErr
}

func (stub *answerSubmitterStub) SubmitAnswer(
	_ context.Context,
	input attempt.SubmitAnswerInput,
) (attempt.SubmitAnswerResult, error) {
	stub.callCount++
	stub.gotInput = input
	return stub.result, stub.err
}

func TestHandler_SubmitAnswer(t *testing.T) {
	t.Parallel()

	nextNodeID := scenario.NodeID("node-next")
	expected := attempt.SubmitAnswerResult{
		AttemptID:   "attempt-1",
		NodeID:      "node-1",
		ChoiceID:    "choice-1",
		Consequence: "safe consequence",
		NextNodeID:  &nextNodeID,
		Completed:   false,
	}
	service := &answerSubmitterStub{result: expected}
	handler := attempt.NewHandler(service, discardLogger())

	request := newSubmitAnswerRequest(t, `{
		"node_id":"node-1",
		"choice_id":"choice-1",
		"idempotency_key":"key-1"
	}`, true)
	response := httptest.NewRecorder()

	handler.SubmitAnswer(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "application/json", response.Header().Get("Content-Type"))

	var got attempt.SubmitAnswerResult
	require.NoError(t, json.NewDecoder(response.Body).Decode(&got))
	require.Equal(t, expected, got)
	require.Equal(t, 1, service.callCount)
	require.Equal(t, attempt.SubmitAnswerInput{
		UserID:         "user-1",
		AttemptID:      "attempt-1",
		NodeID:         "node-1",
		ChoiceID:       "choice-1",
		IdempotencyKey: "key-1",
	}, service.gotInput)
}

func TestHandler_SubmitAnswerRejectsInvalidRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		body          string
		authenticated bool
		path          string
		wantStatus    int
		wantError     string
	}{
		{
			name:       "missing authenticated user",
			body:       `{}`,
			path:       "/api/v1/attempts/attempt-1/answers",
			wantStatus: http.StatusUnauthorized,
			wantError:  "unauthorized",
		},
		{
			name:          "missing attempt id",
			body:          `{}`,
			authenticated: true,
			path:          "/api/v1/attempts//answers",
			wantStatus:    http.StatusBadRequest,
			wantError:     "attempt_id is required",
		},
		{
			name:          "malformed JSON",
			body:          `{`,
			authenticated: true,
			path:          "/api/v1/attempts/attempt-1/answers",
			wantStatus:    http.StatusBadRequest,
			wantError:     "invalid request body",
		},
		{
			name:          "unknown JSON field",
			body:          `{"node_id":"node-1","choice_id":"choice-1","idempotency_key":"key-1","extra":true}`,
			authenticated: true,
			path:          "/api/v1/attempts/attempt-1/answers",
			wantStatus:    http.StatusBadRequest,
			wantError:     "invalid request body",
		},
		{
			name:          "multiple JSON values",
			body:          `{} {}`,
			authenticated: true,
			path:          "/api/v1/attempts/attempt-1/answers",
			wantStatus:    http.StatusBadRequest,
			wantError:     "invalid request body",
		},
		{
			name:          "missing node id",
			body:          `{"choice_id":"choice-1","idempotency_key":"key-1"}`,
			authenticated: true,
			path:          "/api/v1/attempts/attempt-1/answers",
			wantStatus:    http.StatusBadRequest,
			wantError:     "node_id is required",
		},
		{
			name:          "missing choice id",
			body:          `{"node_id":"node-1","idempotency_key":"key-1"}`,
			authenticated: true,
			path:          "/api/v1/attempts/attempt-1/answers",
			wantStatus:    http.StatusBadRequest,
			wantError:     "choice_id is required",
		},
		{
			name:          "missing idempotency key",
			body:          `{"node_id":"node-1","choice_id":"choice-1"}`,
			authenticated: true,
			path:          "/api/v1/attempts/attempt-1/answers",
			wantStatus:    http.StatusBadRequest,
			wantError:     "idempotency_key is required",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			service := &answerSubmitterStub{}
			handler := attempt.NewHandler(service, discardLogger())
			request := httptest.NewRequest(http.MethodPost, test.path, bytes.NewBufferString(test.body))
			request.SetPathValue("attemptID", pathAttemptID(test.path))
			if test.authenticated {
				request = request.WithContext(auth.ContextWithUserID(request.Context(), "user-1"))
			}
			response := httptest.NewRecorder()

			handler.SubmitAnswer(response, request)

			require.Equal(t, test.wantStatus, response.Code)
			require.Equal(t, test.wantError, decodeError(t, response))
			require.Zero(t, service.callCount)
		})
	}
}

func TestHandler_SubmitAnswerMapsServiceErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantError  string
	}{
		{
			name:       "attempt not found",
			err:        attempt.ErrAttemptNotFound,
			wantStatus: http.StatusNotFound,
			wantError:  "attempt or scenario not found",
		},
		{
			name:       "scenario not found through wrapped error",
			err:        errors.Join(errors.New("repository failure"), scenario.ErrScenarioNotFound),
			wantStatus: http.StatusNotFound,
			wantError:  "attempt or scenario not found",
		},
		{
			name:       "unknown choice",
			err:        engine.ErrChoiceNotFound,
			wantStatus: http.StatusBadRequest,
			wantError:  "choice not found in current node",
		},
		{
			name:       "idempotency conflict",
			err:        attempt.ErrIdempotencyConflict,
			wantStatus: http.StatusConflict,
			wantError:  "answer conflicts with attempt state",
		},
		{
			name:       "node already answered",
			err:        attempt.ErrNodeAlreadyAnswered,
			wantStatus: http.StatusConflict,
			wantError:  "answer conflicts with attempt state",
		},
		{
			name:       "attempt completed",
			err:        attempt.ErrAttemptNotInProgress,
			wantStatus: http.StatusConflict,
			wantError:  "answer conflicts with attempt state",
		},
		{
			name:       "current node mismatch",
			err:        attempt.ErrAttemptNodeMismatch,
			wantStatus: http.StatusConflict,
			wantError:  "answer conflicts with attempt state",
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

			service := &answerSubmitterStub{err: test.err}
			handler := attempt.NewHandler(service, discardLogger())
			request := newSubmitAnswerRequest(t, `{
				"node_id":"node-1",
				"choice_id":"choice-1",
				"idempotency_key":"key-1"
			}`, true)
			response := httptest.NewRecorder()

			handler.SubmitAnswer(response, request)

			require.Equal(t, test.wantStatus, response.Code)
			require.Equal(t, test.wantError, decodeError(t, response))
			require.Equal(t, 1, service.callCount)
		})
	}
}

func newSubmitAnswerRequest(t *testing.T, body string, authenticated bool) *http.Request {
	t.Helper()

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/attempts/attempt-1/answers",
		bytes.NewBufferString(body),
	)
	request.SetPathValue("attemptID", "attempt-1")
	if authenticated {
		request = request.WithContext(auth.ContextWithUserID(request.Context(), "user-1"))
	}

	return request
}

func decodeError(t *testing.T, response *httptest.ResponseRecorder) string {
	t.Helper()

	var payload map[string]string
	require.NoError(t, json.NewDecoder(response.Body).Decode(&payload))
	return payload["error"]
}

func pathAttemptID(path string) string {
	if path == "/api/v1/attempts//answers" {
		return ""
	}
	return "attempt-1"
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
