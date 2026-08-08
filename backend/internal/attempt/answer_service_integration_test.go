package attempt_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/marlendd/anti-scam-trainer/internal/attempt"
	"github.com/marlendd/anti-scam-trainer/internal/engine"
	"github.com/marlendd/anti-scam-trainer/internal/scenario"
	"github.com/marlendd/anti-scam-trainer/internal/testfixture"
	"github.com/stretchr/testify/require"
)

func TestService_SubmitAnswerLifecycle_Integration(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	testUserID := insertTestUser(t, ctx, db)
	testScenarioID := insertValidTestScenario(t, ctx, db)
	registerFixtureCleanup(t, db, testUserID, testScenarioID)

	attemptRepository := attempt.NewPgRepository(db)
	scenarioRepository := scenario.NewPgRepository(db)
	service := attempt.NewService(
		attemptRepository,
		attemptRepository,
		&scenarioRepository,
	)

	currentAttempt, err := service.Start(ctx, testUserID, testScenarioID)
	require.NoError(t, err)
	require.NotNil(t, currentAttempt.CurrentNodeID)
	require.Equal(t, testfixture.StartNodeID, *currentAttempt.CurrentNodeID)

	invalidInput := attempt.SubmitAnswerInput{
		UserID:         testUserID,
		AttemptID:      currentAttempt.ID,
		NodeID:         testfixture.StartNodeID,
		ChoiceID:       "unknown-choice",
		IdempotencyKey: newIdempotencyKey(t, ctx, db),
	}

	_, err = service.SubmitAnswer(ctx, invalidInput)
	require.ErrorIs(t, err, engine.ErrChoiceNotFound)
	require.Equal(t, 0, countAttemptAnswers(t, ctx, db, currentAttempt.ID))

	afterRollback, err := attemptRepository.GetByID(ctx, currentAttempt.ID, testUserID)
	require.NoError(t, err)
	require.Equal(t, attempt.StatusInProgress, afterRollback.Status)
	require.NotNil(t, afterRollback.CurrentNodeID)
	require.Equal(t, testfixture.StartNodeID, *afterRollback.CurrentNodeID)

	firstInput := attempt.SubmitAnswerInput{
		UserID:         testUserID,
		AttemptID:      currentAttempt.ID,
		NodeID:         testfixture.StartNodeID,
		ChoiceID:       testfixture.StartChoiceID,
		IdempotencyKey: newIdempotencyKey(t, ctx, db),
	}

	firstResult, err := service.SubmitAnswer(ctx, firstInput)
	require.NoError(t, err)
	require.False(t, firstResult.Completed)
	require.NotNil(t, firstResult.NextNodeID)
	require.Equal(t, testfixture.MiddleNodeID, *firstResult.NextNodeID)
	require.Equal(t, 1, countAttemptAnswers(t, ctx, db, currentAttempt.ID))

	afterFirstAnswer, err := attemptRepository.GetByID(ctx, currentAttempt.ID, testUserID)
	require.NoError(t, err)
	require.Equal(t, attempt.StatusInProgress, afterFirstAnswer.Status)
	require.NotNil(t, afterFirstAnswer.CurrentNodeID)
	require.Equal(t, testfixture.MiddleNodeID, *afterFirstAnswer.CurrentNodeID)

	replayedResult, err := service.SubmitAnswer(ctx, firstInput)
	require.NoError(t, err)
	require.Equal(t, firstResult, replayedResult)
	require.Equal(t, 1, countAttemptAnswers(t, ctx, db, currentAttempt.ID))

	conflictingInput := firstInput
	conflictingInput.ChoiceID = "start-choice-2"
	_, err = service.SubmitAnswer(ctx, conflictingInput)
	require.ErrorIs(t, err, attempt.ErrIdempotencyConflict)

	alreadyAnsweredInput := firstInput
	alreadyAnsweredInput.IdempotencyKey = newIdempotencyKey(t, ctx, db)
	_, err = service.SubmitAnswer(ctx, alreadyAnsweredInput)
	require.ErrorIs(t, err, attempt.ErrNodeAlreadyAnswered)
	require.Equal(t, 1, countAttemptAnswers(t, ctx, db, currentAttempt.ID))

	secondInput := attempt.SubmitAnswerInput{
		UserID:         testUserID,
		AttemptID:      currentAttempt.ID,
		NodeID:         testfixture.MiddleNodeID,
		ChoiceID:       "middle-choice-1",
		IdempotencyKey: newIdempotencyKey(t, ctx, db),
	}

	secondResult, err := service.SubmitAnswer(ctx, secondInput)
	require.NoError(t, err)
	require.False(t, secondResult.Completed)
	require.NotNil(t, secondResult.NextNodeID)
	require.Equal(t, testfixture.FinalNodeID, *secondResult.NextNodeID)

	finalInput := attempt.SubmitAnswerInput{
		UserID:         testUserID,
		AttemptID:      currentAttempt.ID,
		NodeID:         testfixture.FinalNodeID,
		ChoiceID:       testfixture.FinalChoiceID,
		IdempotencyKey: newIdempotencyKey(t, ctx, db),
	}

	finalResult, err := service.SubmitAnswer(ctx, finalInput)
	require.NoError(t, err)
	require.True(t, finalResult.Completed)
	require.Nil(t, finalResult.NextNodeID)
	require.NotNil(t, finalResult.EndingID)
	require.Equal(t, testfixture.SafeEndingID, *finalResult.EndingID)
	require.NotNil(t, finalResult.Score)
	require.Equal(t, 100, *finalResult.Score)
	require.Equal(t, 3, countAttemptAnswers(t, ctx, db, currentAttempt.ID))

	completedAttempt, err := attemptRepository.GetByID(
		ctx,
		currentAttempt.ID,
		testUserID,
	)
	require.NoError(t, err)
	require.Equal(t, attempt.StatusCompleted, completedAttempt.Status)
	require.Nil(t, completedAttempt.CurrentNodeID)
	require.NotNil(t, completedAttempt.EndingID)
	require.Equal(t, testfixture.SafeEndingID, *completedAttempt.EndingID)
	require.NotNil(t, completedAttempt.Score)
	require.Equal(t, 100, *completedAttempt.Score)
	require.NotNil(t, completedAttempt.CompletedAt)

	savedFinalAnswer, err := attemptRepository.GetAnswerByAttemptNode(
		ctx,
		currentAttempt.ID,
		testUserID,
		testfixture.FinalNodeID,
	)
	require.NoError(t, err)
	require.Equal(t, scenario.WeightLow, savedFinalAnswer.Weight)
	require.Equal(t, scenario.ScoreSafe, savedFinalAnswer.ChoiceScore)
	require.Equal(t, finalResult, savedFinalAnswer.Response)
}

func TestService_SubmitAnswerConcurrentReplay_Integration(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	testUserID := insertTestUser(t, ctx, db)
	testScenarioID := insertValidTestScenario(t, ctx, db)
	registerFixtureCleanup(t, db, testUserID, testScenarioID)

	attemptRepository := attempt.NewPgRepository(db)
	scenarioRepository := scenario.NewPgRepository(db)
	service := attempt.NewService(
		attemptRepository,
		attemptRepository,
		&scenarioRepository,
	)

	currentAttempt, err := service.Start(ctx, testUserID, testScenarioID)
	require.NoError(t, err)

	input := attempt.SubmitAnswerInput{
		UserID:         testUserID,
		AttemptID:      currentAttempt.ID,
		NodeID:         testfixture.StartNodeID,
		ChoiceID:       testfixture.StartChoiceID,
		IdempotencyKey: newIdempotencyKey(t, ctx, db),
	}

	type submission struct {
		result attempt.SubmitAnswerResult
		err    error
	}

	start := make(chan struct{})
	submissions := make(chan submission, 2)
	for range 2 {
		go func() {
			<-start
			result, err := service.SubmitAnswer(ctx, input)
			submissions <- submission{result: result, err: err}
		}()
	}
	close(start)

	first := <-submissions
	second := <-submissions

	require.NoError(t, first.err)
	require.NoError(t, second.err)
	require.Equal(t, first.result, second.result)
	require.Equal(t, 1, countAttemptAnswers(t, ctx, db, currentAttempt.ID))

	afterSubmission, err := attemptRepository.GetByID(
		ctx,
		currentAttempt.ID,
		testUserID,
	)
	require.NoError(t, err)
	require.Equal(t, attempt.StatusInProgress, afterSubmission.Status)
	require.NotNil(t, afterSubmission.CurrentNodeID)
	require.Equal(t, testfixture.MiddleNodeID, *afterSubmission.CurrentNodeID)
}

func insertValidTestScenario(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
) scenario.ScenarioID {
	t.Helper()

	fixture := testfixture.ValidScenario()
	contentJSON, err := json.Marshal(scenario.Content{
		StartNodeID: fixture.StartNodeID,
		Nodes:       fixture.Nodes,
		Endings:     fixture.Endings,
	})
	require.NoError(t, err)

	const query = `INSERT INTO scenario_versions (
						logical_id,
						version,
						role,
						title,
						description,
						content
					)
					VALUES (gen_random_uuid(), $1, $2, $3, $4, $5::jsonb)
					RETURNING id
	`

	var scenarioID scenario.ScenarioID
	err = db.QueryRowContext(
		ctx,
		query,
		fixture.Version,
		fixture.Role,
		fixture.Title,
		fixture.Description,
		string(contentJSON),
	).Scan(&scenarioID)
	require.NoError(t, err)

	return scenarioID
}

func countAttemptAnswers(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
	attemptID attempt.AttemptID,
) int {
	t.Helper()

	var count int
	err := db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM answers WHERE attempt_id = $1`,
		attemptID,
	).Scan(&count)
	require.NoError(t, err)

	return count
}
