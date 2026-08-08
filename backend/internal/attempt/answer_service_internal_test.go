package attempt

import (
	"context"
	"errors"
	"testing"

	"github.com/marlendd/anti-scam-trainer/internal/scenario"
	"github.com/stretchr/testify/require"
)

type idempotencyAnswerGetterStub struct {
	answer Answer
	err    error
}

func (s idempotencyAnswerGetterStub) GetAnswerByIdempotencyKey(
	context.Context,
	string,
	IdempotencyKey,
) (Answer, error) {
	return s.answer, s.err
}

func TestFindIdempotentResult(t *testing.T) {
	t.Parallel()

	input := SubmitAnswerInput{
		UserID:         "user-1",
		AttemptID:      "attempt-1",
		NodeID:         "node-1",
		ChoiceID:       "choice-1",
		IdempotencyKey: "11111111-1111-1111-1111-111111111111",
	}
	expectedResult := SubmitAnswerResult{
		AttemptID:   input.AttemptID,
		NodeID:      input.NodeID,
		ChoiceID:    input.ChoiceID,
		Consequence: "Saved consequence",
	}

	t.Run("returns saved result for matching payload", func(t *testing.T) {
		t.Parallel()

		repository := idempotencyAnswerGetterStub{answer: Answer{
			AttemptID: input.AttemptID,
			NodeID:    input.NodeID,
			ChoiceID:  input.ChoiceID,
			Response:  expectedResult,
		}}

		result, found, err := findIdempotentResult(context.Background(), repository, input)

		require.NoError(t, err)
		require.True(t, found)
		require.Equal(t, expectedResult, result)
	})

	t.Run("reports missing key as a cache miss", func(t *testing.T) {
		t.Parallel()

		repository := idempotencyAnswerGetterStub{err: ErrAnswerNotFound}

		result, found, err := findIdempotentResult(context.Background(), repository, input)

		require.NoError(t, err)
		require.False(t, found)
		require.Empty(t, result)
	})

	t.Run("returns repository failure", func(t *testing.T) {
		t.Parallel()

		repositoryErr := errors.New("repository failed")
		repository := idempotencyAnswerGetterStub{err: repositoryErr}

		result, found, err := findIdempotentResult(context.Background(), repository, input)

		require.ErrorIs(t, err, repositoryErr)
		require.False(t, found)
		require.Empty(t, result)
	})

	cases := []struct {
		name   string
		answer Answer
	}{
		{
			name: "rejects a different attempt",
			answer: Answer{
				AttemptID: "attempt-2",
				NodeID:    input.NodeID,
				ChoiceID:  input.ChoiceID,
			},
		},
		{
			name: "rejects a different node",
			answer: Answer{
				AttemptID: input.AttemptID,
				NodeID:    scenario.NodeID("node-2"),
				ChoiceID:  input.ChoiceID,
			},
		},
		{
			name: "rejects a different choice",
			answer: Answer{
				AttemptID: input.AttemptID,
				NodeID:    input.NodeID,
				ChoiceID:  "choice-2",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, found, err := findIdempotentResult(
				context.Background(),
				idempotencyAnswerGetterStub{answer: tc.answer},
				input,
			)

			require.ErrorIs(t, err, ErrIdempotencyConflict)
			require.False(t, found)
			require.Empty(t, result)
		})
	}
}
