package engine_test

import (
	"testing"

	"github.com/marlendd/anti-scam-trainer/internal/engine"
	"github.com/stretchr/testify/require"
)

func TestCalculateScore(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name             string
		weightedScoreSum int
		weightSum        int
		expectedScore    int
	}{
		{
			name:             "returns exact weighted average",
			weightedScoreSum: 300,
			weightSum:        4,
			expectedScore:    75,
		},
		{
			name:             "rounds fractional score down",
			weightedScoreSum: 100,
			weightSum:        3,
			expectedScore:    33,
		},
		{
			name:             "rounds fractional score up",
			weightedScoreSum: 250,
			weightSum:        6,
			expectedScore:    42,
		},
		{
			name:             "returns zero score",
			weightedScoreSum: 0,
			weightSum:        3,
			expectedScore:    0,
		},
		{
			name:             "returns maximum score",
			weightedScoreSum: 300,
			weightSum:        3,
			expectedScore:    100,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			score, err := engine.CalculateScore(tc.weightedScoreSum, tc.weightSum)

			require.NoError(t, err)
			require.Equal(t, tc.expectedScore, score)
		})
	}
}

func TestCalculateScoreRejectsInvalidTotalWeight(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		weightSum int
	}{
		{name: "zero total weight", weightSum: 0},
		{name: "negative total weight", weightSum: -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			score, err := engine.CalculateScore(100, tc.weightSum)

			require.ErrorIs(t, err, engine.ErrInvalidTotalWeight)
			require.Zero(t, score)
		})
	}
}
