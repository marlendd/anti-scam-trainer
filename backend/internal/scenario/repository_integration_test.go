package scenario_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/marlendd/anti-scam-trainer/internal/platform/postgres"
	"github.com/marlendd/anti-scam-trainer/internal/scenario"
	"github.com/marlendd/anti-scam-trainer/internal/testfixture"
	"github.com/stretchr/testify/require"
)

const defaultTestDatabaseURL = "postgres://postgres:postgres@127.0.0.1:5433/antiscam?sslmode=disable"

func TestPgRepository_GetScenario_Integration(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	activeID := insertScenario(t, ctx, db, true)
	inactiveID := insertScenario(t, ctx, db, false)
	t.Cleanup(func() {
		deleteScenario(t, ctx, db, activeID)
		deleteScenario(t, ctx, db, inactiveID)
	})

	repository := scenario.NewPgRepository(db)

	t.Run("gets an active scenario", func(t *testing.T) {
		got, err := repository.GetActiveByID(ctx, activeID)

		require.NoError(t, err)
		require.Equal(t, activeID, got.ID)
		require.Equal(t, scenario.RoleBuyer, got.Role)
		require.Equal(t, testfixture.StartNodeID, got.StartNodeID)
		require.Len(t, got.Nodes, 3)
		require.Len(t, got.Endings, 2)
	})

	t.Run("loads an inactive scenario but rejects it for a new attempt", func(t *testing.T) {
		got, err := repository.GetByID(ctx, inactiveID)
		require.NoError(t, err)
		require.Equal(t, inactiveID, got.ID)

		_, err = repository.GetActiveByID(ctx, inactiveID)
		require.ErrorIs(t, err, scenario.ErrScenarioInactive)
	})

	t.Run("returns not found for an unknown ID", func(t *testing.T) {
		_, err := repository.GetByID(ctx, scenario.ScenarioID("00000000-0000-0000-0000-000000000000"))

		require.ErrorIs(t, err, scenario.ErrScenarioNotFound)
	})
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

func insertScenario(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
	isActive bool,
) scenario.ScenarioID {
	t.Helper()

	fixture := testfixture.ValidScenario()
	content, err := json.Marshal(scenario.Content{
		StartNodeID: fixture.StartNodeID,
		Nodes:       fixture.Nodes,
		Endings:     fixture.Endings,
	})
	require.NoError(t, err)

	const query = `
		INSERT INTO scenario_versions (
			logical_id, version, role, title, description, is_active, content
		)
		VALUES (gen_random_uuid(), 1, $1, $2, $3, $4, $5)
		RETURNING id
	`

	var id scenario.ScenarioID
	err = db.QueryRowContext(
		ctx,
		query,
		fixture.Role,
		fixture.Title,
		fixture.Description,
		isActive,
		string(content),
	).Scan(&id)
	require.NoError(t, err)

	return id
}

func deleteScenario(t *testing.T, ctx context.Context, db *sql.DB, id scenario.ScenarioID) {
	t.Helper()

	_, err := db.ExecContext(ctx, `DELETE FROM scenario_versions WHERE id = $1`, id)
	require.NoError(t, err)
}
