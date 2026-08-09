package feedback

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

func (h *Handler) GetAttemptFeedback(w http.ResponseWriter, r *http.Request) {
	attemptID := r.PathValue("id")

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	res, err := h.service.Generate(r.Context(), userID, attemptID)
	if err != nil {
		if err.Error() == "not_found" {
			h.respondError(w, http.StatusNotFound, "попытка не найдена или нет ответов")
			return
		}
		h.log.Error("failed to generate feedback", "attempt_id", attemptID, "error", err)
		h.respondError(w, http.StatusInternalServerError, "ошибка генерации фидбека")
		return
	}

	h.respondJSON(w, http.StatusOK, res)
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}
