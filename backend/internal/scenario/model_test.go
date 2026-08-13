package scenario_test

import (
	"testing"

	"github.com/marlendd/anti-scam-trainer/internal/testfixture"
	"github.com/stretchr/testify/require"
)

func TestScenarioIsSuccessfulEnding(t *testing.T) {
	t.Parallel()

	currentScenario := testfixture.ValidScenario()

	require.True(t, currentScenario.IsSuccessfulEnding(testfixture.SafeEndingID))
	require.False(t, currentScenario.IsSuccessfulEnding(testfixture.RiskyEndingID))
}
