package scenario_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/marlendd/anti-scam-trainer/internal/scenario"
	"github.com/marlendd/anti-scam-trainer/internal/testfixture"
)

func TestPgRepository_ListActiveByRole_Integration(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()
	repository := scenario.NewPgRepository(db)

	userID := insertCatalogUser(t, ctx, db)
	completedLogicalID := newUUID(t, ctx, db)
	oldCompletedVersionID := insertCatalogScenario(
		t,
		ctx,
		db,
		completedLogicalID,
		1,
		false,
	)
	activeCompletedVersionID := insertCatalogScenario(
		t,
		ctx,
		db,
		completedLogicalID,
		2,
		true,
	)
	inProgressVersionID := insertCatalogScenario(
		t,
		ctx,
		db,
		newUUID(t, ctx, db),
		1,
		true,
	)
	notStartedVersionID := insertCatalogScenario(
		t,
		ctx,
		db,
		newUUID(t, ctx, db),
		1,
		true,
	)

	insertCompletedCatalogAttempt(t, ctx, db, userID, oldCompletedVersionID, 70)
	insertInProgressCatalogAttempt(t, ctx, db, userID, inProgressVersionID)

	t.Cleanup(func() {
		_, err := db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, userID)
		require.NoError(t, err)
		for _, scenarioID := range []scenario.ScenarioID{
			oldCompletedVersionID,
			activeCompletedVersionID,
			inProgressVersionID,
			notStartedVersionID,
		} {
			_, err := db.ExecContext(ctx, `DELETE FROM scenario_versions WHERE id = $1`, scenarioID)
			require.NoError(t, err)
		}
	})

	items, err := repository.ListActiveByRole(ctx, userID, scenario.RoleBuyer)
	require.NoError(t, err)

	itemsByID := make(map[scenario.ScenarioID]scenario.CatalogItem, len(items))
	for _, item := range items {
		itemsByID[item.ID] = item
		require.Equal(t, scenario.RoleBuyer, item.Role)
	}

	completed := itemsByID[activeCompletedVersionID]
	require.Equal(t, scenario.ProgressCompleted, completed.Status)
	require.Equal(t, 2, completed.Version)
	require.Equal(t, testfixture.ValidScenario().Product, completed.Product)
	_, oldVersionExposed := itemsByID[oldCompletedVersionID]
	require.False(t, oldVersionExposed)

	inProgress := itemsByID[inProgressVersionID]
	require.Equal(t, scenario.ProgressInProgress, inProgress.Status)

	notStarted := itemsByID[notStartedVersionID]
	require.Equal(t, scenario.ProgressNotStarted, notStarted.Status)
}

func insertCatalogUser(t *testing.T, ctx context.Context, db *sql.DB) string {
	t.Helper()

	var userID string
	err := db.QueryRowContext(
		ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, 'hash') RETURNING id`,
		fmt.Sprintf("catalog-%d@example.com", time.Now().UnixNano()),
	).Scan(&userID)
	require.NoError(t, err)

	return userID
}

func newUUID(t *testing.T, ctx context.Context, db *sql.DB) string {
	t.Helper()

	var id string
	require.NoError(t, db.QueryRowContext(ctx, `SELECT gen_random_uuid()`).Scan(&id))
	return id
}

func insertCatalogScenario(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
	logicalID string,
	version int,
	isActive bool,
) scenario.ScenarioID {
	t.Helper()

	fixture := testfixture.ValidScenario()
	content, err := json.Marshal(scenario.Content{
		Product:     fixture.Product,
		StartNodeID: fixture.StartNodeID,
		Nodes:       fixture.Nodes,
		Endings:     fixture.Endings,
	})
	require.NoError(t, err)

	var scenarioID scenario.ScenarioID
	err = db.QueryRowContext(ctx, `
		INSERT INTO scenario_versions (
			logical_id, version, role, title, description, is_active, content
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb)
		RETURNING id
	`,
		logicalID,
		version,
		fixture.Role,
		fmt.Sprintf("Catalog scenario %s v%d", logicalID, version),
		fixture.Description,
		isActive,
		content,
	).Scan(&scenarioID)
	require.NoError(t, err)

	return scenarioID
}

func insertCompletedCatalogAttempt(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
	userID string,
	scenarioID scenario.ScenarioID,
	score int,
) {
	t.Helper()

	_, err := db.ExecContext(ctx, `
		INSERT INTO attempts (
			user_id, scenario_id, status, current_node_id,
			ending_id, score, completed_at
		)
		VALUES ($1, $2, 'completed', NULL, $3, $4, now())
	`, userID, scenarioID, testfixture.SafeEndingID, score)
	require.NoError(t, err)
}

func insertInProgressCatalogAttempt(
	t *testing.T,
	ctx context.Context,
	db *sql.DB,
	userID string,
	scenarioID scenario.ScenarioID,
) {
	t.Helper()

	_, err := db.ExecContext(ctx, `
		INSERT INTO attempts (user_id, scenario_id, status, current_node_id)
		VALUES ($1, $2, 'in_progress', $3)
	`, userID, scenarioID, testfixture.StartNodeID)
	require.NoError(t, err)
}
