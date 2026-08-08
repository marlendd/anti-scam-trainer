package progress

import "time"

type RoleStats struct {
	Role            string `json:"role"`
	CompletedCount  int64  `json:"completed_count"`
	InProgressCount int64  `json:"in_progress_count"`
	TotalStarted    int64  `json:"total_started"`
}

type PuzzleFragment struct {
	FragmentID string    `json:"fragment_id"`
	EarnedAt   time.Time `json:"earned_at"`
}

type PuzzleProgress struct {
	EarnedCount int              `json:"earned_count"`
	TotalCount  int              `json:"total_count"`
	Fragments   []PuzzleFragment `json:"fragments"`
}

type CategoryStat struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

type DashboardData struct {
	TotalCompleted int            `json:"total_completed"`
	Stats          []CategoryStat `json:"stats"`
}

type LeaderboardEntry struct {
	Rank       int    `json:"rank"`
	Player     string `json:"player"`
	Fragments  int    `json:"fragments"`
	Score      int    `json:"score"`
	RankChange *int   `json:"rank_change"`
}

type LeaderboardResponse struct {
	Entries []LeaderboardEntry `json:"entries"`
}

type RankHistoryPoint struct {
	Date time.Time `json:"date"`
	Rank int       `json:"rank"`
}

type RankHistoryResponse struct {
	History []RankHistoryPoint `json:"history"`
}
