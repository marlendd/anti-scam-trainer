package scenario

import "fmt"

func Validate(s Scenario) error {
	if s.StartNodeID == "" {
		return fmt.Errorf("empty StartNodeID")
	}

	seenNodes := make(map[NodeID]struct{}, len(s.Nodes))
	for _, node := range s.Nodes {
		if node.ID == "" {
			return fmt.Errorf("node ID is empty")
		}

		if _, exists := seenNodes[node.ID]; exists {
			return fmt.Errorf("node ID %q is duplicated", node.ID)
		}
		seenNodes[node.ID] = struct{}{}
	}

	seenEndings := make(map[EndingID]struct{}, len(s.Endings))
	for _, ending := range s.Endings {
		if ending.ID == "" {
			return fmt.Errorf("ending ID is empty")
		}
		if _, exists := seenEndings[ending.ID]; exists {
			return fmt.Errorf("ending ID %q is duplicated", ending.ID)
		}
		seenEndings[ending.ID] = struct{}{}
	}

	if _, ok := seenNodes[s.StartNodeID]; !ok {
		return fmt.Errorf("expected StartNode included in Nodes")
	}

	seenChoices := make(map[ChoiceID]struct{})
	for _, node := range s.Nodes {
		for _, choice := range node.Choices {
			if choice.ID == "" {
				return fmt.Errorf("choice ID is empty")
			}

			if _, exists := seenChoices[choice.ID]; exists {
				return fmt.Errorf("choice ID %q is duplicated", choice.ID)
			}
			seenChoices[choice.ID] = struct{}{}

			hasNextNode := choice.NextNodeID != ""
			hasEnding := choice.EndingID != ""
			if hasNextNode == hasEnding {
				return fmt.Errorf("choice must point to exactly one target")
			}

			if hasNextNode {
				if _, exists := seenNodes[choice.NextNodeID]; !exists {
					return fmt.Errorf(
						"choice %q refers to unknown node %q",
						choice.ID,
						choice.NextNodeID,
					)
				}
			}

			if hasEnding {
				if _, exists := seenEndings[choice.EndingID]; !exists {
					return fmt.Errorf(
						"choice %q refers to unknown ending %q",
						choice.ID,
						choice.EndingID,
					)
				}
			}

			if choice.Weight < WeightLow || choice.Weight > WeightHigh {
				return fmt.Errorf("expected weight in range 1..3, got: %d", choice.Weight)
			}

			if choice.Score != ScoreCritical && choice.Score != ScoreDangerous && choice.Score != ScoreSafe {
				return fmt.Errorf("expected score in [0,50,100], got: %d", choice.Score)
			}

			if choice.Score != ScoreSafe {
				if len(choice.RiskCategories) == 0 {
					return fmt.Errorf("empty RiskCategories when score is unsafe")
				}
			}
		}
	}

	return nil
}
