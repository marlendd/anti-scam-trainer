package evaluation_test

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/marlendd/anti-scam-trainer/internal/evaluation"
	"github.com/marlendd/anti-scam-trainer/internal/feedback"
	"github.com/marlendd/anti-scam-trainer/internal/platform/postgres"
	"github.com/marlendd/anti-scam-trainer/internal/progress"
	"github.com/stretchr/testify/require"
)

type mockLLM struct {
	shouldFail bool
	respJSON   string
}

func (m *mockLLM) GenerateJSON(ctx context.Context, sys, user string) (string, error) {
	if m.shouldFail {
		return "", errors.New("mock llm timeout or error")
	}
	return m.respJSON, nil
}

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

	t.Run("Category Dashboard Progress", func(t *testing.T) {
		_, err := db.Exec("TRUNCATE users, scenario_versions, attempts, answers CASCADE")
		require.NoError(t, err)

		seedSQL := `
			INSERT INTO users (id, email, password_hash) VALUES ('00000000-0000-0000-0000-000000000001', 'cat@test.com', 'hash');
			INSERT INTO scenario_versions (id, logical_id, version, role, title, description, content)
			VALUES ('00000000-0000-0000-0000-000000000002', gen_random_uuid(), 1, 'buyer', 'title', 'desc', '{"risk_categories": ["phishing", "fake_payment"]}'::jsonb);
			INSERT INTO attempts (id, user_id, scenario_id, status, score, completed_at, ending_id)
			VALUES ('120b7935-62bf-4fd8-828a-6bbe7ef7a19a', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', 'completed', 100, now(), 'end1');
		`
		_, err = db.Exec(seedSQL)
		require.NoError(t, err)

		dash, err := svc.GetUserCategoryDashboard(ctx, "00000000-0000-0000-0000-000000000001")
		require.NoError(t, err)

		require.Equal(t, 1, dash.TotalCompleted)
		require.Len(t, dash.Stats, 2)
		require.ElementsMatch(t, []string{"phishing", "fake_payment"}, []string{dash.Stats[0].Category, dash.Stats[1].Category})
	})

	t.Run("Leaderboard and Rank History", func(t *testing.T) {
		_, err := db.Exec("TRUNCATE users, scenario_versions, attempts, answers, user_inventory CASCADE")
		require.NoError(t, err)

		u1 := "11111111-0000-0000-0000-000000000001"
		u2 := "22222222-0000-0000-0000-000000000002"

		seedSQL := `
			INSERT INTO users (id, email, password_hash) VALUES 
			  ('` + u1 + `', 'u1@test.com', 'h'),
			  ('` + u2 + `', 'u2@test.com', 'h');
			
			INSERT INTO scenario_versions (id, logical_id, version, role, title, description, content) VALUES 
			  ('aaaa0000-0000-0000-0000-000000000000', gen_random_uuid(), 1, 'buyer', 't1', 'd', '{}'::jsonb),
			  ('bbbb0000-0000-0000-0000-000000000000', gen_random_uuid(), 1, 'buyer', 't2', 'd', '{}'::jsonb);

			-- Вчерашние попытки (INTERVAL '2 days', чтобы точно считалось как "до вчерашнего дня")
			INSERT INTO attempts (id, user_id, scenario_id, status, score, completed_at, ending_id) VALUES 
			  (gen_random_uuid(), '` + u2 + `', 'aaaa0000-0000-0000-0000-000000000000', 'completed', 100, NOW() - INTERVAL '2 days', 'end'),
			  (gen_random_uuid(), '` + u1 + `', 'aaaa0000-0000-0000-0000-000000000000', 'completed', 50, NOW() - INTERVAL '2 days', 'end');

			-- Сегодняшняя попытка: Юзер 1 прошел второй сценарий и набрал 100 очков (итого 150)
			INSERT INTO attempts (id, user_id, scenario_id, status, score, completed_at, ending_id) VALUES 
			  (gen_random_uuid(), '` + u1 + `', 'bbbb0000-0000-0000-0000-000000000000', 'completed', 100, NOW(), 'end');
		`
		_, err = db.Exec(seedSQL)
		require.NoError(t, err)

		lb, err := svc.GetLeaderboard(ctx, 10, 0)
		require.NoError(t, err)
		require.Len(t, lb.Entries, 2)

		require.Equal(t, 1, lb.Entries[0].Rank)
		require.Equal(t, "u1@test.com", lb.Entries[0].Player)
		require.Equal(t, 150, lb.Entries[0].Score)
		require.NotNil(t, lb.Entries[0].RankChange)
		require.Equal(t, 1, *lb.Entries[0].RankChange)

		require.Equal(t, 2, lb.Entries[1].Rank)
		require.Equal(t, "u2@test.com", lb.Entries[1].Player)
		require.Equal(t, 100, lb.Entries[1].Score)
		require.NotNil(t, lb.Entries[1].RankChange)
		require.Equal(t, -1, *lb.Entries[1].RankChange)

		history, err := svc.GetMyRankHistory(ctx, u1)
		require.NoError(t, err)
		require.NotEmpty(t, history.History)

		lastIdx := len(history.History) - 1
		require.Equal(t, 1, history.History[lastIdx].Rank)
	})

	t.Run("Feedback Fallback Logic", func(t *testing.T) {
		_, err := db.Exec("TRUNCATE users, scenario_versions, attempts, answers CASCADE")
		require.NoError(t, err)

		uid := "fbfbfbfb-0000-0000-0000-000000000001"
		aid := "a1a1a1a1-0000-0000-0000-000000000001"
		scenarioID := "33333333-3333-3333-3333-333333333333"

		seedSQL := `
			INSERT INTO users (id, email, password_hash) VALUES ('` + uid + `', 'fb@test.com', 'h');
			INSERT INTO scenario_versions (id, logical_id, version, role, title, description, content)
			VALUES ('` + scenarioID + `', gen_random_uuid(), 1, 'buyer', 'FB', 'desc', '{}'::jsonb);
			
			INSERT INTO attempts (id, user_id, scenario_id, status, score, completed_at, ending_id)
			VALUES ('` + aid + `', '` + uid + `', '` + scenarioID + `', 'completed', 50, now(), 'end');

			-- Плохой ответ (попался на фишинг)
			INSERT INTO answers (attempt_id, node_id, choice_id, idempotency_key, weight, choice_score, risk_categories, consequence, explanation, response)
			VALUES ('` + aid + `', 'node1', 'c1', gen_random_uuid(), 1, 0, '["phishing"]'::jsonb, 'Перешел по поддельной ссылке', 'Сайт украл данные', '{"node_question":"q1", "choice_text":"c1"}');
			
			-- Хороший ответ (не перевел деньги)
			INSERT INTO answers (attempt_id, node_id, choice_id, idempotency_key, weight, choice_score, risk_categories, consequence, explanation, response)
			VALUES ('` + aid + `', 'node2', 'c2', gen_random_uuid(), 1, 100, '[]'::jsonb, 'Деньги сохранены', 'Не перевел предоплату', '{"node_question":"q2", "choice_text":"c2"}');
		`
		_, err = db.Exec(seedSQL)
		require.NoError(t, err)

		fbRepo := feedback.NewPgRepository(db, slog.Default())
		fbService := feedback.NewService(fbRepo, &mockLLM{shouldFail: true}, slog.Default())

		fb, err := fbService.Generate(ctx, uid, aid)
		require.NoError(t, err)

		require.Contains(t, fb.Weaknesses, "Перешел по поддельной ссылке")
		require.Contains(t, fb.Strengths, "Не перевел предоплату")
		require.Equal(t, "phishing", fb.RiskProfile.DominantRisk)
		require.Equal(t, 1, fb.RiskProfile.RiskCount)
		require.NotEmpty(t, fb.Recommendations)
	})
}
