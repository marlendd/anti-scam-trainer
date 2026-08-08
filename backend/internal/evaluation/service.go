package evaluation

import "math"

type Evaluator struct{}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

func (e *Evaluator) CalculateScore(answers []AnswerData) int {
	if len(answers) == 0 {
		return 0
	}

	var sum int
	var weightSum int
	for _, val := range answers {
		sum += int(val.Weight * val.ChoiceScore)
		weightSum += int(val.Weight)
	}

	return int(math.Round(float64(sum) / float64(weightSum)))
}
