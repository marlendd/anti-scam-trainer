package engine_test

import (
	"testing"

	"github.com/marlendd/anti-scam-trainer/internal/engine"
	"github.com/marlendd/anti-scam-trainer/internal/scenario"
	"github.com/marlendd/anti-scam-trainer/internal/testfixture"
	"github.com/stretchr/testify/require"
)

func getValidInput() engine.Input {
	return engine.Input{
		Scenario:      testfixture.ValidScenario(),
		CurrentNodeID: testfixture.StartNodeID,
		ChoiceID:      testfixture.StartChoiceID,
	}
}

func TestTransition_ToNextNode(t *testing.T) {
	in := getValidInput()

	require.NoError(
		t,
		scenario.Validate(in.Scenario),
		"getValidInput returned invalid scenario",
	)

	expectedChoice := in.Scenario.Nodes[0].Choices[0]

	result, err := engine.Transition(in)

	require.NoError(t, err)
	require.Equal(t, testfixture.MiddleNodeID, result.NextNodeID)
	require.Empty(t, result.EndingID)
	require.False(t, result.Completed)

	require.Equal(t, expectedChoice, result.Choice)
	require.Equal(t, expectedChoice.Consequence, result.Consequence)
}

func TestTransition_ToEnding(t *testing.T) {
	in := getValidInput()
	in.CurrentNodeID = testfixture.FinalNodeID
	in.ChoiceID = testfixture.FinalChoiceID

	require.NoError(
		t,
		scenario.Validate(in.Scenario),
		"getValidInput returned invalid scenario",
	)

	expectedChoice := in.Scenario.Nodes[2].Choices[0]

	result, err := engine.Transition(in)

	require.NoError(t, err)
	require.Empty(t, result.NextNodeID)
	require.Equal(t, testfixture.SafeEndingID, result.EndingID)
	require.True(t, result.Completed)

	require.Equal(t, expectedChoice, result.Choice)
	require.Equal(t, expectedChoice.Consequence, result.Consequence)
}

func TestTransition_CurrentNodeNotFound(t *testing.T) {
	in := getValidInput()
	in.CurrentNodeID = "unknown-node"

	result, err := engine.Transition(in)

	require.ErrorIs(t, err, engine.ErrCurrentNodeNotFound)
	require.Empty(t, result)
}

func TestTransition_ChoiceNotFoundInCurrentNode(t *testing.T) {
	in := getValidInput()
	in.ChoiceID = testfixture.FinalChoiceID

	result, err := engine.Transition(in)

	require.ErrorIs(t, err, engine.ErrChoiceNotFound)
	require.Empty(t, result)
}

func TestTransition_InvalidChoiceTarget(t *testing.T) {
	tests := []struct {
		name       string
		nextNodeID scenario.NodeID
		endingID   scenario.EndingID
	}{
		{
			name: "target is missing",
		},
		{
			name:       "both targets are set",
			nextNodeID: testfixture.MiddleNodeID,
			endingID:   testfixture.SafeEndingID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := getValidInput()
			in.Scenario.Nodes[0].Choices[0].NextNodeID = tt.nextNodeID
			in.Scenario.Nodes[0].Choices[0].EndingID = tt.endingID

			result, err := engine.Transition(in)

			require.ErrorIs(t, err, engine.ErrInvalidChoiceTarget)
			require.Empty(t, result)
		})
	}
}
