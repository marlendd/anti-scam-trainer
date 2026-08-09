package scenario

import (
	"fmt"
	"strings"
)

// Validate - валидатор данных сценария, проверяющий, удовлетворяет ли сценарий следующим условиям:
// 1) ID старта, узлов, концовок и вариантов заполнены и уникальны;
// 2) каждый узел содержит хотя бы четыре варианта ответа;
// 3) каждый вариант имеет ровно одну цель, допустимые вес и оценку, а небезопасный вариант также имеет категорию риска.
func Validate(s Scenario) error {
	if s.StartNodeID == "" {
		return ErrEmptyStartNodeID
	}
	if s.RewardFragmentID != "" && strings.TrimSpace(string(s.RewardFragmentID)) == "" {
		return ErrInvalidRewardFragment
	}

	nodesByID := make(map[NodeID]Node, len(s.Nodes))
	for _, node := range s.Nodes {
		if node.ID == "" {
			return ErrEmptyNodeID
		}
		if strings.TrimSpace(node.Text) == "" {
			return fmt.Errorf("%w: node %q", ErrEmptyNodeText, node.ID)
		}

		if _, exists := nodesByID[node.ID]; exists {
			return fmt.Errorf("%w: %q", ErrDuplicateNodeID, node.ID)
		}
		nodesByID[node.ID] = node
	}

	seenEndings := make(map[EndingID]struct{}, len(s.Endings))
	for _, ending := range s.Endings {
		if ending.ID == "" {
			return ErrEmptyEndingID
		}
		if strings.TrimSpace(ending.Header) == "" {
			return fmt.Errorf("%w: ending %q", ErrEmptyEndingHeader, ending.ID)
		}
		if strings.TrimSpace(ending.Result) == "" {
			return fmt.Errorf("%w: ending %q", ErrEmptyEndingResult, ending.ID)
		}
		if _, exists := seenEndings[ending.ID]; exists {
			return fmt.Errorf("%w: %q", ErrDuplicateEndingID, ending.ID)
		}
		seenEndings[ending.ID] = struct{}{}
	}

	seenSuccessfulEndings := make(map[EndingID]struct{}, len(s.SuccessfulEndingIDs))
	for _, endingID := range s.SuccessfulEndingIDs {
		if _, exists := seenEndings[endingID]; !exists {
			return fmt.Errorf("%w: %q", ErrUnknownSuccessfulEnding, endingID)
		}
		if _, exists := seenSuccessfulEndings[endingID]; exists {
			return fmt.Errorf("%w: %q", ErrDuplicateSuccessfulEnding, endingID)
		}
		seenSuccessfulEndings[endingID] = struct{}{}
	}
	if s.RewardFragmentID != "" && len(seenSuccessfulEndings) == 0 {
		return ErrMissingSuccessfulEnding
	}

	if _, ok := nodesByID[s.StartNodeID]; !ok {
		return fmt.Errorf("%w: %q", ErrUnknownStartNode, s.StartNodeID)
	}

	seenChoices := make(map[ChoiceID]struct{})
	for _, node := range s.Nodes {
		if len(node.Choices) < 4 {
			return fmt.Errorf(
				"%w: node %q has %d choices",
				ErrTooFewNodeChoices,
				node.ID,
				len(node.Choices),
			)
		}

		for _, choice := range node.Choices {
			if choice.ID == "" {
				return ErrEmptyChoiceID
			}
			if strings.TrimSpace(choice.Text) == "" {
				return fmt.Errorf("%w: choice %q", ErrEmptyChoiceText, choice.ID)
			}
			if strings.TrimSpace(choice.Consequence) == "" {
				return fmt.Errorf("%w: choice %q", ErrEmptyConsequence, choice.ID)
			}
			if strings.TrimSpace(choice.Explanation) == "" {
				return fmt.Errorf("%w: choice %q", ErrEmptyExplanation, choice.ID)
			}

			if _, exists := seenChoices[choice.ID]; exists {
				return fmt.Errorf("%w: %q", ErrDuplicateChoiceID, choice.ID)
			}
			seenChoices[choice.ID] = struct{}{}

			hasNextNode := choice.NextNodeID != ""
			hasEnding := choice.EndingID != ""
			if hasNextNode == hasEnding {
				return fmt.Errorf("%w: choice %q", ErrInvalidChoiceTarget, choice.ID)
			}

			if hasNextNode {
				if _, exists := nodesByID[choice.NextNodeID]; !exists {
					return fmt.Errorf(
						"%w: choice %q refers to %q",
						ErrUnknownNode,
						choice.ID,
						choice.NextNodeID,
					)
				}
			}

			if hasEnding {
				if _, exists := seenEndings[choice.EndingID]; !exists {
					return fmt.Errorf(
						"%w: choice %q refers to %q",
						ErrUnknownEnding,
						choice.ID,
						choice.EndingID,
					)
				}
			}

			if choice.Weight < WeightLow || choice.Weight > WeightHigh {
				return fmt.Errorf("%w: got %d", ErrInvalidWeight, choice.Weight)
			}

			if choice.Score != ScoreCritical && choice.Score != ScoreRisky && choice.Score != ScoreSafe {
				return fmt.Errorf("%w: got %d", ErrInvalidScore, choice.Score)
			}

			if choice.Score != ScoreSafe && len(choice.RiskCategories) == 0 {
				return fmt.Errorf("%w: choice %q", ErrMissingRiskCategory, choice.ID)
			}

			for _, category := range choice.RiskCategories {
				if !isKnownRiskCategory(category) {
					return fmt.Errorf(
						"%w: choice %q has %q",
						ErrUnknownRiskCategory,
						choice.ID,
						category,
					)
				}
			}
		}
	}

	return validateGraph(s, nodesByID, seenEndings)
}

// isKnownRiskCategory - проверяет принадлежность категории риска закрытому набору значений.
func isKnownRiskCategory(category RiskCategory) bool {
	switch category {
	case RiskExternalMessenger,
		RiskExternalLink,
		RiskSMSCode,
		RiskOffPlatformPayment,
		RiskFakePaymentDelivery,
		RiskUrgencyPressure,
		RiskDisableProtection,
		RiskUnverifiedEvidence:
		return true
	default:
		return false
	}
}

// validateGraph - валидатор графа сценария, проверяющий, удовлетворяет ли граф следующим условиям:
// 1) в графе нет циклов, то есть не существует пути из узла в самого себя;
// 2) все узлы и концовки графа достижимы из стартового узла;
// 3) граф имеет хотя бы две достижимые концовки.
func validateGraph(
	s Scenario,
	nodesByID map[NodeID]Node,
	endingsByID map[EndingID]struct{},
) error {
	visited := make(map[NodeID]struct{})
	visiting := make(map[NodeID]struct{})
	reachedEndings := make(map[EndingID]struct{})

	var dfs func(NodeID) error

	dfs = func(nodeID NodeID) error {
		if _, exists := visiting[nodeID]; exists {
			return fmt.Errorf("%w: %q", ErrCycleDetected, nodeID)
		}

		if _, exists := visited[nodeID]; exists {
			return nil
		}

		node := nodesByID[nodeID]
		visiting[nodeID] = struct{}{}

		for _, choice := range node.Choices {
			if choice.EndingID != "" {
				reachedEndings[choice.EndingID] = struct{}{}
				continue
			}

			if err := dfs(choice.NextNodeID); err != nil {
				return err
			}
		}

		delete(visiting, nodeID)
		visited[nodeID] = struct{}{}

		return nil
	}
	if err := dfs(s.StartNodeID); err != nil {
		return err
	}
	if err := validatePaths(s.StartNodeID, nodesByID); err != nil {
		return err
	}

	for nodeID := range nodesByID {
		if _, exists := visited[nodeID]; !exists {
			return fmt.Errorf("%w: %q", ErrUnreachableNode, nodeID)
		}
	}

	for endingID := range endingsByID {
		if _, exists := reachedEndings[endingID]; !exists {
			return fmt.Errorf("%w: %q", ErrUnreachableEnding, endingID)
		}
	}

	if len(reachedEndings) < 2 {
		return fmt.Errorf("%w: got %d", ErrTooFewReachableEndings, len(reachedEndings))
	}

	return nil
}

// validatePaths - валидатор путей сценария, проверяющий, удовлетворяет ли каждый путь следующим условиям:
// 1) путь от стартового узла до концовки содержит хотя бы три выбора;
// 2) на пути есть хотя бы один узел, варианты которого ведут к минимум двум разным целям.
func validatePaths(startNodeID NodeID, nodesByID map[NodeID]Node) error {
	var walk func(NodeID, int, bool) error

	walk = func(nodeID NodeID, choiceCount int, pathHasBranching bool) error {
		node := nodesByID[nodeID]
		pathHasBranching = pathHasBranching || choicesHaveDifferentTargets(node.Choices)

		for _, choice := range node.Choices {
			pathChoiceCount := choiceCount + 1
			if choice.EndingID != "" {
				if pathChoiceCount < 3 {
					return fmt.Errorf(
						"%w: ending %q reached after %d choices",
						ErrPathTooShort,
						choice.EndingID,
						pathChoiceCount,
					)
				}
				if !pathHasBranching {
					return fmt.Errorf(
						"%w: ending %q",
						ErrPathWithoutBranching,
						choice.EndingID,
					)
				}

				continue
			}

			if err := walk(choice.NextNodeID, pathChoiceCount, pathHasBranching); err != nil {
				return err
			}
		}

		return nil
	}

	return walk(startNodeID, 0, false)
}

// choicesHaveDifferentTargets - проверка настоящего ветвления узла:
// 1) варианты узла ведут к минимум двум разным узлам или концовкам.
func choicesHaveDifferentTargets(choices []Choice) bool {
	type target struct {
		nodeID   NodeID
		endingID EndingID
	}

	targets := make(map[target]struct{}, len(choices))
	for _, choice := range choices {
		targets[target{
			nodeID:   choice.NextNodeID,
			endingID: choice.EndingID,
		}] = struct{}{}

		if len(targets) >= 2 {
			return true
		}
	}

	return false
}
