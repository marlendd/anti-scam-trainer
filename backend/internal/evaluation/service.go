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

func (s *Service) GetUserPuzzleProgress(ctx context.Context, userID string) (PuzzleProgress, error) {
	fragments, err := s.repo.GetUserFragments(ctx, userID)
	if err != nil {
		return PuzzleProgress{}, err
	}

	total, err := s.repo.GetTotalAvailableFragments(ctx)
	if err != nil {
		return PuzzleProgress{}, err
	}

	return PuzzleProgress{
		EarnedCount: len(fragments),
		TotalCount:  total,
		Fragments:   fragments,
	}, nil
}

// Эту функцию нужно вызывать в конце сценария
func (s *Service) GrantRewardIfEligible(ctx context.Context, userID, fragmentID string, isSuccess bool) error {
	if !isSuccess || fragmentID == "" {
		return nil
	}
	return s.repo.SaveReward(ctx, userID, fragmentID)
}
