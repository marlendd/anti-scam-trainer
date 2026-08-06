package testfixture

import (
	"fmt"

	"github.com/marlendd/anti-scam-trainer/internal/scenario"
)

const (
	StartNodeID   scenario.NodeID   = "node-start"
	MiddleNodeID  scenario.NodeID   = "node-middle"
	FinalNodeID   scenario.NodeID   = "node-final"
	StartChoiceID scenario.ChoiceID = "start-choice-1"
	FinalChoiceID scenario.ChoiceID = "final-choice-1"
	SafeEndingID  scenario.EndingID = "ending-safe"
	RiskyEndingID scenario.EndingID = "ending-risk"
)

// ValidScenario - валидный сценарий для тестов backend-модулей.
func ValidScenario() scenario.Scenario {
	return scenario.Scenario{
		ID:          "scenario-v1",
		LogicalID:   "scenario",
		Version:     1,
		Role:        scenario.RoleBuyer,
		Title:       "Тестовый сценарий",
		Description: "Валидный граф для тестов backend-модулей.",
		StartNodeID: StartNodeID,
		Nodes: []scenario.Node{
			{
				ID:      StartNodeID,
				Author:  "seller",
				Text:    "Первый узел",
				Choices: choicesToNode("start", MiddleNodeID),
			},
			{
				ID:      MiddleNodeID,
				Author:  "seller",
				Text:    "Второй узел",
				Choices: choicesToNode("middle", FinalNodeID),
			},
			{
				ID:     FinalNodeID,
				Author: "seller",
				Text:   "Третий узел",
				Choices: choicesToEndings(
					"final",
					SafeEndingID,
					RiskyEndingID,
				),
			},
		},
		Endings: []scenario.Ending{
			{
				ID:     SafeEndingID,
				Header: "Безопасная концовка",
				Result: "Пользователь избежал риска.",
			},
			{
				ID:     RiskyEndingID,
				Header: "Рискованная концовка",
				Result: "Пользователь столкнулся с риском.",
			},
		},
	}
}

func choicesToNode(prefix string, nextNodeID scenario.NodeID) []scenario.Choice {
	choices := make([]scenario.Choice, 0, 4)
	for i := 1; i <= 4; i++ {
		choices = append(choices, scenario.Choice{
			ID:          scenario.ChoiceID(fmt.Sprintf("%s-choice-%d", prefix, i)),
			Text:        fmt.Sprintf("Вариант %d", i),
			Consequence: fmt.Sprintf("Последствие %d", i),
			Explanation: fmt.Sprintf("Объяснение %d", i),
			Weight:      scenario.WeightLow,
			Score:       scenario.ScoreSafe,
			NextNodeID:  nextNodeID,
		})
	}

	return choices
}

func choicesToEndings(
	prefix string,
	firstEndingID scenario.EndingID,
	secondEndingID scenario.EndingID,
) []scenario.Choice {
	choices := make([]scenario.Choice, 0, 4)
	for i := 1; i <= 4; i++ {
		endingID := firstEndingID
		if i > 2 {
			endingID = secondEndingID
		}

		choices = append(choices, scenario.Choice{
			ID:          scenario.ChoiceID(fmt.Sprintf("%s-choice-%d", prefix, i)),
			Text:        fmt.Sprintf("Вариант %d", i),
			Consequence: fmt.Sprintf("Последствие %d", i),
			Explanation: fmt.Sprintf("Объяснение %d", i),
			Weight:      scenario.WeightLow,
			Score:       scenario.ScoreSafe,
			EndingID:    endingID,
		})
	}

	return choices
}
