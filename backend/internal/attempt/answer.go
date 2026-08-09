package attempt

import (
	"time"

	"github.com/marlendd/anti-scam-trainer/internal/scenario"
)

type AnswerID string
type IdempotencyKey string

type ScoreTotals struct {
	WeightedScoreSum int
	WeightSum        int
}

type Answer struct {
	ID             AnswerID
	AttemptID      AttemptID
	NodeID         scenario.NodeID
	ChoiceID       scenario.ChoiceID
	IdempotencyKey IdempotencyKey
	Weight         scenario.Weight
	ChoiceScore    scenario.ChoiceScore
	RiskCategories []scenario.RiskCategory
	Consequence    string
	Explanation    string
	Response       SubmitAnswerResult
	CreatedAt      time.Time
}

type SubmitAnswerInput struct {
	UserID         string            `json:"-"`
	AttemptID      AttemptID         `json:"attempt_id"`
	NodeID         scenario.NodeID   `json:"node_id"`
	ChoiceID       scenario.ChoiceID `json:"choice_id"`
	IdempotencyKey IdempotencyKey    `json:"idempotency_key"`
}

type SubmitAnswerResult struct {
	AttemptID        AttemptID            `json:"attempt_id"`
	NodeID           scenario.NodeID      `json:"node_id"`
	ChoiceID         scenario.ChoiceID    `json:"choice_id"`
	Consequence      string               `json:"consequence"`
	NextNodeID       *scenario.NodeID     `json:"next_node_id,omitempty"`
	EndingID         *scenario.EndingID   `json:"ending_id,omitempty"`
	Completed        bool                 `json:"completed"`
	Score            *int                 `json:"score,omitempty"`
	RewardFragmentID *scenario.FragmentID `json:"reward_fragment_id,omitempty"`
}
