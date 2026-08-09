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
	require.NotEmpty(t, seed.Content.Nodes[0].Messages)
	expectedMessages := append([]scenario.Message(nil), seed.Content.Nodes[0].Messages...)
	legacyContent := seed.Content
	legacyContent.Nodes = append([]scenario.Node(nil), seed.Content.Nodes...)
	legacyContent.Nodes[0].Author = expectedMessages[len(expectedMessages)-1].Author
	legacyContent.Nodes[0].Text = "legacy combined dialogue"
	legacyContent.Nodes[0].Messages = nil
	legacyContentJSON, err := json.Marshal(legacyContent)
	require.NoError(t, err)
	_, err = db.ExecContext(
		ctx,
		`UPDATE scenario_versions SET content = $1 WHERE id = $2`,
		legacyContentJSON,
		seed.ID,
	)
	require.NoError(t, err)

	_, err = scenario.ApplySeedFiles(ctx, db, "../../seeds")
	require.NoError(t, err)

	var backfilledContentJSON json.RawMessage
	err = db.QueryRowContext(
		ctx,
		`SELECT content FROM scenario_versions WHERE id = $1`,
		seed.ID,
	).Scan(&backfilledContentJSON)
	require.NoError(t, err)
	var backfilledContent scenario.Content
	require.NoError(t, json.Unmarshal(backfilledContentJSON, &backfilledContent))
	require.Equal(t, expectedMessages, backfilledContent.Nodes[0].Messages)
	require.Empty(t, backfilledContent.Nodes[0].Author)
	require.Empty(t, backfilledContent.Nodes[0].Text)

	_, err = db.ExecContext(
		ctx,
		`UPDATE scenario_versions
		 SET reward_fragment_id = NULL,
		     content = content - 'successful_ending_ids' - 'product'
		 WHERE id = $1`,
		seed.ID,
	)
	require.NoError(t, err)

	_, err = scenario.ApplySeedFiles(ctx, db, "../../seeds")
	require.NoError(t, err)

	var (
		rewardFragmentID        string
		successfulEndingIDsJSON json.RawMessage
		productJSON             json.RawMessage
	)
	err = db.QueryRowContext(
		ctx,
		`SELECT reward_fragment_id,
		        content -> 'successful_ending_ids',
		        content -> 'product'
		 FROM scenario_versions
		 WHERE id = $1`,
		seed.ID,
	).Scan(&rewardFragmentID, &successfulEndingIDsJSON, &productJSON)
	require.NoError(t, err)

	var successfulEndingIDs []string
	require.NoError(t, json.Unmarshal(successfulEndingIDsJSON, &successfulEndingIDs))
	require.Equal(t, string(seed.RewardFragmentID), rewardFragmentID)
	require.Equal(t, []string{"ending_safe", "ending_recovered"}, successfulEndingIDs)
	var product scenario.Product
	require.NoError(t, json.Unmarshal(productJSON, &product))
	require.Equal(t, seed.Content.Product, product)
}
