package scenario_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/marlendd/anti-scam-trainer/internal/auth"
	"github.com/marlendd/anti-scam-trainer/internal/scenario"
)

type catalogServiceStub struct {
	items     []scenario.CatalogItem
	err       error
	gotUserID string
	gotRole   scenario.Role
	called    bool
}

func (stub *catalogServiceStub) ListActiveByRole(
	_ context.Context,
	userID string,
	role scenario.Role,
) ([]scenario.CatalogItem, error) {
	stub.called = true
	stub.gotUserID = userID
	stub.gotRole = role
	return stub.items, stub.err
}

func TestCatalogHandler_List(t *testing.T) {
	t.Parallel()

	bestScore := 85
	service := &catalogServiceStub{items: []scenario.CatalogItem{
		{
			ID:          "scenario-1",
			LogicalID:   "logical-1",
			Version:     2,
			Role:        scenario.RoleBuyer,
			Title:       "Safe delivery",
			Description: "Recognize a fake delivery link",
			Status:      scenario.ProgressCompleted,
			BestScore:   &bestScore,
		},
	}}
	handler := scenario.NewCatalogHandler(service, catalogDiscardLogger())
	request := httptest.NewRequest(http.MethodGet, "/api/v1/scenarios?role=buyer", nil)
	request = request.WithContext(auth.ContextWithUserID(request.Context(), "user-1"))
	response := httptest.NewRecorder()

	handler.List(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "application/json", response.Header().Get("Content-Type"))
	require.True(t, service.called)
	require.Equal(t, "user-1", service.gotUserID)
	require.Equal(t, scenario.RoleBuyer, service.gotRole)

	var payload struct {
		Items []struct {
			ID        string `json:"id"`
			Status    string `json:"status"`
			BestScore *int   `json:"best_score"`
		} `json:"items"`
	}
	require.NoError(t, json.NewDecoder(response.Body).Decode(&payload))
	require.Len(t, payload.Items, 1)
	require.Equal(t, "scenario-1", payload.Items[0].ID)
	require.Equal(t, "completed", payload.Items[0].Status)
	require.NotNil(t, payload.Items[0].BestScore)
	require.Equal(t, 85, *payload.Items[0].BestScore)
}

func TestCatalogHandler_ListRejectsUnauthorizedRequest(t *testing.T) {
	t.Parallel()

	service := &catalogServiceStub{}
	handler := scenario.NewCatalogHandler(service, catalogDiscardLogger())
	response := httptest.NewRecorder()

	handler.List(
		response,
		httptest.NewRequest(http.MethodGet, "/api/v1/scenarios?role=buyer", nil),
	)

	require.Equal(t, http.StatusUnauthorized, response.Code)
	require.False(t, service.called)
	require.Equal(t, "unauthorized", catalogErrorMessage(t, response))
}

func TestCatalogHandler_ListMapsErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantError  string
	}{
		{
			name:       "invalid role",
			err:        scenario.ErrInvalidRole,
			wantStatus: http.StatusBadRequest,
			wantError:  "role must be buyer or seller",
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

			service := &catalogServiceStub{err: test.err}
			handler := scenario.NewCatalogHandler(service, catalogDiscardLogger())
			request := httptest.NewRequest(http.MethodGet, "/api/v1/scenarios?role=invalid", nil)
			request = request.WithContext(auth.ContextWithUserID(request.Context(), "user-1"))
			response := httptest.NewRecorder()

			handler.List(response, request)

			require.Equal(t, test.wantStatus, response.Code)
			require.Equal(t, test.wantError, catalogErrorMessage(t, response))
		})
	}
}

func catalogErrorMessage(t *testing.T, response *httptest.ResponseRecorder) string {
	t.Helper()

	var payload map[string]string
	require.NoError(t, json.NewDecoder(response.Body).Decode(&payload))
	return payload["error"]
}

func catalogDiscardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
