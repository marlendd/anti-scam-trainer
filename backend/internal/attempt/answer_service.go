package attempt

import (
	"context"

	"github.com/marlendd/anti-scam-trainer/internal/scenario"
)

type AnswerRepository interface {
	WithinAnswerTransaction(ctx context.Context, fn func(AnswerRepository) error) error

	GetAttemptForUpdate(
		ctx context.Context,
		attemptID AttemptID,
		userID string,
	) (Attempt, error)

	GetAnswerByIdempotencyKey(
		ctx context.Context,
		userID string,
		idempotencyKey IdempotencyKey,
	) (Answer, error)

	GetAnswerByAttemptNode(ctx context.Context,
		attemptID AttemptID,
		userID string,
		nodeID scenario.NodeID,
	) (Answer, error)

	GetScoreTotals(
		ctx context.Context,
		attemptID AttemptID,
	) (ScoreTotals, error)

	CreateAnswer(
		ctx context.Context,
		answer Answer,
	) (Answer, error)

	AdvanceAttempt(
		ctx context.Context,
		attemptID AttemptID,
		userID string,
		currentNodeID scenario.NodeID,
		nextNodeID scenario.NodeID,
	) error

	CompleteAttempt(
		ctx context.Context,
		attemptID AttemptID,
		userID string,
		currentNodeID scenario.NodeID,
		endingID scenario.EndingID,
		score int,
	) error
}
