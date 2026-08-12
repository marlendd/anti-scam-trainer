package attempt

import (
	"context"
	"fmt"

	"github.com/marlendd/anti-scam-trainer/internal/evaluation"
	"github.com/marlendd/anti-scam-trainer/internal/scenario"
)

type ScoreEvaluator interface {
	CalculateScoreFromTotals(weightedScoreSum, weightSum int) (int, error)
}

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

	GetLatestCompleted(
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

	GetByID(
		ctx context.Context,
		scenario scenario.ScenarioID,
	) (scenario.Scenario, error)
}

type Service struct {
	attempts  AttemptRepository
	answers   AnswerRepository
	scenarios ScenarioProvider
	evaluator ScoreEvaluator
}

func NewService(
	attempts AttemptRepository,
	answers AnswerRepository,
	scenario ScenarioProvider,
	evaluators ...ScoreEvaluator,
) *Service {
	evaluator := ScoreEvaluator(evaluation.NewEvaluator())
	if len(evaluators) > 0 {
		evaluator = evaluators[0]
	}

	return &Service{
		attempts:  attempts,
		answers:   answers,
		scenarios: scenario,
		evaluator: evaluator,
	}
}

func (s *Service) GetState(
	ctx context.Context,
	userID string,
	attemptID AttemptID,
) (State, error) {
	currentAttempt, err := s.attempts.GetByID(ctx, attemptID, userID)
	if err != nil {
		return State{}, fmt.Errorf("get attempt by id: %w", err)
	}

	currentScenario, err := s.scenarios.GetByID(ctx, currentAttempt.ScenarioID)
	if err != nil {
		return State{}, fmt.Errorf("get scenario by id: %w", err)
	}

	if err := scenario.Validate(currentScenario); err != nil {
		return State{}, fmt.Errorf("validate scenario: %w", err)
	}

	answers, err := s.answers.ListAnswersByAttempt(ctx, attemptID, userID)
	if err != nil {
		return State{}, fmt.Errorf("list attempt answers: %w", err)
	}

	state := State{
		Attempt: currentAttempt,
		Scenario: ScenarioHeader{
			Title:       currentScenario.Title,
			Description: currentScenario.Description,
			Role:        currentScenario.Role,
			Product:     currentScenario.Product,
		},
		History: make([]HistoryItem, 0, len(answers)),
	}

	for _, answer := range answers {
		node, found := findScenarioNode(currentScenario, answer.NodeID)
		if !found {
			return State{}, ErrInvalidAttemptState
		}

		choice, found := findNodeChoice(node, answer.ChoiceID)
		if !found {
			return State{}, ErrInvalidAttemptState
		}

		messages := node.DialogueMessages()
		lastMessage := messages[len(messages)-1]
		state.History = append(state.History, HistoryItem{
			Node: HistoryNode{
				ID:       node.ID,
				Author:   lastMessage.Author,
				Text:     lastMessage.Text,
				Messages: messages,
			},
			SelectedChoice: ChoiceOption{
				ID:   choice.ID,
				Text: choice.Text,
			},
			Consequence: answer.Consequence,
			AnsweredAt:  answer.CreatedAt,
		})
	}

	if currentAttempt.Status != StatusInProgress {
		return state, nil
	}

	if currentAttempt.CurrentNodeID == nil {
		return State{}, ErrInvalidAttemptState
	}

	node, found := findScenarioNode(currentScenario, *currentAttempt.CurrentNodeID)
	if !found {
		return State{}, ErrInvalidAttemptState
	}

	choices := make([]ChoiceOption, 0, len(node.Choices))
	for _, choice := range node.Choices {
		choices = append(choices, ChoiceOption{
			ID:   choice.ID,
			Text: choice.Text,
		})
	}
	messages := node.DialogueMessages()
	lastMessage := messages[len(messages)-1]

	state.CurrentNode = &CurrentNode{
		ID:       node.ID,
		Author:   lastMessage.Author,
		Text:     lastMessage.Text,
		Messages: messages,
		Choices:  choices,
	}

	return state, nil
}

func findScenarioNode(
	currentScenario scenario.Scenario,
	nodeID scenario.NodeID,
) (scenario.Node, bool) {
	for _, node := range currentScenario.Nodes {
		if node.ID == nodeID {
			return node, true
		}
	}

	return scenario.Node{}, false
}

func findNodeChoice(node scenario.Node, choiceID scenario.ChoiceID) (scenario.Choice, bool) {
	for _, choice := range node.Choices {
		if choice.ID == choiceID {
			return choice, true
		}
	}

	return scenario.Choice{}, false
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

func (s *Service) GetLatestCompleted(
	ctx context.Context,
	userID string,
	scenarioID scenario.ScenarioID,
) (Attempt, error) {
	latestAttempt, err := s.attempts.GetLatestCompleted(ctx, userID, scenarioID)
	if err != nil {
		return Attempt{}, fmt.Errorf("get latest completed attempt: %w", err)
	}

	return latestAttempt, nil
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
