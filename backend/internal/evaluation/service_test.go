package evaluation_test

import (
	"testing"

	"github.com/marlendd/anti-scam-trainer/internal/evaluation"
	"github.com/stretchr/testify/require"
)

func TestEvaluatorCalculateScore(t *testing.T) {
	t.Parallel()

	evaluator := evaluation.NewEvaluator()
	score, err := evaluator.CalculateScore([]evaluation.AnswerData{
		{Weight: 2, ChoiceScore: 50},
		{Weight: 1, ChoiceScore: 100},
	})

	require.NoError(t, err)
	require.Equal(t, 67, score)
}

func TestEvaluatorCalculateScoreFromTotals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		weightedScoreSum int
		weightSum        int
		wantScore        int
	}{
		{name: "exact weighted average", weightedScoreSum: 300, weightSum: 4, wantScore: 75},
		{name: "rounds fractional score down", weightedScoreSum: 100, weightSum: 3, wantScore: 33},
		{name: "rounds fractional score up", weightedScoreSum: 250, weightSum: 6, wantScore: 42},
		{name: "zero score", weightedScoreSum: 0, weightSum: 3, wantScore: 0},
		{name: "maximum score", weightedScoreSum: 300, weightSum: 3, wantScore: 100},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			score, err := evaluation.NewEvaluator().CalculateScoreFromTotals(
				test.weightedScoreSum,
				test.weightSum,
			)

			require.NoError(t, err)
			require.Equal(t, test.wantScore, score)
		})
	}
}

func TestEvaluatorRejectsInvalidTotalWeight(t *testing.T) {
	t.Parallel()

	for _, weightSum := range []int{0, -1} {
		score, err := evaluation.NewEvaluator().CalculateScoreFromTotals(100, weightSum)

		require.ErrorIs(t, err, evaluation.ErrInvalidTotalWeight)
		require.Zero(t, score)
	}
}
