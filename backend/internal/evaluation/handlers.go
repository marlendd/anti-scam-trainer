package evaluation

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

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

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	res, err := h.service.GetAttemptResults(r.Context(), userID, attemptID)
	if err != nil {
		h.log.Error("failed to get attempt stats", "attempt_id", attemptID, "error", err)
		h.respondError(w, http.StatusInternalServerError, "database error")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]int{"score": res})
}

func (h *Handler) GetMyRoleStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	stats, err := h.service.GetUserRoleStats(r.Context(), userID)
	if err != nil {
		h.log.Error("failed to get personal role stats", "error", err)
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.respondJSON(w, http.StatusOK, stats)
}

func (h *Handler) GetMyCategoryDashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	data, err := h.service.GetUserCategoryDashboard(r.Context(), userID)
	if err != nil {
		h.log.Error("failed to get personal dashboard stats", "error", err)
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.respondJSON(w, http.StatusOK, data)
}

func (h *Handler) GetMyPuzzleProgress(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	data, err := h.service.GetLeaderboard(r.Context(), limit, offset)
	if err != nil {
		h.log.Error("failed to get leaderboard", "error", err)
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.respondJSON(w, http.StatusOK, data)
}

func (h *Handler) GetMyRankHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	data, err := h.service.GetMyRankHistory(r.Context(), userID)
	if err != nil {
		h.log.Error("failed to get rank history", "user_id", userID, "error", err)
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.respondJSON(w, http.StatusOK, data)
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}
