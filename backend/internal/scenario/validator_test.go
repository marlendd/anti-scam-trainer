package scenario_test

import (
	"fmt"
	"testing"

	"github.com/marlendd/anti-scam-trainer/internal/scenario"
	"github.com/marlendd/anti-scam-trainer/internal/testfixture"
	"github.com/stretchr/testify/require"
)

func TestValidate_ValidScenarioWithThreeChoicesAndBranching(t *testing.T) {
	s := testfixture.ValidScenario()

	err := scenario.Validate(s)

	if err != nil {
		t.Fatalf("expected valid scenario, got error: %v", err)
	}
}

func TestValidate_KnownRiskCategories(t *testing.T) {
	categories := []scenario.RiskCategory{
		scenario.RiskExternalMessenger,
		scenario.RiskExternalLink,
		scenario.RiskSMSCode,
		scenario.RiskOffPlatformPayment,
		scenario.RiskFakePaymentDelivery,
		scenario.RiskUrgencyPressure,
		scenario.RiskDisableProtection,
	}

	for _, category := range categories {
		t.Run(string(category), func(t *testing.T) {
			s := testfixture.ValidScenario()
			s.Nodes[0].Choices[0].Score = scenario.ScoreRisky
			s.Nodes[0].Choices[0].RiskCategories = []scenario.RiskCategory{category}

			require.NoError(t, scenario.Validate(s))
		})
	}
}

func TestValidate_InvalidScenarios(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*scenario.Scenario)
		wantErr error
	}{
		{
			name: "empty start node ID",
			mutate: func(s *scenario.Scenario) {
				s.StartNodeID = ""
			},
			wantErr: scenario.ErrEmptyStartNodeID,
		},
		{
			name: "unknown start node",
			mutate: func(s *scenario.Scenario) {
				s.StartNodeID = scenario.NodeID("unknown")
			},
			wantErr: scenario.ErrUnknownStartNode,
		},
		{
			name: "empty node ID",
			mutate: func(s *scenario.Scenario) {
				s.Nodes[0].ID = ""
			},
			wantErr: scenario.ErrEmptyNodeID,
		},
		{
			name: "empty node text",
			mutate: func(s *scenario.Scenario) {
				s.Nodes[0].Text = "  "
			},
			wantErr: scenario.ErrEmptyNodeText,
		},
		{
			name: "duplicated node ID",
			mutate: func(s *scenario.Scenario) {
				s.Nodes = append(s.Nodes, s.Nodes[0])
			},
			wantErr: scenario.ErrDuplicateNodeID,
		},
		{
			name: "empty ending ID",
			mutate: func(s *scenario.Scenario) {
				s.Endings[0].ID = ""
			},
			wantErr: scenario.ErrEmptyEndingID,
		},
		{
			name: "empty ending header",
			mutate: func(s *scenario.Scenario) {
				s.Endings[0].Header = "\t"
			},
			wantErr: scenario.ErrEmptyEndingHeader,
		},
		{
			name: "empty ending result",
			mutate: func(s *scenario.Scenario) {
				s.Endings[0].Result = "\n"
			},
			wantErr: scenario.ErrEmptyEndingResult,
		},
		{
			name: "duplicated ending ID",
			mutate: func(s *scenario.Scenario) {
				s.Endings = append(s.Endings, s.Endings[0])
			},
			wantErr: scenario.ErrDuplicateEndingID,
		},
		{
			name: "node has fewer than four choices",
			mutate: func(s *scenario.Scenario) {
				s.Nodes[0].Choices = s.Nodes[0].Choices[:3]
			},
			wantErr: scenario.ErrTooFewNodeChoices,
		},
		{
			name: "empty choice ID",
			mutate: func(s *scenario.Scenario) {
				s.Nodes[0].Choices[0].ID = ""
			},
			wantErr: scenario.ErrEmptyChoiceID,
		},
		{
			name: "empty choice text",
			mutate: func(s *scenario.Scenario) {
				s.Nodes[0].Choices[0].Text = " "
			},
			wantErr: scenario.ErrEmptyChoiceText,
		},
		{
			name: "empty choice consequence",
			mutate: func(s *scenario.Scenario) {
				s.Nodes[0].Choices[0].Consequence = "\t"
			},
			wantErr: scenario.ErrEmptyConsequence,
		},
		{
			name: "empty choice explanation",
			mutate: func(s *scenario.Scenario) {
				s.Nodes[0].Choices[0].Explanation = "\n"
			},
			wantErr: scenario.ErrEmptyExplanation,
		},
		{
			name: "duplicated choice ID",
			mutate: func(s *scenario.Scenario) {
				s.Nodes[0].Choices[1].ID = s.Nodes[0].Choices[0].ID
			},
			wantErr: scenario.ErrDuplicateChoiceID,
		},
		{
			name: "choice has no target",
			mutate: func(s *scenario.Scenario) {
				s.Nodes[0].Choices[0].NextNodeID = ""
			},
			wantErr: scenario.ErrInvalidChoiceTarget,
		},
		{
			name: "choice has two targets",
			mutate: func(s *scenario.Scenario) {
				s.Nodes[0].Choices[0].EndingID = s.Endings[0].ID
			},
			wantErr: scenario.ErrInvalidChoiceTarget,
		},
		{
			name: "choice refers to unknown node",
			mutate: func(s *scenario.Scenario) {
				s.Nodes[0].Choices[0].NextNodeID = "unknown-node"
			},
			wantErr: scenario.ErrUnknownNode,
		},
		{
			name: "choice refers to unknown ending",
			mutate: func(s *scenario.Scenario) {
				s.Nodes[2].Choices[0].EndingID = "unknown-ending"
			},
			wantErr: scenario.ErrUnknownEnding,
		},
		{
			name: "weight is too low",
			mutate: func(s *scenario.Scenario) {
				s.Nodes[0].Choices[0].Weight = scenario.Weight(0)
			},
			wantErr: scenario.ErrInvalidWeight,
		},
		{
			name: "weight is too high",
			mutate: func(s *scenario.Scenario) {
				s.Nodes[0].Choices[0].Weight = scenario.Weight(4)
			},
			wantErr: scenario.ErrInvalidWeight,
		},
		{
			name: "invalid score",
			mutate: func(s *scenario.Scenario) {
				s.Nodes[0].Choices[0].Score = scenario.ChoiceScore(25)
			},
			wantErr: scenario.ErrInvalidScore,
		},
		{
			name: "unsafe choice has no risk categories",
			mutate: func(s *scenario.Scenario) {
				s.Nodes[0].Choices[0].Score = scenario.ScoreRisky
				s.Nodes[0].Choices[0].RiskCategories = nil
			},
			wantErr: scenario.ErrMissingRiskCategory,
		},
		{
			name: "choice has unknown risk category",
			mutate: func(s *scenario.Scenario) {
				s.Nodes[0].Choices[0].Score = scenario.ScoreRisky
				s.Nodes[0].Choices[0].RiskCategories = []scenario.RiskCategory{
					"externl_link",
				}
			},
			wantErr: scenario.ErrUnknownRiskCategory,
		},
		{
			name: "unreachable node",
			mutate: func(s *scenario.Scenario) {
				choices := make([]scenario.Choice, 4)
				for i := range choices {
					choices[i] = scenario.Choice{
						ID:          scenario.ChoiceID(fmt.Sprintf("choice-disconnected-%d", i)),
						Text:        "Вариант ответа",
						Consequence: "Последствие ответа",
						Explanation: "Объяснение ответа",
						Weight:      scenario.WeightLow,
						Score:       scenario.ScoreSafe,
						NextNodeID:  testfixture.FinalNodeID,
					}
				}
				s.Nodes = append(s.Nodes, scenario.Node{
					ID:      "node-disconnected",
					Text:    "Недостижимый узел",
					Choices: choices,
				})
			},
			wantErr: scenario.ErrUnreachableNode,
		},
		{
			name: "unreachable ending",
			mutate: func(s *scenario.Scenario) {
				s.Endings = append(s.Endings, scenario.Ending{
					ID:     "unreachable-ending",
					Header: "Недостижимая концовка",
					Result: "Недостижимый результат",
				})
			},
			wantErr: scenario.ErrUnreachableEnding,
		},
		{
			name: "too few reachable endings",
			mutate: func(s *scenario.Scenario) {
				alternateNode := s.Nodes[1]
				alternateNode.ID = "node-alternate"
				alternateNode.Choices = append(
					[]scenario.Choice(nil),
					s.Nodes[1].Choices...,
				)
				for i := range alternateNode.Choices {
					alternateNode.Choices[i].ID = scenario.ChoiceID(
						fmt.Sprintf("choice-alternate-%d", i),
					)
				}
				s.Nodes = append(s.Nodes, alternateNode)
				s.Nodes[0].Choices[3].NextNodeID = alternateNode.ID

				s.Endings = s.Endings[:1]
				for i := range s.Nodes[2].Choices {
					s.Nodes[2].Choices[i].EndingID = s.Endings[0].ID
				}
			},
			wantErr: scenario.ErrTooFewReachableEndings,
		},
		{
			name: "cycle in graph",
			mutate: func(s *scenario.Scenario) {
				choice := &s.Nodes[2].Choices[0]

				choice.EndingID = ""
				choice.NextNodeID = s.Nodes[0].ID
			},
			wantErr: scenario.ErrCycleDetected,
		},
		{
			name: "path has only two choices",
			mutate: func(s *scenario.Scenario) {
				for i := range s.Nodes[0].Choices {
					s.Nodes[0].Choices[i].NextNodeID = s.Nodes[2].ID
				}
				s.Nodes = []scenario.Node{s.Nodes[0], s.Nodes[2]}
			},
			wantErr: scenario.ErrPathTooShort,
		},
		{
			name: "all choices on path lead to the same target",
			mutate: func(s *scenario.Scenario) {
				for i := range s.Nodes[0].Choices {
					s.Nodes[0].Choices[i].NextNodeID = s.Nodes[0].Choices[0].NextNodeID
				}
				for i := range s.Nodes[1].Choices {
					s.Nodes[1].Choices[i].NextNodeID = s.Nodes[1].Choices[0].NextNodeID
				}
				for i := range s.Nodes[2].Choices {
					s.Nodes[2].Choices[i].EndingID = s.Nodes[2].Choices[0].EndingID
				}
			},
			wantErr: scenario.ErrPathWithoutBranching,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := testfixture.ValidScenario()
			tc.mutate(&s)

			err := scenario.Validate(s)

			require.Error(t, err)
			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}
