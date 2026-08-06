package engine

import (
	"fmt"

	"github.com/marlendd/anti-scam-trainer/internal/scenario"
)

type Input struct {
	Scenario      scenario.Scenario
	CurrentNodeID scenario.NodeID
	ChoiceID      scenario.ChoiceID
}

type Result struct {
	Choice      scenario.Choice
	Consequence string
	NextNodeID  scenario.NodeID
	EndingID    scenario.EndingID
	Completed   bool
}

func Transition(in Input) (Result, error) {
	currentNode, err := findNode(in.Scenario.Nodes, in.CurrentNodeID)
	if err != nil {
		return Result{}, err
	}

	currentChoice, err := findChoice(currentNode.Choices, in.ChoiceID)
	if err != nil {
		return Result{}, err
	}

	hasNextNode := currentChoice.NextNodeID != ""
	hasEnding := currentChoice.EndingID != ""
	if hasNextNode == hasEnding {
		return Result{}, fmt.Errorf("%w: choice %q", ErrInvalidChoiceTarget, currentChoice.ID)
	}

	result := Result{}
	if currentChoice.NextNodeID != "" {
		result.NextNodeID = currentChoice.NextNodeID
	} else {
		result.EndingID = currentChoice.EndingID
		result.Completed = true
	}
	result.Choice = currentChoice
	result.Consequence = currentChoice.Consequence

	return result, nil
}

// findNode - поиск узла сценария по его ID.
func findNode(nodes []scenario.Node, nodeID scenario.NodeID) (scenario.Node, error) {
	for _, node := range nodes {
		if node.ID == nodeID {
			return node, nil
		}
	}

	return scenario.Node{}, fmt.Errorf("%w: %q", ErrCurrentNodeNotFound, nodeID)
}

// findChoice - поиск варианта ответа только внутри текущего узла.
func findChoice(choices []scenario.Choice, choiceID scenario.ChoiceID) (scenario.Choice, error) {
	for _, choice := range choices {
		if choice.ID == choiceID {
			return choice, nil
		}
	}

	return scenario.Choice{}, fmt.Errorf("%w: %q", ErrChoiceNotFound, choiceID)
}
