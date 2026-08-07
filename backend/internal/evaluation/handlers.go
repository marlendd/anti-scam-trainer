package evaluation

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/marlendd/anti-scam-trainer/internal/auth"
)

type Handler struct {
	service *Service
	log     *slog.Logger
}

func NewHandler(service *Service, log *slog.Logger) *Handler {
	return &Handler{service: service, log: log}
}

func (h *Handler) GetStatsOfAttempt(w http.ResponseWriter, r *http.Request) {
	attemptID := r.PathValue("id")

	res, err := h.service.GetAttemptResults(r.Context(), attemptID)
	if err != nil {
		h.log.Error("failed to get attempt stats", "attempt_id", attemptID, "error", err)

		h.respondError(w, http.StatusInternalServerError, "database error")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	h.respondJSON(w, http.StatusOK, map[string]int{"score": res})
}

func (h *Handler) GetGlobalStatsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stats, err := h.service.GetGlobalStats(ctx)
	if err != nil {
		h.log.Error("failed to get stats", "error", err)
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.respondJSON(w, http.StatusOK, stats)
}

func (h *Handler) GetCategoryStats(w http.ResponseWriter, r *http.Request) {
	data, err := h.service.GetCategoryDashboard(r.Context())
	if err != nil {
		h.log.Error("failed to get dashboard stats", "error", err)
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.respondJSON(w, http.StatusOK, data)
}

func (h *Handler) GetMyPuzzleProgress(w http.ResponseWriter, r *http.Request) {
	// Достаем ID пользователя из контекста
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	progress, err := h.service.GetUserPuzzleProgress(r.Context(), userID)
	if err != nil {
		h.log.Error("failed to get puzzle progress", "user_id", userID, "error", err)
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.respondJSON(w, http.StatusOK, progress)
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}
