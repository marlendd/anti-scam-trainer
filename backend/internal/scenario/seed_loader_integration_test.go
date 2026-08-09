package scenario_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/marlendd/anti-scam-trainer/internal/scenario"
	"github.com/stretchr/testify/require"
)

func TestApplySeedFilesIsIdempotent_Integration(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()

	firstCount, err := scenario.ApplySeedFiles(ctx, db, "../../seeds")
	require.NoError(t, err)
	require.Equal(t, 4, firstCount)

	secondCount, err := scenario.ApplySeedFiles(ctx, db, "../../seeds")
	require.NoError(t, err)
	require.Equal(t, 4, secondCount)

	seeds, err := scenario.LoadSeedFiles("../../seeds")
	require.NoError(t, err)

	for _, seed := range seeds {
		var count int
		err := db.QueryRowContext(
			ctx,
			`SELECT count(*) FROM scenario_versions WHERE id = $1`,
			seed.ID,
		).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count)

		var rewardFragmentID string
		err = db.QueryRowContext(
			ctx,
			`SELECT reward_fragment_id FROM scenario_versions WHERE id = $1`,
			seed.ID,
		).Scan(&rewardFragmentID)
		require.NoError(t, err)
		require.Equal(t, string(seed.RewardFragmentID), rewardFragmentID)
	}

	seed := seeds[0]
	_, err = db.ExecContext(
		ctx,
		`UPDATE scenario_versions
		 SET reward_fragment_id = NULL,
		     content = content - 'successful_ending_ids'
		 WHERE id = $1`,
		seed.ID,
	)
	require.NoError(t, err)

	_, err = scenario.ApplySeedFiles(ctx, db, "../../seeds")
	require.NoError(t, err)

	var (
		rewardFragmentID        string
		successfulEndingIDsJSON json.RawMessage
	)
	err = db.QueryRowContext(
		ctx,
		`SELECT reward_fragment_id,
		        content -> 'successful_ending_ids'
		 FROM scenario_versions
		 WHERE id = $1`,
		seed.ID,
	).Scan(&rewardFragmentID, &successfulEndingIDsJSON)
	require.NoError(t, err)

	var successfulEndingIDs []string
	require.NoError(t, json.Unmarshal(successfulEndingIDsJSON, &successfulEndingIDs))
	require.Equal(t, string(seed.RewardFragmentID), rewardFragmentID)
	require.Equal(t, []string{"ending_safe", "ending_recovered"}, successfulEndingIDs)
}
