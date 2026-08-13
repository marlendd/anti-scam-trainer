package attempt_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/marlendd/anti-scam-trainer/internal/attempt"
	"github.com/marlendd/anti-scam-trainer/internal/platform/postgres"
	"github.com/marlendd/anti-scam-trainer/internal/scenario"
	"github.com/marlendd/anti-scam-trainer/internal/testfixture"
	"github.com/stretchr/testify/require"
)

const defaultTestDatabaseURL = "postgres://postgres:password@127.0.0.1:5433/postgres?sslmode=disable"

func TestPgRepository_AttemptLifecycle_Integration(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	testUserID := insertTestUser(t, ctx, db)
	testScenarioID := insertTestScenario(t, ctx, db)
	registerFixtureCleanup(t, db, testUserID, testScenarioID)

	repository := attempt.NewPgRepository(db)
	startNodeID := scenario.NodeID("node-start")

	created, err := repository.Create(ctx, testUserID, testScenarioID, startNodeID)
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)
	require.Equal(t, testUserID, created.UserID)
	require.Equal(t, testScenarioID, created.ScenarioID)
	require.Equal(t, attempt.StatusInProgress, created.Status)
	require.NotNil(t, created.CurrentNodeID)
	require.Equal(t, startNodeID, *created.CurrentNodeID)
	require.Nil(t, created.EndingID)
	require.Nil(t, created.Score)
	require.Nil(t, created.CompletedAt)
	require.False(t, created.StartedAt.IsZero())
	require.False(t, created.UpdatedAt.IsZero())

	_, err = repository.Create(ctx, testUserID, testScenarioID, startNodeID)
	require.ErrorIs(t, err, attempt.ErrActiveAttemptExists)

	byID, err := repository.GetByID(ctx, created.ID, testUserID)
	require.NoError(t, err)
	require.Equal(t, created, byID)

	// Новый экземпляр репозитория имитирует повторный вход после перезагрузки страницы.
	reloadedRepository := attempt.NewPgRepository(db)
	active, err := reloadedRepository.GetActive(ctx, testUserID, testScenarioID)
	require.NoError(t, err)
	require.Equal(t, created, active)

	err = repository.Abort(ctx, created.ID, testUserID)
	require.NoError(t, err)

	_, err = repository.GetActive(ctx, testUserID, testScenarioID)
	require.ErrorIs(t, err, attempt.ErrActiveAttemptNotFound)

	aborted, err := repository.GetByID(ctx, created.ID, testUserID)
	require.NoError(t, err)
	require.Equal(t, attempt.StatusAborted, aborted.Status)
	require.NotNil(t, aborted.CurrentNodeID)
	require.Equal(t, startNodeID, *aborted.CurrentNodeID)
	require.Nil(t, aborted.EndingID)
	require.Nil(t, aborted.Score)
	require.Nil(t, aborted.CompletedAt)
	require.False(t, aborted.UpdatedAt.Before(created.UpdatedAt))

	err = repository.Abort(ctx, created.ID, testUserID)
	require.ErrorIs(t, err, attempt.ErrAttemptNotInProgress)
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("set RUN_INTEGRATION_TESTS=1 to run integration tests")
	}

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = defaultTestDatabaseURL
	}

	db, err := sql.Open("pgx", databaseURL)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	require.NoError(t, db.PingContext(ctx))
	require.NoError(t, postgres.RunMigrations(databaseURL, "../../migrations"))

	return db
}

func insertTestUser(t *testing.T, ctx context.Context, db *sql.DB) string {
	t.Helper()

	const query = `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`

	email := fmt.Sprintf("attempt_repository_%d@example.com", time.Now().UnixNano())

	var id string
	err := db.QueryRowContext(ctx, query, email, "test-password-hash").Scan(&id)
	require.NoError(t, err)

	return id
}

func insertTestScenario(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
) scenario.ScenarioID {
	t.Helper()

	const query = `
		INSERT INTO scenario_versions (
			logical_id,
			version,
			role,
			title,
			description,
			content
		)
		VALUES (gen_random_uuid(), 1, 'buyer', $1, $2, '{}'::jsonb)
		RETURNING id
	`

	var id scenario.ScenarioID
	err := db.QueryRowContext(
		ctx,
		query,
		"Integration test scenario",
		"Scenario created by attempt repository integration test",
	).Scan(&id)
	require.NoError(t, err)

	return id
}

func registerFixtureCleanup(
	t *testing.T,
	db *sql.DB,
	testUserID string,
	testScenarioID scenario.ScenarioID,
) {
	t.Helper()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, testUserID)
		require.NoError(t, err)

		_, err = db.ExecContext(ctx, `DELETE FROM scenario_versions WHERE id = $1`, testScenarioID)
		require.NoError(t, err)
	})
}

func TestPgRepository_GetByIDRejectsOtherUser_Integration(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	ownerID := insertTestUser(t, ctx, db)
	otherUserID := insertTestUser(t, ctx, db)
	testScenarioID := insertTestScenario(t, ctx, db)
	registerFixtureCleanup(t, db, ownerID, testScenarioID)
	t.Cleanup(func() {
		_, err := db.ExecContext(context.Background(), `DELETE FROM users WHERE id = $1`, otherUserID)
		require.NoError(t, err)
	})

	repository := attempt.NewPgRepository(db)
	created, err := repository.Create(ctx, ownerID, testScenarioID, "node-start")
	require.NoError(t, err)

	_, err = repository.GetByID(ctx, created.ID, otherUserID)
	require.ErrorIs(t, err, attempt.ErrAttemptNotFound)
}

func TestPgRepository_GetLatestCompleted_Integration(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	testUserID := insertTestUser(t, ctx, db)
	testScenarioID := insertTestScenario(t, ctx, db)
	registerFixtureCleanup(t, db, testUserID, testScenarioID)

	olderCompletedAt := time.Now().Add(-time.Hour).UTC()
	newerCompletedAt := olderCompletedAt.Add(30 * time.Minute)

	olderID := insertCompletedAttempt(
		t,
		ctx,
		db,
		testUserID,
		testScenarioID,
		"ending-older",
		50,
		olderCompletedAt,
	)
	newerID := insertCompletedAttempt(
		t,
		ctx,
		db,
		testUserID,
		testScenarioID,
		"ending-newer",
		100,
		newerCompletedAt,
	)

	repository := attempt.NewPgRepository(db)
	latest, err := repository.GetLatestCompleted(ctx, testUserID, testScenarioID)

	require.NoError(t, err)
	require.Equal(t, newerID, latest.ID)
	require.NotEqual(t, olderID, latest.ID)
	require.Equal(t, attempt.StatusCompleted, latest.Status)
	require.Nil(t, latest.CurrentNodeID)
	require.NotNil(t, latest.EndingID)
	require.Equal(t, scenario.EndingID("ending-newer"), *latest.EndingID)
	require.NotNil(t, latest.Score)
	require.Equal(t, 100, *latest.Score)
	require.NotNil(t, latest.CompletedAt)
	require.WithinDuration(t, newerCompletedAt, *latest.CompletedAt, time.Microsecond)
}

func TestPgRepository_GetLatestCompletedNotFound_Integration(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	testUserID := insertTestUser(t, ctx, db)
	testScenarioID := insertTestScenario(t, ctx, db)
	registerFixtureCleanup(t, db, testUserID, testScenarioID)

	repository := attempt.NewPgRepository(db)
	_, err := repository.GetLatestCompleted(ctx, testUserID, testScenarioID)

	require.ErrorIs(t, err, attempt.ErrCompletedAttemptNotFound)
}

func TestService_RestartPreservesCompletedResults_Integration(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	testUserID := insertTestUser(t, ctx, db)
	testScenarioID := insertTestScenario(t, ctx, db)
	registerFixtureCleanup(t, db, testUserID, testScenarioID)

	completedAt := time.Now().Add(-time.Hour).UTC()
	completedID := insertCompletedAttempt(
		t,
		ctx,
		db,
		testUserID,
		testScenarioID,
		testfixture.SafeEndingID,
		100,
		completedAt,
	)
	insertCompletedAnswer(t, ctx, db, completedID)

	repository := attempt.NewPgRepository(db)
	activeBeforeRestart, err := repository.Create(
		ctx,
		testUserID,
		testScenarioID,
		testfixture.MiddleNodeID,
	)
	require.NoError(t, err)

	fixture := testfixture.ValidScenario()
	fixture.ID = testScenarioID
	service := attempt.NewService(
		repository,
		repository,
		integrationScenarioProvider{currentScenario: fixture},
	)

	activeAfterRestart, err := service.Restart(ctx, testUserID, testScenarioID)
	require.NoError(t, err)
	require.NotEqual(t, activeBeforeRestart.ID, activeAfterRestart.ID)
	require.Equal(t, attempt.StatusInProgress, activeAfterRestart.Status)
	require.NotNil(t, activeAfterRestart.CurrentNodeID)
	require.Equal(t, fixture.StartNodeID, *activeAfterRestart.CurrentNodeID)

	aborted, err := repository.GetByID(ctx, activeBeforeRestart.ID, testUserID)
	require.NoError(t, err)
	require.Equal(t, attempt.StatusAborted, aborted.Status)

	latestCompleted, err := service.GetLatestCompleted(ctx, testUserID, testScenarioID)
	require.NoError(t, err)
	require.Equal(t, completedID, latestCompleted.Attempt.ID)
	require.Equal(t, attempt.StatusCompleted, latestCompleted.Attempt.Status)
	require.Len(t, latestCompleted.History, 1)
	require.Equal(t, testfixture.StartChoiceID, latestCompleted.History[0].SelectedChoice.ID)
	require.Equal(t, scenario.ScoreSafe, latestCompleted.History[0].ChoiceScore)
	require.Equal(t, "Сохранённое последствие", latestCompleted.History[0].Consequence)
	require.Equal(t, "Сохранённое объяснение", latestCompleted.History[0].Explanation)
	require.NotNil(t, latestCompleted.Ending)
	require.Equal(t, testfixture.SafeEndingID, latestCompleted.Ending.ID)

	var completedAttempts int
	err = db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM attempts WHERE id = $1 AND status = 'completed'`,
		completedID,
	).Scan(&completedAttempts)
	require.NoError(t, err)
	require.Equal(t, 1, completedAttempts)

	var savedAnswers int
	err = db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM answers WHERE attempt_id = $1`,
		completedID,
	).Scan(&savedAnswers)
	require.NoError(t, err)
	require.Equal(t, 1, savedAnswers)
}

type integrationScenarioProvider struct {
	currentScenario scenario.Scenario
}

func (p integrationScenarioProvider) GetActiveByID(
	context.Context,
	scenario.ScenarioID,
) (scenario.Scenario, error) {
	return p.currentScenario, nil
}

func (p integrationScenarioProvider) GetByID(
	context.Context,
	scenario.ScenarioID,
) (scenario.Scenario, error) {
	return p.currentScenario, nil
}

func insertCompletedAnswer(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
	attemptID attempt.AttemptID,
) {
	t.Helper()

	const query = `
		INSERT INTO answers (
			attempt_id,
			node_id,
			choice_id,
			idempotency_key,
			weight,
			choice_score,
			risk_categories,
			consequence,
			explanation,
			response
		)
		VALUES ($1, $2, $3, gen_random_uuid(), 1, 100, '[]'::jsonb, $4, $5, '{}'::jsonb)
	`

	_, err := db.ExecContext(
		ctx,
		query,
		attemptID,
		testfixture.StartNodeID,
		testfixture.StartChoiceID,
		"Сохранённое последствие",
		"Сохранённое объяснение",
	)
	require.NoError(t, err)
}

func insertCompletedAttempt(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
	userID string,
	scenarioID scenario.ScenarioID,
	endingID scenario.EndingID,
	score int,
	completedAt time.Time,
) attempt.AttemptID {
	t.Helper()

	const query = `
		INSERT INTO attempts (
			user_id,
			scenario_id,
			status,
			ending_id,
			score,
			completed_at,
			updated_at
		)
		VALUES ($1, $2, 'completed', $3, $4, $5, $5)
		RETURNING id
	`

	var id attempt.AttemptID
	err := db.QueryRowContext(
		ctx,
		query,
		userID,
		scenarioID,
		endingID,
		score,
		completedAt,
	).Scan(&id)
	require.NoError(t, err)

	return id
}

func TestPgRepository_WithinTransactionRollsBack_Integration(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	testUserID := insertTestUser(t, ctx, db)
	testScenarioID := insertTestScenario(t, ctx, db)
	registerFixtureCleanup(t, db, testUserID, testScenarioID)

	repository := attempt.NewPgRepository(db)
	callbackErr := errors.New("force transaction rollback")

	var createdAttemptID attempt.AttemptID
	err := repository.WithinTransaction(ctx, func(txRepository attempt.AttemptRepository) error {
		created, err := txRepository.Create(
			ctx,
			testUserID,
			testScenarioID,
			"node-start",
		)
		if err != nil {
			return fmt.Errorf("create attempt inside transaction: %w", err)
		}

		createdAttemptID = created.ID

		return callbackErr
	})

	require.ErrorIs(t, err, callbackErr)
	require.NotEmpty(t, createdAttemptID)

	_, err = repository.GetByID(ctx, createdAttemptID, testUserID)
	require.ErrorIs(t, err, attempt.ErrAttemptNotFound)

	_, err = repository.GetActive(ctx, testUserID, testScenarioID)
	require.ErrorIs(t, err, attempt.ErrActiveAttemptNotFound)
}
