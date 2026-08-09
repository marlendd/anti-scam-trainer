package progress

import (
	"context"

	"github.com/marlendd/anti-scam-trainer/internal/evaluation"
)

type Service struct {
	repo      Repository
	evaluator *evaluation.Evaluator
}

func NewService(repo Repository, ev *evaluation.Evaluator) *Service {
	return &Service{repo: repo, evaluator: ev}
}

func (s *Service) GetAttemptResults(ctx context.Context, userID, attemptID string) (int, error) {
	answers, err := s.repo.GetAnswersByAttempt(ctx, userID, attemptID)
	if err != nil {
		return -1, err
	}

	if len(answers) == 0 {
		return -1, nil
	}

	return s.evaluator.CalculateScore(answers), nil
}

func (s *Service) GetUserRoleStats(ctx context.Context, userID string) ([]RoleStats, error) {
	return s.repo.GetUserStatsByRole(ctx, userID)
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

func (s *Service) GetUserCategoryDashboard(ctx context.Context, userID string) (DashboardData, error) {
	stats, err := s.repo.GetUserStatsByCategory(ctx, userID)
	if err != nil {
		return DashboardData{}, err
	}

	total, err := s.repo.GetUserTotalCompletedCount(ctx, userID)
	if err != nil {
		return DashboardData{}, err
	}

	return DashboardData{
		TotalCompleted: total,
		Stats:          stats,
	}, nil
}

func (s *Service) GetLeaderboard(ctx context.Context, limit, offset int) (LeaderboardResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 50 // дефолт
	}
	if offset < 0 {
		offset = 0
	}

	entries, err := s.repo.GetLeaderboard(ctx, limit, offset)
	if err != nil {
		return LeaderboardResponse{}, err
	}

	if entries == nil {
		entries = []LeaderboardEntry{}
	}

	return LeaderboardResponse{Entries: entries}, nil
}

func (s *Service) GetMyRankHistory(ctx context.Context, userID string) (RankHistoryResponse, error) {
	history, err := s.repo.GetUserRankHistory(ctx, userID, 7) // За последнюю неделю
	if err != nil {
		return RankHistoryResponse{}, err
	}

	if history == nil {
		history = []RankHistoryPoint{}
	}

	return RankHistoryResponse{History: history}, nil
}
