package scenario_test

import (
	"context"
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
	}
}
