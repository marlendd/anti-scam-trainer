package attempt

import (
	"context"
	"fmt"

	"github.com/marlendd/anti-scam-trainer/internal/scenario"
)

type AttemptRepository interface {
	WithinTransaction(
		ctx context.Context,
		fn func(AttemptRepository) error,
	) error
	Create(
		ctx context.Context,
		userID string,
		scenarioID scenario.ScenarioID,
		startNodeID scenario.NodeID,
	) (Attempt, error)
	GetByID(
		ctx context.Context,
		attemptID AttemptID,
		userID string,
	) (Attempt, error)
	GetActive(
		ctx context.Context,
		userID string,
		scenarioID scenario.ScenarioID,
	) (Attempt, error)
	Abort(
		ctx context.Context,
		attemptID AttemptID,
		userID string,
	) error
}

type ScenarioProvider interface {
	GetActiveByID(
		ctx context.Context,
		scenario scenario.ScenarioID,
	) (scenario.Scenario, error)
}

type Service struct {
	attempts  AttemptRepository
	scenarios ScenarioProvider
}

func NewService(
	attempts AttemptRepository,
	scenario ScenarioProvider,
) *Service {
	return &Service{
		attempts:  attempts,
		scenarios: scenario,
	}
}

func (s *Service) Start(
	ctx context.Context,
	userID string,
	scenarioID scenario.ScenarioID,
) (Attempt, error) {
	currentScenario, err := s.scenarios.GetActiveByID(ctx, scenarioID)
	if err != nil {
		return Attempt{}, fmt.Errorf("get active scenario by id: %w", err)
	}

	if err := scenario.Validate(currentScenario); err != nil {
		return Attempt{}, fmt.Errorf("validate scenario: %w", err)
	}

	newAttempt, err := s.attempts.Create(
		ctx,
		userID,
		scenarioID,
		currentScenario.StartNodeID,
	)
	if err != nil {
		return Attempt{}, fmt.Errorf("create new attempt: %w", err)
	}

	return newAttempt, nil
}

func (s *Service) Resume(
	ctx context.Context,
	userID string,
	scenarioID scenario.ScenarioID,
) (Attempt, error) {
	currentAttempt, err := s.attempts.GetActive(ctx, userID, scenarioID)
	if err != nil {
		return Attempt{}, fmt.Errorf("get active attempt: %w", err)
	}

	return currentAttempt, nil
}

func (s *Service) Restart(
	ctx context.Context,
	userID string,
	scenarioID scenario.ScenarioID,
) (Attempt, error) {
	currentScenario, err := s.scenarios.GetActiveByID(ctx, scenarioID)
	if err != nil {
		return Attempt{}, fmt.Errorf("get active scenario by id: %w", err)
	}

	if err := scenario.Validate(currentScenario); err != nil {
		return Attempt{}, fmt.Errorf("validate scenario: %w", err)
	}

	var newAttempt Attempt

	if err := s.attempts.WithinTransaction(
		ctx,
		func(ar AttemptRepository) error {
			currentAttempt, err := ar.GetActive(
				ctx,
				userID,
				scenarioID,
			)
			if err != nil {
				return fmt.Errorf("get active attempt: %w", err)
			}

			if err := ar.Abort(ctx, currentAttempt.ID, userID); err != nil {
				return fmt.Errorf("abort attempt: %w", err)
			}

			createdAttempt, err := ar.Create(
				ctx,
				userID,
				scenarioID,
				currentScenario.StartNodeID,
			)
			if err != nil {
				return fmt.Errorf("create attempt: %w", err)
			}

			newAttempt = createdAttempt

			return nil
		},
	); err != nil {
		return Attempt{}, fmt.Errorf("restart attempt: %w", err)
	}

	return newAttempt, nil
}
