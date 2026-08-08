package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
)

type contextKey string

const userIDContextKey contextKey = "userID"

const SessionCookieName = "session_id"

func RequireAuth(service *Service, log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(SessionCookieName)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := service.GetUserBySession(r.Context(), cookie.Value)
			if err != nil {
				if errors.Is(err, ErrSessionNotFound) || errors.Is(err, ErrUserNotFound) {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
				}
				log.Error("failed to check session", "error", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			ctx := ContextWithUserID(r.Context(), user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDContextKey).(string)
	return id, ok
}
