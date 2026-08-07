package evaluation

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAttemptResults(ctx context.Context, attemptID string) (int, error) {
	dbRes, err := s.repo.GetAnswersByAttempt(ctx, attemptID)
	if err != nil {
		return -1, err
	}

	if len(dbRes) == 0 {
		return 0, nil
	}

	var sum int
	var weightSum int
	for _, val := range dbRes {
		sum += int(val.Weight * val.ChoiceScore)
		weightSum += int(val.Weight)
	}

	res := sum / weightSum

	return res, nil
}

func (s *Service) GetGlobalStats(ctx context.Context) ([]RoleStats, error) {
	return s.repo.GetStatsByRole(ctx)
}
