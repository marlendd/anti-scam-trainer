package scenario

import "fmt"

func Validate(s Scenario) error {
	if s.StartNodeID == "" {
		return ErrEmptyStartNodeID
	}

	nodesByID := make(map[NodeID]Node, len(s.Nodes))
	for _, node := range s.Nodes {
		if node.ID == "" {
			return ErrEmptyNodeID
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
		if _, exists := seenEndings[ending.ID]; exists {
			return fmt.Errorf("%w: %q", ErrDuplicateEndingID, ending.ID)
		}
		seenEndings[ending.ID] = struct{}{}
	}

	if _, ok := nodesByID[s.StartNodeID]; !ok {
		return fmt.Errorf("%w: %q", ErrUnknownStartNode, s.StartNodeID)
	}

	seenChoices := make(map[ChoiceID]struct{})
	for _, node := range s.Nodes {
		for _, choice := range node.Choices {
			if choice.ID == "" {
				return ErrEmptyChoiceID
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

			if choice.Score != ScoreSafe {
				if len(choice.RiskCategories) == 0 {
					return fmt.Errorf("%w: choice %q", ErrMissingRiskCategory, choice.ID)
				}
			}
		}
	}

	return validateGraph(s, nodesByID, seenEndings)
}

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
