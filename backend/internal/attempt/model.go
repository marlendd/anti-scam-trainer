package attempt

import (
	"time"

	"github.com/marlendd/anti-scam-trainer/internal/scenario"
)

type AttemptID string
type Status string

const (
	StatusInProgress Status = "in_progress"
	StatusAborted    Status = "aborted"
	StatusCompleted  Status = "completed"
)

type Attempt struct {
	ID            AttemptID
	UserID        string
	ScenarioID    scenario.ScenarioID
	Status        Status
	CurrentNodeID *scenario.NodeID   // nullable
	EndingID      *scenario.EndingID // nullable
	Score         *int               // nullable
	StartedAt     time.Time
	UpdatedAt     time.Time
	CompletedAt   *time.Time // nullable
}
