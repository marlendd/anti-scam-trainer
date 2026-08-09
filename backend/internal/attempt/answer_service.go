package attempt

import (
	"context"
	"errors"
	"fmt"

	"github.com/marlendd/anti-scam-trainer/internal/engine"
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

	ListAnswersByAttempt(
		ctx context.Context,
		attemptID AttemptID,
		userID string,
	) ([]Answer, error)

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

	GrantFragment(
		ctx context.Context,
		userID string,
		scenarioID scenario.ScenarioID,
		fragmentID scenario.FragmentID,
	) (bool, error)
}

type idempotencyAnswerGetter interface {
	GetAnswerByIdempotencyKey(
		ctx context.Context,
		userID string,
		idempotencyKey IdempotencyKey,
	) (Answer, error)
}

func findIdempotentResult(
	ctx context.Context,
	repository idempotencyAnswerGetter,
	input SubmitAnswerInput,
) (SubmitAnswerResult, bool, error) {
	answer, err := repository.GetAnswerByIdempotencyKey(
		ctx,
		input.UserID,
		input.IdempotencyKey,
	)
	if err != nil {
		if errors.Is(err, ErrAnswerNotFound) {
			return SubmitAnswerResult{}, false, nil
		}

		return SubmitAnswerResult{}, false, fmt.Errorf(
			"get answer by idempotency key: %w",
			err,
		)
	}

	if answer.AttemptID != input.AttemptID ||
		answer.NodeID != input.NodeID ||
		answer.ChoiceID != input.ChoiceID {
		return SubmitAnswerResult{}, false, ErrIdempotencyConflict
	}

	return answer.Response, true, nil
}

func (s *Service) SubmitAnswer(ctx context.Context, input SubmitAnswerInput) (SubmitAnswerResult, error) {
	var result SubmitAnswerResult

	if err := s.answers.WithinAnswerTransaction(ctx, func(ar AnswerRepository) error {
		idempotentResult, found, err := findIdempotentResult(
			ctx,
			ar,
			input,
		)
		if err != nil {
			return fmt.Errorf("check idempotent result: %w", err)
		}

		if found {
			result = idempotentResult
			return nil
		}

		currentAttempt, err := ar.GetAttemptForUpdate(
			ctx,
			input.AttemptID,
			input.UserID,
		)
		if err != nil {
			return fmt.Errorf("get attempt for update: %w", err)
		}

		idempotentResult, found, err = findIdempotentResult(
			ctx,
			ar,
			input,
		)
		if err != nil {
			return fmt.Errorf("recheck idempotent result: %w", err)
		}

		if found {
			result = idempotentResult
			return nil
		}

		_, err = ar.GetAnswerByAttemptNode(
			ctx,
			currentAttempt.ID,
			currentAttempt.UserID,
			input.NodeID,
		)
		if err == nil {
			return ErrNodeAlreadyAnswered
		}

		if !errors.Is(err, ErrAnswerNotFound) {
			return fmt.Errorf("get answer by attempt node: %w", err)
		}

		if currentAttempt.Status != StatusInProgress {
			return ErrAttemptNotInProgress
		}

		if currentAttempt.CurrentNodeID == nil {
			return ErrInvalidAttemptState
		}

		if *currentAttempt.CurrentNodeID != input.NodeID {
			return ErrAttemptNodeMismatch
		}

		currentScenario, err := s.scenarios.GetByID(ctx, currentAttempt.ScenarioID)
		if err != nil {
			return fmt.Errorf("get scenario by id: %w", err)
		}

		if err := scenario.Validate(currentScenario); err != nil {
			return fmt.Errorf(
				"validate scenario: %w",
				err,
			)
		}

		transitionResult, err := engine.Transition(engine.Input{
			Scenario:      currentScenario,
			CurrentNodeID: *currentAttempt.CurrentNodeID,
			ChoiceID:      input.ChoiceID,
		})
		if err != nil {
			return fmt.Errorf("calculate scenario transition: %w", err)
		}

		result = SubmitAnswerResult{
			AttemptID:   currentAttempt.ID,
			ChoiceID:    transitionResult.Choice.ID,
			NodeID:      input.NodeID,
			Consequence: transitionResult.Consequence,
			Completed:   transitionResult.Completed,
		}

		if transitionResult.Completed {
			totals, err := ar.GetScoreTotals(ctx, currentAttempt.ID)
			if err != nil {
				return fmt.Errorf("get score totals: %w", err)
			}

			weightedScoreSum := totals.WeightedScoreSum +
				int(transitionResult.Choice.Weight)*
					int(transitionResult.Choice.Score)

			weightSum := totals.WeightSum +
				int(transitionResult.Choice.Weight)

			score, err := engine.CalculateScore(weightedScoreSum, weightSum)
			if err != nil {
				return fmt.Errorf("calculate attempt score: %w", err)
			}

			endingID := transitionResult.EndingID
			result.EndingID = &endingID
			result.Score = &score
		} else {
			nextNodeID := transitionResult.NextNodeID
			result.NextNodeID = &nextNodeID
		}

		if transitionResult.Completed &&
			currentScenario.RewardFragmentID != "" &&
			currentScenario.IsSuccessfulEnding(transitionResult.EndingID) {
			awarded, err := ar.GrantFragment(
				ctx,
				currentAttempt.UserID,
				currentAttempt.ScenarioID,
				currentScenario.RewardFragmentID,
			)
			if err != nil {
				return fmt.Errorf("grant reward fragment: %w", err)
			}
			if awarded {
				fragmentID := currentScenario.RewardFragmentID
				result.RewardFragmentID = &fragmentID
			}
		}

		answer := Answer{
			AttemptID:      currentAttempt.ID,
			NodeID:         input.NodeID,
			ChoiceID:       transitionResult.Choice.ID,
			IdempotencyKey: input.IdempotencyKey,
			Weight:         transitionResult.Choice.Weight,
			ChoiceScore:    transitionResult.Choice.Score,
			RiskCategories: transitionResult.Choice.RiskCategories,
			Consequence:    transitionResult.Choice.Consequence,
			Explanation:    transitionResult.Choice.Explanation,
			Response:       result,
		}

		if _, err := ar.CreateAnswer(ctx, answer); err != nil {
			return fmt.Errorf("create answer: %w", err)
		}

		if transitionResult.Completed {
			if err := ar.CompleteAttempt(
				ctx,
				currentAttempt.ID,
				currentAttempt.UserID,
				input.NodeID,
				transitionResult.EndingID,
				*result.Score,
			); err != nil {
				return fmt.Errorf(
					"complete attempt: %w",
					err,
				)
			}
			return nil
		}

		if err := ar.AdvanceAttempt(
			ctx,
			currentAttempt.ID,
			currentAttempt.UserID,
			input.NodeID,
			transitionResult.NextNodeID,
		); err != nil {
			return fmt.Errorf(
				"advance attempt: %w",
				err,
			)
		}

		return nil
	}); err != nil {
		return SubmitAnswerResult{}, fmt.Errorf("submit answer: %w", err)
	}

	return result, nil
}
