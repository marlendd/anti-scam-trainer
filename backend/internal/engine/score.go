package engine

import "math"

func CalculateScore(weightedScoreSum, weightSum int) (int, error) {
	if weightSum <= 0 {
		return 0, ErrInvalidTotalWeight
	}

	score := float64(weightedScoreSum) / float64(weightSum)

	return int(math.Round(score)), nil
}
