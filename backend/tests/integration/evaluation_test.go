package evaluation_test

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/marlendd/anti-scam-trainer/internal/evaluation"
	"github.com/marlendd/anti-scam-trainer/internal/platform/postgres"
	"github.com/marlendd/anti-scam-trainer/internal/progress"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T, dbURL string) *sql.DB {
	t.Helper()
	db, err := sql.Open("postgres", dbURL)
	require.NoError(t, err)

	return db
}

func TestEvaluation_Integration(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("set RUN_INTEGRATION_TESTS=1 to run integration tests")
	}

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@127.0.0.1:5433/antiscam?sslmode=disable"
	}
	err := postgres.RunMigrations(dbURL, "../../migrations")
	require.NoError(t, err, "failed to run migrations")

	db := setupTestDB(t, dbURL)

	defer func() {
		err := db.Close()
		require.NoError(t, err)
	}()

	evaluator := evaluation.NewEvaluator()

	repo := progress.NewPgRepository(db, slog.Default())
	svc := progress.NewService(repo, evaluator)
	ctx := context.Background()

	_, err = db.Exec("TRUNCATE users, scenario_versions, attempts, answers CASCADE")
	require.NoError(t, err)

	userID := "00000000-0000-0000-0000-000000000001"
	attemptID := "120b7935-62bf-4fd8-828a-6bbe7ef7a19a"

	t.Run("Seed and Calculate Score", func(t *testing.T) {
		seedSQL := `
			INSERT INTO users (id, email, password_hash) VALUES ('00000000-0000-0000-0000-000000000001', 'test@test.com', 'hash');
			INSERT INTO scenario_versions (id, logical_id, version, role, title, description, reward_fragment_id, content)
			VALUES ('00000000-0000-0000-0000-000000000002', gen_random_uuid(), 1, 'buyer', 'title', 'desc', 'safe-deal-piece-test', '{}'::jsonb);
			INSERT INTO attempts (id, user_id, scenario_id, status, current_node_id) 
			VALUES ('120b7935-62bf-4fd8-828a-6bbe7ef7a19a', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', 'in_progress', 'start_node');
			
			-- Тут оценка 50, значит риск ОБЯЗАТЕЛЕН
			INSERT INTO answers (attempt_id, node_id, choice_id, idempotency_key, weight, choice_score, risk_categories, consequence, explanation, response)
			VALUES ('120b7935-62bf-4fd8-828a-6bbe7ef7a19a', 'node1', 'c1', gen_random_uuid(), 2, 50, '["suspicious_link"]'::jsonb, 'cons', 'expl', '{}');
			
			-- Тут оценка 100, риск может быть пустым
			INSERT INTO answers (attempt_id, node_id, choice_id, idempotency_key, weight, choice_score, risk_categories, consequence, explanation, response)
			VALUES ('120b7935-62bf-4fd8-828a-6bbe7ef7a19a', 'node2', 'c2', gen_random_uuid(), 1, 100, '[]'::jsonb, 'cons', 'expl', '{}');
		`
		_, err := db.Exec(seedSQL)
		require.NoError(t, err)

		score, err := svc.GetAttemptResults(ctx, userID, attemptID)
		require.NoError(t, err)
		require.Equal(t, 67, score)
	})

	t.Run("Verify Personal Progress Stats", func(t *testing.T) {
		stats, err := repo.GetUserStatsByRole(ctx, userID)
		require.NoError(t, err)
		require.NotEmpty(t, stats)

		require.Len(t, stats, 1)

		require.Equal(t, "buyer", stats[0].Role)
		require.Equal(t, int64(1), stats[0].InProgressCount)
		require.Equal(t, int64(0), stats[0].CompletedCount)
	})

	t.Run("Read Puzzle Progress", func(t *testing.T) {
		_, err := db.Exec(`
			INSERT INTO user_inventory (user_id, scenario_id, fragment_id)
			VALUES (
				'00000000-0000-0000-0000-000000000001',
				'00000000-0000-0000-0000-000000000002',
				'safe-deal-piece-test'
			)
		`)
		require.NoError(t, err)

		puzzleProgress, err := svc.GetUserPuzzleProgress(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, 1, puzzleProgress.EarnedCount)
		require.Equal(t, 1, puzzleProgress.TotalCount)
		require.Len(t, puzzleProgress.Fragments, 1)
		require.Equal(
			t,
			"00000000-0000-0000-0000-000000000002",
			puzzleProgress.Fragments[0].ScenarioID,
		)
		require.Equal(t, "safe-deal-piece-test", puzzleProgress.Fragments[0].FragmentID)
	})

	t.Run("Verify Leaderboard Empty State", func(t *testing.T) {
		resp, err := svc.GetLeaderboard(ctx, 10, 0)
		require.NoError(t, err)

		require.Empty(t, resp.Entries)
	})
}
