package scenario

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/marlendd/anti-scam-trainer/internal/auth"
)

type catalogLister interface {
	ListActiveByRole(
		ctx context.Context,
		userID string,
		role Role,
	) ([]CatalogItem, error)
}

type CatalogHandler struct {
	service catalogLister
	log     *slog.Logger
}

func NewCatalogHandler(service catalogLister, log *slog.Logger) *CatalogHandler {
	return &CatalogHandler{service: service, log: log}
}

type catalogItemResponse struct {
	ID          ScenarioID        `json:"id"`
	LogicalID   LogicalScenarioID `json:"logical_id"`
	Version     int               `json:"version"`
	Role        Role              `json:"role"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Product     Product           `json:"product"`
	Status      ProgressStatus    `json:"status"`
}

type catalogResponse struct {
	Items []catalogItemResponse `json:"items"`
}

func (handler *CatalogHandler) List(w http.ResponseWriter, request *http.Request) {
	userID, ok := auth.UserIDFromContext(request.Context())
	if !ok || userID == "" {
		writeCatalogError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	role := Role(request.URL.Query().Get("role"))
	items, err := handler.service.ListActiveByRole(request.Context(), userID, role)
	if err != nil {
		if errors.Is(err, ErrInvalidRole) {
			writeCatalogError(w, http.StatusBadRequest, "role must be buyer or seller")
			return
		}

		handler.log.Error("list scenario catalog failed", "error", err)
		writeCatalogError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responseItems := make([]catalogItemResponse, 0, len(items))
	for _, item := range items {
		responseItems = append(responseItems, catalogItemResponse(item))
	}

	writeCatalogJSON(w, http.StatusOK, catalogResponse{Items: responseItems})
}

func writeCatalogJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeCatalogError(w http.ResponseWriter, status int, message string) {
	writeCatalogJSON(w, status, map[string]string{"error": message})
}
