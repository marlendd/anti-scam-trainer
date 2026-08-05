package scenario

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func validScenario() Scenario {
	endings := []Ending{
		{
			ID:     "ending-safe",
			Header: "Сделка остановлена вовремя",
			Result: "Пользователь не передал и сохранил защиту платформы.",
		},
		{
			ID:     "ending-scam",
			Header: "Данные переданы мошеннику",
			Result: "Пользователь сообщил секретный код вне защищённого канала.",
		},
	}

	choices := []Choice{
		{
			ID:         "choice-safe",
			Weight:     WeightLow,
			Score:      ScoreSafe,
			NextNodeID: "node-next",
		},
		{
			ID:             "choice-risky",
			Weight:         WeightMedium,
			Score:          ScoreDangerous,
			RiskCategories: []RiskCategory{RiskExternalMessenger},
			NextNodeID:     "node-next",
		},
		{
			ID:       "choice-safe-finish",
			Weight:   WeightLow,
			Score:    ScoreSafe,
			EndingID: "ending-safe",
		},
		{
			ID:             "choice-danger-finish",
			Weight:         WeightHigh,
			Score:          ScoreCritical,
			RiskCategories: []RiskCategory{RiskSMSCode, RiskUrgencyPressure},
			EndingID:       "ending-scam",
		},
	}
	nodes := []Node{
		{
			ID: "node-start",
			Choices: []Choice{
				choices[0], choices[1],
			},
		},
		{
			ID: "node-next",
			Choices: []Choice{
				choices[2], choices[3],
			},
		},
	}

	s := Scenario{
		ID:          "1",
		LogicalID:   "valid-scenario",
		StartNodeID: nodes[0].ID,
		Nodes:       nodes,
		Endings:     endings,
	}

	return s
}

func TestValidate_ValidScenario(t *testing.T) {
	s := validScenario()

	err := Validate(s)

	if err != nil {
		t.Fatalf("expected valid scenario, got error: %v", err)
	}
}

func TestValidate_InvalidScenarios(t *testing.T) {
	cases := []struct {
		name      string
		mutate    func(*Scenario)
		errorPart string
	}{
		{
			name: "empty start node ID",
			mutate: func(s *Scenario) {
				s.StartNodeID = ""
			},
			errorPart: "empty StartNodeID",
		},
		{
			name: "unknown start node",
			mutate: func(s *Scenario) {
				s.StartNodeID = NodeID("unknown")
			},
			errorPart: "StartNode",
		},
		{
			name: "empty node ID",
			mutate: func(s *Scenario) {
				s.Nodes[0].ID = ""
			},
			errorPart: "node ID is empty",
		},
		{
			name: "duplicated node ID",
			mutate: func(s *Scenario) {
				s.Nodes = append(s.Nodes, s.Nodes[0])
			},
			errorPart: "node ID",
		},
		{
			name: "empty ending ID",
			mutate: func(s *Scenario) {
				s.Endings[0].ID = ""
			},
			errorPart: "ending ID is empty",
		},
		{
			name: "duplicated ending ID",
			mutate: func(s *Scenario) {
				s.Endings = append(s.Endings, s.Endings[0])
			},
			errorPart: "ending ID",
		},
		{
			name: "empty choice ID",
			mutate: func(s *Scenario) {
				s.Nodes[0].Choices[0].ID = ""
			},
			errorPart: "choice ID is empty",
		},
		{
			name: "duplicated choice ID",
			mutate: func(s *Scenario) {
				s.Nodes[0].Choices[1].ID = s.Nodes[0].Choices[0].ID
			},
			errorPart: "choice ID",
		},
		{
			name: "choice has no target",
			mutate: func(s *Scenario) {
				s.Nodes[0].Choices[0].NextNodeID = ""
			},
			errorPart: "exactly one target",
		},
		{
			name: "choice has two targets",
			mutate: func(s *Scenario) {
				s.Nodes[0].Choices[0].EndingID = s.Endings[0].ID
			},
			errorPart: "exactly one target",
		},
		{
			name: "choice refers to unknown node",
			mutate: func(s *Scenario) {
				s.Nodes[0].Choices[0].NextNodeID = "unknown-node"
			},
			errorPart: "unknown node",
		},
		{
			name: "choice refers to unknown ending",
			mutate: func(s *Scenario) {
				s.Nodes[1].Choices[0].EndingID = "unknown-ending"
			},
			errorPart: "unknown ending",
		},
		{
			name: "weight is too low",
			mutate: func(s *Scenario) {
				s.Nodes[0].Choices[0].Weight = Weight(0)
			},
			errorPart: "weight",
		},
		{
			name: "weight is too high",
			mutate: func(s *Scenario) {
				s.Nodes[0].Choices[0].Weight = Weight(4)
			},
			errorPart: "weight",
		},
		{
			name: "invalid score",
			mutate: func(s *Scenario) {
				s.Nodes[0].Choices[0].Score = ChoiceScore(25)
			},
			errorPart: "score",
		},
		{
			name: "unsafe choice has no risk categories",
			mutate: func(s *Scenario) {
				s.Nodes[0].Choices[1].RiskCategories = nil
			},
			errorPart: "RiskCategories",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := validScenario()
			tc.mutate(&s)

			err := Validate(s)

			require.Error(t, err)
			require.ErrorContains(t, err, tc.errorPart)
		})
	}
}
