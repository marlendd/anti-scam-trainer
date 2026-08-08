package seeds_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/marlendd/anti-scam-trainer/internal/scenario"
	"github.com/stretchr/testify/require"
)

type seedFile struct {
	ID          scenario.ScenarioID        `json:"id"`
	LogicalID   scenario.LogicalScenarioID `json:"logical_id"`
	Version     int                        `json:"version"`
	Role        scenario.Role              `json:"role"`
	Title       string                     `json:"title"`
	Description string                     `json:"description"`
	IsActive    bool                       `json:"is_active"`
	Content     scenario.Content           `json:"content"`
}

func TestScenarioSeedsAreValid(t *testing.T) {
	paths, err := filepath.Glob("*.json")
	require.NoError(t, err)
	require.NotEmpty(t, paths, "at least one scenario seed is required")

	for _, path := range paths {
		path := path
		t.Run(path, func(t *testing.T) {
			raw, err := os.ReadFile(path)
			require.NoError(t, err)

			var seed seedFile
			require.NoError(t, json.Unmarshal(raw, &seed))

			require.NotEmpty(t, seed.ID)
			require.NotEmpty(t, seed.LogicalID)
			require.Positive(t, seed.Version)
			require.Contains(t, []scenario.Role{scenario.RoleBuyer, scenario.RoleSeller}, seed.Role)
			require.NotEmpty(t, seed.Title)
			require.NotEmpty(t, seed.Description)

			s := scenario.Scenario{
				ID:          seed.ID,
				LogicalID:   seed.LogicalID,
				Version:     seed.Version,
				Role:        seed.Role,
				Title:       seed.Title,
				Description: seed.Description,
				StartNodeID: seed.Content.StartNodeID,
				Nodes:       seed.Content.Nodes,
				Endings:     seed.Content.Endings,
			}

			require.NoError(t, scenario.Validate(s))
		})
	}
}
