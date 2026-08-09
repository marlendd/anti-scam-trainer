package attempt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/marlendd/anti-scam-trainer/internal/auth"
	"github.com/marlendd/anti-scam-trainer/internal/engine"
	"github.com/marlendd/anti-scam-trainer/internal/scenario"
)

type handlerService interface {
	GetState(
		ctx context.Context,
		userID string,
		attemptID AttemptID,
	) (State, error)

	Start(
		ctx context.Context,
		userID string,
		scenarioID scenario.ScenarioID,
	) (Attempt, error)

	Resume(
		ctx context.Context,
		userID string,
		scenarioID scenario.ScenarioID,
	) (Attempt, error)

	Restart(
		ctx context.Context,
		userID string,
		scenarioID scenario.ScenarioID,
	) (Attempt, error)

	SubmitAnswer(
		ctx context.Context,
		input SubmitAnswerInput,
	) (SubmitAnswerResult, error)
}

type choiceOptionResponse struct {
	ID   scenario.ChoiceID `json:"id"`
	Text string            `json:"text"`
}

type currentNodeResponse struct {
	ID       scenario.NodeID        `json:"id"`
	Author   scenario.AuthorID      `json:"author"`
	Text     string                 `json:"text"`
	Messages []messageResponse      `json:"messages"`
	Choices  []choiceOptionResponse `json:"choices"`
}

type messageResponse struct {
	Author scenario.AuthorID `json:"author"`
	Text   string            `json:"text"`
}

type historyNodeResponse struct {
	ID       scenario.NodeID   `json:"id"`
	Author   scenario.AuthorID `json:"author"`
	Text     string            `json:"text"`
	Messages []messageResponse `json:"messages"`
}

type historyItemResponse struct {
	Node           historyNodeResponse  `json:"node"`
	SelectedChoice choiceOptionResponse `json:"selected_choice"`
	Consequence    string               `json:"consequence"`
	AnsweredAt     time.Time            `json:"answered_at"`
}

type attemptStateResponse struct {
	attemptResponse
	CurrentNode *currentNodeResponse  `json:"current_node,omitempty"`
	History     []historyItemResponse `json:"history"`
}

type Handler struct {
	service handlerService
	log     *slog.Logger
}

func NewHandler(service handlerService, log *slog.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

type attemptResponse struct {
	ID            AttemptID           `json:"id"`
	ScenarioID    scenario.ScenarioID `json:"scenario_id"`
	Status        Status              `json:"status"`
	CurrentNodeID *scenario.NodeID    `json:"current_node_id,omitempty"`
	EndingID      *scenario.EndingID  `json:"ending_id,omitempty"`
	Score         *int                `json:"score,omitempty"`
	StartedAt     time.Time           `json:"started_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	CompletedAt   *time.Time          `json:"completed_at,omitempty"`
}

func newAttemptResponse(currentAttempt Attempt) attemptResponse {
	return attemptResponse{
		ID:            currentAttempt.ID,
		ScenarioID:    currentAttempt.ScenarioID,
		Status:        currentAttempt.Status,
		CurrentNodeID: currentAttempt.CurrentNodeID,
		EndingID:      currentAttempt.EndingID,
		Score:         currentAttempt.Score,
		StartedAt:     currentAttempt.StartedAt,
		UpdatedAt:     currentAttempt.UpdatedAt,
		CompletedAt:   currentAttempt.CompletedAt,
	}
}

func newAttemptStateResponse(state State) attemptStateResponse {
	response := attemptStateResponse{
		attemptResponse: newAttemptResponse(state.Attempt),
		History:         make([]historyItemResponse, 0, len(state.History)),
	}

	for _, item := range state.History {
		response.History = append(response.History, historyItemResponse{
			Node: historyNodeResponse{
				ID:       item.Node.ID,
				Author:   item.Node.Author,
				Text:     item.Node.Text,
				Messages: newMessageResponses(item.Node.Messages),
			},
			SelectedChoice: choiceOptionResponse(item.SelectedChoice),
			Consequence:    item.Consequence,
			AnsweredAt:     item.AnsweredAt,
		})
	}

	if state.CurrentNode == nil {
		return response
	}

	choices := make([]choiceOptionResponse, 0, len(state.CurrentNode.Choices))
	for _, choice := range state.CurrentNode.Choices {
		choices = append(choices, choiceOptionResponse(choice))
	}

	response.CurrentNode = &currentNodeResponse{
		ID:       state.CurrentNode.ID,
		Author:   state.CurrentNode.Author,
		Text:     state.CurrentNode.Text,
		Messages: newMessageResponses(state.CurrentNode.Messages),
		Choices:  choices,
	}

	return response
}

func newMessageResponses(messages []scenario.Message) []messageResponse {
	response := make([]messageResponse, 0, len(messages))
	for _, message := range messages {
		response = append(response, messageResponse{
			Author: message.Author,
			Text:   message.Text,
		})
	}

	return response
}

func (h *Handler) GetState(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	attemptID := AttemptID(r.PathValue("attemptID"))
	if attemptID == "" {
		respondError(w, http.StatusBadRequest, "attempt_id is required")
		return
	}

	state, err := h.service.GetState(r.Context(), userID, attemptID)
	if err != nil {
		h.respondAttemptError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, newAttemptStateResponse(state))
}

func (h *Handler) Start(w http.ResponseWriter, r *http.Request) {
	h.handleAttemptOperation(w, r, http.StatusCreated, h.service.Start)
}

func (h *Handler) Resume(w http.ResponseWriter, r *http.Request) {
	h.handleAttemptOperation(w, r, http.StatusOK, h.service.Resume)
}

func (h *Handler) Restart(w http.ResponseWriter, r *http.Request) {
	h.handleAttemptOperation(w, r, http.StatusCreated, h.service.Restart)
}

func (h *Handler) handleAttemptOperation(
	w http.ResponseWriter,
	r *http.Request,
	successStatus int,
	operation func(context.Context, string, scenario.ScenarioID) (Attempt, error),
) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	scenarioID := scenario.ScenarioID(r.PathValue("scenarioID"))
	if scenarioID == "" {
		respondError(w, http.StatusBadRequest, "scenario_id is required")
		return
	}

	currentAttempt, err := operation(r.Context(), userID, scenarioID)
	if err != nil {
		h.respondAttemptError(w, err)
		return
	}

	respondJSON(w, successStatus, newAttemptResponse(currentAttempt))
}

func (h *Handler) respondAttemptError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, scenario.ErrScenarioNotFound):
		respondError(w, http.StatusNotFound, "scenario not found")
	case errors.Is(err, ErrActiveAttemptNotFound),
		errors.Is(err, ErrAttemptNotFound):
		respondError(w, http.StatusNotFound, "attempt not found")
	case errors.Is(err, scenario.ErrScenarioInactive):
		respondError(w, http.StatusConflict, "scenario is inactive")
	case errors.Is(err, ErrActiveAttemptExists):
		respondError(w, http.StatusConflict, "active attempt already exists")
	default:
		h.log.Error("attempt operation failed", "error", err)
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}

type submitAnswerRequest struct {
	NodeID         scenario.NodeID   `json:"node_id"`
	ChoiceID       scenario.ChoiceID `json:"choice_id"`
	IdempotencyKey IdempotencyKey    `json:"idempotency_key"`
}

func (h *Handler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	attemptID := AttemptID(r.PathValue("attemptID"))
	if attemptID == "" {
		respondError(w, http.StatusBadRequest, "attempt_id is required")
		return
	}

	var request submitAnswerRequest
	if err := decodeJSON(r, &request); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	switch {
	case request.NodeID == "":
		respondError(w, http.StatusBadRequest, "node_id is required")
		return
	case request.ChoiceID == "":
		respondError(w, http.StatusBadRequest, "choice_id is required")
		return
	case request.IdempotencyKey == "":
		respondError(w, http.StatusBadRequest, "idempotency_key is required")
		return
	}

	result, err := h.service.SubmitAnswer(r.Context(), SubmitAnswerInput{
		UserID:         userID,
		AttemptID:      attemptID,
		NodeID:         request.NodeID,
		ChoiceID:       request.ChoiceID,
		IdempotencyKey: request.IdempotencyKey,
	})
	if err != nil {
		h.respondSubmitAnswerError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, result)
}

func (h *Handler) respondSubmitAnswerError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrAttemptNotFound),
		errors.Is(err, scenario.ErrScenarioNotFound):
		respondError(w, http.StatusNotFound, "attempt or scenario not found")
	case errors.Is(err, engine.ErrChoiceNotFound):
		respondError(w, http.StatusBadRequest, "choice not found in current node")
	case errors.Is(err, ErrIdempotencyConflict),
		errors.Is(err, ErrNodeAlreadyAnswered),
		errors.Is(err, ErrAttemptNotInProgress),
		errors.Is(err, ErrAttemptNodeMismatch):
		respondError(w, http.StatusConflict, "answer conflicts with attempt state")
	default:
		h.log.Error("submit answer failed", "error", err)
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}

func decodeJSON(r *http.Request, destination any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("decode JSON: %w", err)
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON value")
	}

	return nil
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
