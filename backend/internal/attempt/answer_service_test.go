package attempt_test

import (
	"context"
	"errors"
	"testing"

	"github.com/marlendd/anti-scam-trainer/internal/attempt"
	"github.com/marlendd/anti-scam-trainer/internal/engine"
	"github.com/marlendd/anti-scam-trainer/internal/scenario"
	"github.com/marlendd/anti-scam-trainer/internal/testfixture"
	"github.com/stretchr/testify/require"
)

type answerLookupResult struct {
	answer attempt.Answer
	err    error
}

type advanceAttemptCall struct {
	attemptID     attempt.AttemptID
	userID        string
	currentNodeID scenario.NodeID
	nextNodeID    scenario.NodeID
}

type completeAttemptCall struct {
	attemptID     attempt.AttemptID
	userID        string
	currentNodeID scenario.NodeID
	endingID      scenario.EndingID
	score         int
}

type grantFragmentCall struct {
	userID     string
	scenarioID scenario.ScenarioID
	fragmentID scenario.FragmentID
}

type answerRepositoryFake struct {
	transactionErr   error
	transactionCalls int

	idempotencyResults []answerLookupResult
	idempotencyCalls   int

	currentAttempt  attempt.Attempt
	getAttemptErr   error
	getAttemptCalls int

	nodeAnswer attempt.Answer
	nodeErr    error
	nodeCalls  int

	historyAnswers []attempt.Answer
	historyErr     error
	historyCalls   int

	scoreTotals attempt.ScoreTotals
	scoreErr    error
	scoreCalls  int

	createdAnswer *attempt.Answer
	createErr     error
	createCalls   int

	advanceCall *advanceAttemptCall
	advanceErr  error

	completeCall *completeAttemptCall
	completeErr  error

	grantCall    *grantFragmentCall
	grantAwarded bool
	grantErr     error
}

func (r *answerRepositoryFake) WithinAnswerTransaction(
	ctx context.Context,
	fn func(attempt.AnswerRepository) error,
) error {
	r.transactionCalls++
	if r.transactionErr != nil {
		return r.transactionErr
	}

	return fn(r)
}

func (r *answerRepositoryFake) GetAttemptForUpdate(
	context.Context,
	attempt.AttemptID,
	string,
) (attempt.Attempt, error) {
	r.getAttemptCalls++
	return r.currentAttempt, r.getAttemptErr
}

func (r *answerRepositoryFake) GetAnswerByIdempotencyKey(
	context.Context,
	string,
	attempt.IdempotencyKey,
) (attempt.Answer, error) {
	callIndex := r.idempotencyCalls
	r.idempotencyCalls++

	if callIndex < len(r.idempotencyResults) {
		result := r.idempotencyResults[callIndex]
		return result.answer, result.err
	}

	return attempt.Answer{}, attempt.ErrAnswerNotFound
}

func (r *answerRepositoryFake) GetAnswerByAttemptNode(
	context.Context,
	attempt.AttemptID,
	string,
	scenario.NodeID,
) (attempt.Answer, error) {
	r.nodeCalls++
	return r.nodeAnswer, r.nodeErr
}

func (r *answerRepositoryFake) ListAnswersByAttempt(
	context.Context,
	attempt.AttemptID,
	string,
) ([]attempt.Answer, error) {
	r.historyCalls++
	return r.historyAnswers, r.historyErr
}

func (r *answerRepositoryFake) GetScoreTotals(
	context.Context,
	attempt.AttemptID,
) (attempt.ScoreTotals, error) {
	r.scoreCalls++
	return r.scoreTotals, r.scoreErr
}

func (r *answerRepositoryFake) CreateAnswer(
	_ context.Context,
	answer attempt.Answer,
) (attempt.Answer, error) {
	r.createCalls++
	r.createdAnswer = &answer
	if r.createErr != nil {
		return attempt.Answer{}, r.createErr
	}

	answer.ID = "answer-1"
	return answer, nil
}

func (r *answerRepositoryFake) AdvanceAttempt(
	_ context.Context,
	attemptID attempt.AttemptID,
	userID string,
	currentNodeID scenario.NodeID,
	nextNodeID scenario.NodeID,
) error {
	r.advanceCall = &advanceAttemptCall{
		attemptID:     attemptID,
		userID:        userID,
		currentNodeID: currentNodeID,
		nextNodeID:    nextNodeID,
	}

	return r.advanceErr
}

func (r *answerRepositoryFake) CompleteAttempt(
	_ context.Context,
	attemptID attempt.AttemptID,
	userID string,
	currentNodeID scenario.NodeID,
	endingID scenario.EndingID,
	score int,
) error {
	r.completeCall = &completeAttemptCall{
		attemptID:     attemptID,
		userID:        userID,
		currentNodeID: currentNodeID,
		endingID:      endingID,
		score:         score,
	}

	return r.completeErr
}

func (r *answerRepositoryFake) GrantFragment(
	_ context.Context,
	userID string,
	scenarioID scenario.ScenarioID,
	fragmentID scenario.FragmentID,
) (bool, error) {
	r.grantCall = &grantFragmentCall{
		userID:     userID,
		scenarioID: scenarioID,
		fragmentID: fragmentID,
	}

	return r.grantAwarded, r.grantErr
}

type submitAnswerSetup struct {
	input      attempt.SubmitAnswerInput
	scenario   scenario.Scenario
	repository *answerRepositoryFake
	provider   *scenarioProviderStub
	service    *attempt.Service
}

func newSubmitAnswerSetup() *submitAnswerSetup {
	currentScenario := testfixture.ValidScenario()
	currentNodeID := currentScenario.StartNodeID
	input := attempt.SubmitAnswerInput{
		UserID:         userID,
		AttemptID:      "attempt-1",
		NodeID:         currentNodeID,
		ChoiceID:       testfixture.StartChoiceID,
		IdempotencyKey: "11111111-1111-1111-1111-111111111111",
	}
	repository := &answerRepositoryFake{
		currentAttempt: attempt.Attempt{
			ID:            input.AttemptID,
			UserID:        input.UserID,
			ScenarioID:    currentScenario.ID,
			Status:        attempt.StatusInProgress,
			CurrentNodeID: &currentNodeID,
		},
		nodeErr:      attempt.ErrAnswerNotFound,
		grantAwarded: true,
	}
	provider := &scenarioProviderStub{
		getByIDFn: func(
			context.Context,
			scenario.ScenarioID,
		) (scenario.Scenario, error) {
			return currentScenario, nil
		},
	}

	return &submitAnswerSetup{
		input:      input,
		scenario:   currentScenario,
		repository: repository,
		provider:   provider,
		service: attempt.NewService(
			&attemptRepositoryStub{},
			repository,
			provider,
		),
	}
}

func (s *submitAnswerSetup) moveAttemptToFinalNode() {
	finalNodeID := testfixture.FinalNodeID
	s.input.NodeID = finalNodeID
	s.input.ChoiceID = testfixture.FinalChoiceID
	s.repository.currentAttempt.CurrentNodeID = &finalNodeID
}

func TestServiceSubmitAnswer_AdvancesAttempt(t *testing.T) {
	t.Parallel()

	setup := newSubmitAnswerSetup()

	result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

	require.NoError(t, err)
	require.Equal(t, setup.input.AttemptID, result.AttemptID)
	require.Equal(t, setup.input.NodeID, result.NodeID)
	require.Equal(t, setup.input.ChoiceID, result.ChoiceID)
	require.Equal(t, "Последствие 1", result.Consequence)
	require.False(t, result.Completed)
	require.NotNil(t, result.NextNodeID)
	require.Equal(t, testfixture.MiddleNodeID, *result.NextNodeID)
	require.Nil(t, result.EndingID)
	require.Nil(t, result.Score)

	require.Equal(t, 2, setup.repository.idempotencyCalls)
	require.Equal(t, 1, setup.repository.getAttemptCalls)
	require.Equal(t, 1, setup.repository.nodeCalls)
	require.Zero(t, setup.repository.scoreCalls)
	require.Equal(t, 1, setup.repository.createCalls)
	require.NotNil(t, setup.repository.createdAnswer)
	require.Equal(t, result, setup.repository.createdAnswer.Response)
	require.Equal(t, scenario.WeightLow, setup.repository.createdAnswer.Weight)
	require.Equal(t, scenario.ScoreSafe, setup.repository.createdAnswer.ChoiceScore)
	require.Equal(t, "Объяснение 1", setup.repository.createdAnswer.Explanation)

	require.NotNil(t, setup.repository.advanceCall)
	require.Equal(t, setup.input.AttemptID, setup.repository.advanceCall.attemptID)
	require.Equal(t, setup.input.UserID, setup.repository.advanceCall.userID)
	require.Equal(t, setup.input.NodeID, setup.repository.advanceCall.currentNodeID)
	require.Equal(t, testfixture.MiddleNodeID, setup.repository.advanceCall.nextNodeID)
	require.Nil(t, setup.repository.completeCall)
}

func TestServiceSubmitAnswer_CompletesAttemptAndCalculatesScore(t *testing.T) {
	t.Parallel()

	setup := newSubmitAnswerSetup()
	setup.moveAttemptToFinalNode()
	setup.repository.scoreTotals = attempt.ScoreTotals{
		WeightedScoreSum: 100,
		WeightSum:        2,
	}

	result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

	require.NoError(t, err)
	require.True(t, result.Completed)
	require.Nil(t, result.NextNodeID)
	require.NotNil(t, result.EndingID)
	require.Equal(t, testfixture.SafeEndingID, *result.EndingID)
	require.NotNil(t, result.Score)
	require.Equal(t, 67, *result.Score)
	require.NotNil(t, result.RewardFragmentID)
	require.Equal(t, testfixture.RewardFragmentID, *result.RewardFragmentID)
	require.Equal(t, 1, setup.repository.scoreCalls)
	require.NotNil(t, setup.repository.createdAnswer)
	require.Equal(t, result, setup.repository.createdAnswer.Response)
	require.Nil(t, setup.repository.advanceCall)
	require.NotNil(t, setup.repository.completeCall)
	require.Equal(t, testfixture.SafeEndingID, setup.repository.completeCall.endingID)
	require.Equal(t, 67, setup.repository.completeCall.score)
	require.NotNil(t, setup.repository.grantCall)
	require.Equal(t, setup.input.UserID, setup.repository.grantCall.userID)
	require.Equal(t, setup.scenario.ID, setup.repository.grantCall.scenarioID)
	require.Equal(t, testfixture.RewardFragmentID, setup.repository.grantCall.fragmentID)
}

func TestServiceSubmitAnswer_DoesNotGrantFragmentForRiskyEnding(t *testing.T) {
	t.Parallel()

	setup := newSubmitAnswerSetup()
	setup.moveAttemptToFinalNode()
	setup.input.ChoiceID = testfixture.RiskyFinalChoiceID

	result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

	require.NoError(t, err)
	require.True(t, result.Completed)
	require.NotNil(t, result.EndingID)
	require.Equal(t, testfixture.RiskyEndingID, *result.EndingID)
	require.Nil(t, result.RewardFragmentID)
	require.Nil(t, setup.repository.grantCall)
}

func TestServiceSubmitAnswer_Idempotency(t *testing.T) {
	t.Parallel()

	t.Run("returns an existing result before locking", func(t *testing.T) {
		t.Parallel()

		setup := newSubmitAnswerSetup()
		expected := attempt.SubmitAnswerResult{
			AttemptID:   setup.input.AttemptID,
			NodeID:      setup.input.NodeID,
			ChoiceID:    setup.input.ChoiceID,
			Consequence: "Saved result",
		}
		setup.repository.idempotencyResults = []answerLookupResult{{
			answer: attempt.Answer{
				AttemptID: setup.input.AttemptID,
				NodeID:    setup.input.NodeID,
				ChoiceID:  setup.input.ChoiceID,
				Response:  expected,
			},
		}}

		result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

		require.NoError(t, err)
		require.Equal(t, expected, result)
		require.Equal(t, 1, setup.repository.idempotencyCalls)
		require.Zero(t, setup.repository.getAttemptCalls)
		require.Zero(t, setup.repository.createCalls)
	})

	t.Run("returns a result committed while waiting for the lock", func(t *testing.T) {
		t.Parallel()

		setup := newSubmitAnswerSetup()
		expected := attempt.SubmitAnswerResult{
			AttemptID:   setup.input.AttemptID,
			NodeID:      setup.input.NodeID,
			ChoiceID:    setup.input.ChoiceID,
			Consequence: "Concurrent result",
		}
		setup.repository.idempotencyResults = []answerLookupResult{
			{err: attempt.ErrAnswerNotFound},
			{answer: attempt.Answer{
				AttemptID: setup.input.AttemptID,
				NodeID:    setup.input.NodeID,
				ChoiceID:  setup.input.ChoiceID,
				Response:  expected,
			}},
		}

		result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

		require.NoError(t, err)
		require.Equal(t, expected, result)
		require.Equal(t, 2, setup.repository.idempotencyCalls)
		require.Equal(t, 1, setup.repository.getAttemptCalls)
		require.Zero(t, setup.repository.nodeCalls)
		require.Zero(t, setup.repository.createCalls)
	})

	t.Run("rejects the same key with a different payload", func(t *testing.T) {
		t.Parallel()

		setup := newSubmitAnswerSetup()
		setup.repository.idempotencyResults = []answerLookupResult{{
			answer: attempt.Answer{
				AttemptID: setup.input.AttemptID,
				NodeID:    setup.input.NodeID,
				ChoiceID:  "another-choice",
			},
		}}

		result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

		require.ErrorIs(t, err, attempt.ErrIdempotencyConflict)
		require.Empty(t, result)
		require.Zero(t, setup.repository.getAttemptCalls)
	})
}

func TestServiceSubmitAnswer_RejectsAnsweredNodeAndInvalidAttempt(t *testing.T) {
	t.Parallel()

	t.Run("rejects another key for an answered node", func(t *testing.T) {
		t.Parallel()

		setup := newSubmitAnswerSetup()
		setup.repository.nodeAnswer = attempt.Answer{ID: "answer-existing"}
		setup.repository.nodeErr = nil

		result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

		require.ErrorIs(t, err, attempt.ErrNodeAlreadyAnswered)
		require.Empty(t, result)
		require.Zero(t, setup.repository.createCalls)
	})

	tests := []struct {
		name        string
		mutate      func(*submitAnswerSetup)
		expectedErr error
	}{
		{
			name: "rejects a completed attempt",
			mutate: func(setup *submitAnswerSetup) {
				setup.repository.currentAttempt.Status = attempt.StatusCompleted
			},
			expectedErr: attempt.ErrAttemptNotInProgress,
		},
		{
			name: "rejects an in-progress attempt without a current node",
			mutate: func(setup *submitAnswerSetup) {
				setup.repository.currentAttempt.CurrentNodeID = nil
			},
			expectedErr: attempt.ErrInvalidAttemptState,
		},
		{
			name: "rejects a request for another node",
			mutate: func(setup *submitAnswerSetup) {
				setup.input.NodeID = testfixture.MiddleNodeID
			},
			expectedErr: attempt.ErrAttemptNodeMismatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			setup := newSubmitAnswerSetup()
			tt.mutate(setup)

			result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

			require.ErrorIs(t, err, tt.expectedErr)
			require.Empty(t, result)
			require.Zero(t, setup.repository.createCalls)
		})
	}
}

func TestServiceSubmitAnswer_PreservesDependencyErrors(t *testing.T) {
	t.Parallel()

	t.Run("transaction failure", func(t *testing.T) {
		t.Parallel()

		setup := newSubmitAnswerSetup()
		transactionErr := errors.New("transaction failed")
		setup.repository.transactionErr = transactionErr

		result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

		require.ErrorIs(t, err, transactionErr)
		require.Empty(t, result)
	})

	t.Run("attempt lock failure", func(t *testing.T) {
		t.Parallel()

		setup := newSubmitAnswerSetup()
		lockErr := errors.New("lock failed")
		setup.repository.getAttemptErr = lockErr

		result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

		require.ErrorIs(t, err, lockErr)
		require.Empty(t, result)
	})

	t.Run("answered node lookup failure", func(t *testing.T) {
		t.Parallel()

		setup := newSubmitAnswerSetup()
		lookupErr := errors.New("node lookup failed")
		setup.repository.nodeErr = lookupErr

		result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

		require.ErrorIs(t, err, lookupErr)
		require.Empty(t, result)
	})

	t.Run("scenario provider failure", func(t *testing.T) {
		t.Parallel()

		setup := newSubmitAnswerSetup()
		providerErr := errors.New("scenario provider failed")
		setup.provider.getByIDFn = func(
			context.Context,
			scenario.ScenarioID,
		) (scenario.Scenario, error) {
			return scenario.Scenario{}, providerErr
		}

		result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

		require.ErrorIs(t, err, providerErr)
		require.Empty(t, result)
	})

	t.Run("invalid scenario", func(t *testing.T) {
		t.Parallel()

		setup := newSubmitAnswerSetup()
		invalidScenario := setup.scenario
		invalidScenario.StartNodeID = ""
		setup.provider.getByIDFn = func(
			context.Context,
			scenario.ScenarioID,
		) (scenario.Scenario, error) {
			return invalidScenario, nil
		}

		result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

		require.ErrorIs(t, err, scenario.ErrEmptyStartNodeID)
		require.Empty(t, result)
	})

	t.Run("unknown choice", func(t *testing.T) {
		t.Parallel()

		setup := newSubmitAnswerSetup()
		setup.input.ChoiceID = "unknown-choice"

		result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

		require.ErrorIs(t, err, engine.ErrChoiceNotFound)
		require.Empty(t, result)
	})

	t.Run("answer creation failure", func(t *testing.T) {
		t.Parallel()

		setup := newSubmitAnswerSetup()
		createErr := errors.New("create answer failed")
		setup.repository.createErr = createErr

		result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

		require.ErrorIs(t, err, createErr)
		require.Empty(t, result)
		require.Nil(t, setup.repository.advanceCall)
	})

	t.Run("advance failure", func(t *testing.T) {
		t.Parallel()

		setup := newSubmitAnswerSetup()
		advanceErr := errors.New("advance failed")
		setup.repository.advanceErr = advanceErr

		result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

		require.ErrorIs(t, err, advanceErr)
		require.Empty(t, result)
		require.Equal(t, 1, setup.repository.createCalls)
	})

	t.Run("score totals failure", func(t *testing.T) {
		t.Parallel()

		setup := newSubmitAnswerSetup()
		setup.moveAttemptToFinalNode()
		totalsErr := errors.New("score totals failed")
		setup.repository.scoreErr = totalsErr

		result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

		require.ErrorIs(t, err, totalsErr)
		require.Empty(t, result)
		require.Zero(t, setup.repository.createCalls)
	})

	t.Run("complete failure", func(t *testing.T) {
		t.Parallel()

		setup := newSubmitAnswerSetup()
		setup.moveAttemptToFinalNode()
		completeErr := errors.New("complete failed")
		setup.repository.completeErr = completeErr

		result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

		require.ErrorIs(t, err, completeErr)
		require.Empty(t, result)
		require.Equal(t, 1, setup.repository.createCalls)
	})

	t.Run("fragment grant failure", func(t *testing.T) {
		t.Parallel()

		setup := newSubmitAnswerSetup()
		setup.moveAttemptToFinalNode()
		grantErr := errors.New("fragment grant failed")
		setup.repository.grantErr = grantErr

		result, err := setup.service.SubmitAnswer(context.Background(), setup.input)

		require.ErrorIs(t, err, grantErr)
		require.Empty(t, result)
		require.Zero(t, setup.repository.createCalls)
		require.Nil(t, setup.repository.completeCall)
	})
}
