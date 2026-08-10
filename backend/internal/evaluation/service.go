package evaluation

import (
	"errors"
	"math"
)

var ErrInvalidTotalWeight = errors.New("total weight must be positive")

type Evaluator struct{}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

func (e *Evaluator) CalculateScore(answers []AnswerData) (int, error) {
	var weightedScoreSum int
	var weightSum int
	for _, answer := range answers {
		weightedScoreSum += int(answer.Weight) * int(answer.ChoiceScore)
		weightSum += int(answer.Weight)
	}

	return e.CalculateScoreFromTotals(weightedScoreSum, weightSum)
}

func (e *Evaluator) CalculateScoreFromTotals(weightedScoreSum, weightSum int) (int, error) {
	if weightSum <= 0 {
		return 0, ErrInvalidTotalWeight
	}

	return int(math.Round(float64(weightedScoreSum) / float64(weightSum))), nil
}
