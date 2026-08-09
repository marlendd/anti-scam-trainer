package attempt_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/marlendd/anti-scam-trainer/internal/attempt"
	"github.com/marlendd/anti-scam-trainer/internal/scenario"
	"github.com/stretchr/testify/require"
)

func TestPgRepository_AnswerLifecycle_Integration(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	testUserID := insertTestUser(t, ctx, db)
	testScenarioID := insertTestScenario(t, ctx, db)
	registerFixtureCleanup(t, db, testUserID, testScenarioID)

	repository := attempt.NewPgRepository(db)
	startNodeID := scenario.NodeID("node-start")
	nextNodeID := scenario.NodeID("node-next")
	endingID := scenario.EndingID("ending-safe")

	createdAttempt, err := repository.Create(
		ctx,
		testUserID,
		testScenarioID,
		startNodeID,
	)
	require.NoError(t, err)

	firstKey := newIdempotencyKey(t, ctx, db)
	firstResponse := attempt.SubmitAnswerResult{
		AttemptID:   createdAttempt.ID,
		NodeID:      startNodeID,
		ChoiceID:    "choice-safe",
		Consequence: "Safe consequence",
		NextNodeID:  &nextNodeID,
	}
	firstAnswer := attempt.Answer{
		AttemptID:      createdAttempt.ID,
		NodeID:         startNodeID,
		ChoiceID:       firstResponse.ChoiceID,
		IdempotencyKey: firstKey,
		Weight:         scenario.WeightMedium,
		ChoiceScore:    scenario.ScoreSafe,
		RiskCategories: []scenario.RiskCategory{},
		Consequence:    firstResponse.Consequence,
		Explanation:    "Safe explanation",
		Response:       firstResponse,
	}

	var savedFirstAnswer attempt.Answer
	err = repository.WithinAnswerTransaction(
		ctx,
		func(txRepository attempt.AnswerRepository) error {
			lockedAttempt, err := txRepository.GetAttemptForUpdate(
				ctx,
				createdAttempt.ID,
				testUserID,
			)
			if err != nil {
				return fmt.Errorf("lock attempt: %w", err)
			}
			if lockedAttempt.ID != createdAttempt.ID {
				return fmt.Errorf("locked unexpected attempt %q", lockedAttempt.ID)
			}

			_, err = txRepository.GetAnswerByIdempotencyKey(ctx, testUserID, firstKey)
			if !errors.Is(err, attempt.ErrAnswerNotFound) {
				return fmt.Errorf("expected answer not found, got %w", err)
			}

			savedFirstAnswer, err = txRepository.CreateAnswer(ctx, firstAnswer)
			if err != nil {
				return fmt.Errorf("create first answer: %w", err)
			}

			return txRepository.AdvanceAttempt(
				ctx,
				createdAttempt.ID,
				testUserID,
				startNodeID,
				nextNodeID,
			)
		},
	)
	require.NoError(t, err)
	require.NotEmpty(t, savedFirstAnswer.ID)
	require.False(t, savedFirstAnswer.CreatedAt.IsZero())
	require.Equal(t, firstResponse, savedFirstAnswer.Response)

	byKey, err := repository.GetAnswerByIdempotencyKey(ctx, testUserID, firstKey)
	require.NoError(t, err)
	require.Equal(t, savedFirstAnswer, byKey)

	byNode, err := repository.GetAnswerByAttemptNode(
		ctx,
		createdAttempt.ID,
		testUserID,
		startNodeID,
	)
	require.NoError(t, err)
	require.Equal(t, savedFirstAnswer, byNode)

	totals, err := repository.GetScoreTotals(ctx, createdAttempt.ID)
	require.NoError(t, err)
	require.Equal(t, 200, totals.WeightedScoreSum)
	require.Equal(t, 2, totals.WeightSum)

	advancedAttempt, err := repository.GetByID(ctx, createdAttempt.ID, testUserID)
	require.NoError(t, err)
	require.NotNil(t, advancedAttempt.CurrentNodeID)
	require.Equal(t, nextNodeID, *advancedAttempt.CurrentNodeID)

	secondKey := newIdempotencyKey(t, ctx, db)
	finalScore := 75
	secondResponse := attempt.SubmitAnswerResult{
		AttemptID:   createdAttempt.ID,
		NodeID:      nextNodeID,
		ChoiceID:    "choice-risky",
		Consequence: "Risky consequence",
		EndingID:    &endingID,
		Completed:   true,
		Score:       &finalScore,
	}
	secondAnswer := attempt.Answer{
		AttemptID:      createdAttempt.ID,
		NodeID:         nextNodeID,
		ChoiceID:       secondResponse.ChoiceID,
		IdempotencyKey: secondKey,
		Weight:         scenario.WeightMedium,
		ChoiceScore:    scenario.ScoreRisky,
		RiskCategories: []scenario.RiskCategory{scenario.RiskExternalLink},
		Consequence:    secondResponse.Consequence,
		Explanation:    "Risky explanation",
		Response:       secondResponse,
	}

	err = repository.WithinAnswerTransaction(
		ctx,
		func(txRepository attempt.AnswerRepository) error {
			if _, err := txRepository.GetAttemptForUpdate(
				ctx,
				createdAttempt.ID,
				testUserID,
			); err != nil {
				return fmt.Errorf("lock advanced attempt: %w", err)
			}

			if _, err := txRepository.CreateAnswer(ctx, secondAnswer); err != nil {
				return fmt.Errorf("create second answer: %w", err)
			}

			return txRepository.CompleteAttempt(
				ctx,
				createdAttempt.ID,
				testUserID,
				nextNodeID,
				endingID,
				finalScore,
			)
		},
	)
	require.NoError(t, err)

	completedAttempt, err := repository.GetByID(ctx, createdAttempt.ID, testUserID)
	require.NoError(t, err)
	require.Equal(t, attempt.StatusCompleted, completedAttempt.Status)
	require.Nil(t, completedAttempt.CurrentNodeID)
	require.NotNil(t, completedAttempt.EndingID)
	require.Equal(t, endingID, *completedAttempt.EndingID)
	require.NotNil(t, completedAttempt.Score)
	require.Equal(t, finalScore, *completedAttempt.Score)
	require.NotNil(t, completedAttempt.CompletedAt)

	duplicateKeyAnswer := secondAnswer
	duplicateKeyAnswer.NodeID = "another-node"
	_, err = repository.CreateAnswer(ctx, duplicateKeyAnswer)
	require.ErrorIs(t, err, attempt.ErrIdempotencyConflict)

	duplicateNodeAnswer := firstAnswer
	duplicateNodeAnswer.IdempotencyKey = newIdempotencyKey(t, ctx, db)
	_, err = repository.CreateAnswer(ctx, duplicateNodeAnswer)
	require.ErrorIs(t, err, attempt.ErrNodeAlreadyAnswered)
}

func TestPgRepository_GrantFragmentIsIdempotent_Integration(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	testUserID := insertTestUser(t, ctx, db)
	testScenarioID := insertTestScenario(t, ctx, db)
	registerFixtureCleanup(t, db, testUserID, testScenarioID)

	repository := attempt.NewPgRepository(db)
	fragmentID := scenario.FragmentID("safe-deal-piece-test")

	awarded, err := repository.GrantFragment(
		ctx,
		testUserID,
		testScenarioID,
		fragmentID,
	)
	require.NoError(t, err)
	require.True(t, awarded)

	awarded, err = repository.GrantFragment(
		ctx,
		testUserID,
		testScenarioID,
		fragmentID,
	)
	require.NoError(t, err)
	require.False(t, awarded)

	var (
		savedScenarioID scenario.ScenarioID
		savedFragmentID scenario.FragmentID
		count           int
	)
	err = db.QueryRowContext(
		ctx,
		`SELECT MIN(scenario_id::text), MIN(fragment_id), COUNT(*)
		 FROM user_inventory
		 WHERE user_id = $1`,
		testUserID,
	).Scan(&savedScenarioID, &savedFragmentID, &count)
	require.NoError(t, err)
	require.Equal(t, testScenarioID, savedScenarioID)
	require.Equal(t, fragmentID, savedFragmentID)
	require.Equal(t, 1, count)
}

func newIdempotencyKey(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
) attempt.IdempotencyKey {
	t.Helper()

	var key attempt.IdempotencyKey
	require.NoError(t, db.QueryRowContext(ctx, `SELECT gen_random_uuid()`).Scan(&key))

	return key
}
