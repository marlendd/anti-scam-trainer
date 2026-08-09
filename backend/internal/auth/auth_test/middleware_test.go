package auth_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/marlendd/anti-scam-trainer/internal/auth"
	"github.com/stretchr/testify/require"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestRequireAuth(t *testing.T) {
	t.Parallel()

	t.Run("returns 401 when no cookie present", func(t *testing.T) {
		t.Parallel()

		service := newTestService(&userRepositoryStub{}, &sessionRepositoryStub{})
		middleware := auth.RequireAuth(service, testLogger())

		nextCalled := false
		handler := middleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			nextCalled = true
		}))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
		require.False(t, nextCalled, "next handler must not be called without a session")
	})

	t.Run("returns 401 when session is not found", func(t *testing.T) {
		t.Parallel()

		sessions := &sessionRepositoryStub{
			getByIDFn: func(context.Context, string) (auth.Session, error) {
				return auth.Session{}, auth.ErrSessionNotFound
			},
		}
		service := newTestService(&userRepositoryStub{}, sessions)
		middleware := auth.RequireAuth(service, testLogger())

		nextCalled := false
		handler := middleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			nextCalled = true
		}))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
		req.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: "invalid-session"})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
		require.False(t, nextCalled)
	})

	t.Run("returns 500 on unexpected repository error", func(t *testing.T) {
		t.Parallel()

		sessions := &sessionRepositoryStub{
			getByIDFn: func(context.Context, string) (auth.Session, error) {
				return auth.Session{}, errUnexpected{}
			},
		}
		service := newTestService(&userRepositoryStub{}, sessions)
		middleware := auth.RequireAuth(service, testLogger())

		handler := middleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			t.Fatal("next handler must not be called on internal error")
		}))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
		req.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: "some-session"})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("passes userID to context and calls next on valid session", func(t *testing.T) {
		t.Parallel()

		sessions := &sessionRepositoryStub{
			getByIDFn: func(context.Context, string) (auth.Session, error) {
				return auth.Session{
					ID:        "session-1",
					UserID:    "user-1",
					ExpiresAt: time.Now().Add(time.Hour),
				}, nil
			},
		}
		users := &userRepositoryStub{
			getByIDFn: func(context.Context, string) (auth.User, error) {
				return auth.User{ID: "user-1", Email: "test@example.com"}, nil
			},
		}
		service := newTestService(users, sessions)
		middleware := auth.RequireAuth(service, testLogger())

		var capturedUserID string
		var ok bool
		handler := middleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			capturedUserID, ok = auth.UserIDFromContext(r.Context())
		}))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
		req.AddCookie(&http.Cookie{Name: auth.SessionCookieName, Value: "session-1"})
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		require.True(t, ok, "userID must be present in context")
		require.Equal(t, "user-1", capturedUserID)
	})
}

type errUnexpected struct{}

func (errUnexpected) Error() string { return "unexpected repository failure" }
