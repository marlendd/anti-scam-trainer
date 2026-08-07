package evaluation

import "time"

type AnswerData struct {
	Weight      int16 `json:"weight"`
	ChoiceScore int16 `json:"choice_score"`
}

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
