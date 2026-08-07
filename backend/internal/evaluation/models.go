package evaluation

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
