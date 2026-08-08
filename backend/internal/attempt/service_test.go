package attempt_test

import (
	"context"
	"errors"
	"testing"

	"github.com/marlendd/anti-scam-trainer/internal/attempt"
	"github.com/marlendd/anti-scam-trainer/internal/scenario"
	"github.com/marlendd/anti-scam-trainer/internal/testfixture"
	"github.com/stretchr/testify/require"
)

const (
	userID     = "user-1"
	scenarioID = scenario.ScenarioID("scenario-v1")
)

type attemptRepositoryStub struct {
	withinTransactionFn func(
		ctx context.Context,
		fn func(attempt.AttemptRepository) error,
	) error
	createFn func(
		ctx context.Context,
		userID string,
		scenarioID scenario.ScenarioID,
		startNodeID scenario.NodeID,
	) (attempt.Attempt, error)
	getActiveFn func(
		ctx context.Context,
		userID string,
		scenarioID scenario.ScenarioID,
	) (attempt.Attempt, error)
	abortFn func(
		ctx context.Context,
		attemptID attempt.AttemptID,
		userID string,
	) error
}

func (s *attemptRepositoryStub) WithinTransaction(
	ctx context.Context,
	fn func(attempt.AttemptRepository) error,
) error {
	if s.withinTransactionFn == nil {
		panic("unexpected WithinTransaction call")
	}

	return s.withinTransactionFn(ctx, fn)
}

func (s *attemptRepositoryStub) Create(
	ctx context.Context,
	userID string,
	scenarioID scenario.ScenarioID,
	startNodeID scenario.NodeID,
) (attempt.Attempt, error) {
	if s.createFn == nil {
		panic("unexpected Create call")
	}
	return s.createFn(ctx, userID, scenarioID, startNodeID)
}

func (s *attemptRepositoryStub) GetByID(
	context.Context,
	attempt.AttemptID,
	string,
) (attempt.Attempt, error) {
	panic("unexpected GetByID call")
}

func (s *attemptRepositoryStub) GetActive(
	ctx context.Context,
	userID string,
	scenarioID scenario.ScenarioID,
) (attempt.Attempt, error) {
	if s.getActiveFn == nil {
		panic("unexpected GetActive call")
	}

	return s.getActiveFn(ctx, userID, scenarioID)
}

func (s *attemptRepositoryStub) Abort(
	ctx context.Context,
	attemptID attempt.AttemptID,
	userID string,
) error {
	if s.abortFn == nil {
		panic("unexpected Abort call")
	}

	return s.abortFn(ctx, attemptID, userID)
}

type scenarioProviderStub struct {
	getActiveByIDFn func(
		ctx context.Context,
		scenarioID scenario.ScenarioID,
	) (scenario.Scenario, error)
	getByIDFn func(
		ctx context.Context,
		scenarioID scenario.ScenarioID,
	) (scenario.Scenario, error)
}

func (s *scenarioProviderStub) GetActiveByID(
	ctx context.Context,
	scenarioID scenario.ScenarioID,
) (scenario.Scenario, error) {
	return s.getActiveByIDFn(ctx, scenarioID)
}

func (s *scenarioProviderStub) GetByID(
	ctx context.Context,
	scenarioID scenario.ScenarioID,
) (scenario.Scenario, error) {
	if s.getByIDFn == nil {
		panic("unexpected GetByID call")
	}

	return s.getByIDFn(ctx, scenarioID)
}

func TestServiceStart(t *testing.T) {
	t.Parallel()

	t.Run("starts attempt from scenario start node", func(t *testing.T) {
		t.Parallel()

		currentScenario := testfixture.ValidScenario()
		expectedAttempt := attempt.Attempt{
			ID:         "attempt-1",
			UserID:     userID,
			ScenarioID: currentScenario.ID,
			Status:     attempt.StatusInProgress,
		}

		provider := &scenarioProviderStub{
			getActiveByIDFn: func(
				_ context.Context,
				scenarioID scenario.ScenarioID,
			) (scenario.Scenario, error) {
				require.Equal(t, currentScenario.ID, scenarioID)
				return currentScenario, nil
			},
		}

		repository := &attemptRepositoryStub{
			createFn: func(
				_ context.Context,
				gotUserID string,
				scenarioID scenario.ScenarioID,
				startNodeID scenario.NodeID,
			) (attempt.Attempt, error) {
				require.Equal(t, userID, gotUserID)
				require.Equal(t, currentScenario.ID, scenarioID)
				require.Equal(t, currentScenario.StartNodeID, startNodeID)
				return expectedAttempt, nil
			},
		}

		service := attempt.NewService(repository, provider)

		actualAttempt, err := service.Start(context.Background(), userID, currentScenario.ID)

		require.NoError(t, err)
		require.Equal(t, expectedAttempt, actualAttempt)
	})

	t.Run("does not create attempt when scenario provider fails", func(t *testing.T) {
		t.Parallel()

		providerErr := errors.New("scenario provider failed")
		provider := &scenarioProviderStub{
			getActiveByIDFn: func(
				context.Context,
				scenario.ScenarioID,
			) (scenario.Scenario, error) {
				return scenario.Scenario{}, providerErr
			},
		}

		repository := &attemptRepositoryStub{
			createFn: func(
				context.Context,
				string,
				scenario.ScenarioID,
				scenario.NodeID,
			) (attempt.Attempt, error) {
				t.Fatal("Create must not be called")
				return attempt.Attempt{}, nil
			},
		}

		service := attempt.NewService(repository, provider)

		_, err := service.Start(context.Background(), userID, "scenario-v1")

		require.ErrorIs(t, err, providerErr)
	})

	t.Run("does not create attempt for invalid scenario", func(t *testing.T) {
		t.Parallel()

		invalidScenario := testfixture.ValidScenario()
		invalidScenario.StartNodeID = ""

		provider := &scenarioProviderStub{
			getActiveByIDFn: func(
				context.Context,
				scenario.ScenarioID,
			) (scenario.Scenario, error) {
				return invalidScenario, nil
			},
		}

		repository := &attemptRepositoryStub{
			createFn: func(
				context.Context,
				string,
				scenario.ScenarioID,
				scenario.NodeID,
			) (attempt.Attempt, error) {
				t.Fatal("Create must not be called")
				return attempt.Attempt{}, nil
			},
		}

		service := attempt.NewService(repository, provider)

		_, err := service.Start(context.Background(), userID, invalidScenario.ID)

		require.ErrorIs(t, err, scenario.ErrEmptyStartNodeID)
	})

	t.Run("preserves active attempt conflict", func(t *testing.T) {
		t.Parallel()

		currentScenario := testfixture.ValidScenario()
		provider := &scenarioProviderStub{
			getActiveByIDFn: func(
				context.Context,
				scenario.ScenarioID,
			) (scenario.Scenario, error) {
				return currentScenario, nil
			},
		}

		repository := &attemptRepositoryStub{
			createFn: func(
				context.Context,
				string,
				scenario.ScenarioID,
				scenario.NodeID,
			) (attempt.Attempt, error) {
				return attempt.Attempt{}, attempt.ErrActiveAttemptExists
			},
		}

		service := attempt.NewService(repository, provider)

		_, err := service.Start(context.Background(), userID, currentScenario.ID)

		require.ErrorIs(t, err, attempt.ErrActiveAttemptExists)
	})
}

func TestService_Resume(t *testing.T) {
	t.Parallel()

	t.Run("active attempt found", func(t *testing.T) {
		t.Parallel()

		currentNodeID := scenario.NodeID("node-start")
		expectedAttempt := attempt.Attempt{
			ID:            "attempt-1",
			UserID:        "user-1",
			ScenarioID:    "scenario-v1",
			Status:        attempt.StatusInProgress,
			CurrentNodeID: &currentNodeID,
		}

		repository := &attemptRepositoryStub{
			getActiveFn: func(
				ctx context.Context,
				gotUserID string,
				gotScenarioID scenario.ScenarioID,
			) (attempt.Attempt, error) {
				require.Equal(t, userID, gotUserID)
				require.Equal(t, scenarioID, gotScenarioID)

				return expectedAttempt, nil
			},
		}

		provider := &scenarioProviderStub{}
		service := attempt.NewService(repository, provider)

		actualAttempt, err := service.Resume(
			context.Background(),
			userID,
			scenarioID,
		)

		require.NoError(t, err)
		require.Equal(t, expectedAttempt, actualAttempt)
	})

	t.Run("active attempt not found", func(t *testing.T) {
		t.Parallel()

		repository := &attemptRepositoryStub{
			getActiveFn: func(
				_ context.Context,
				_ string,
				_ scenario.ScenarioID,
			) (attempt.Attempt, error) {
				return attempt.Attempt{}, attempt.ErrActiveAttemptNotFound
			},
		}

		provider := &scenarioProviderStub{}

		service := attempt.NewService(repository, provider)

		actualAttempt, err := service.Resume(
			context.Background(),
			userID,
			scenarioID,
		)

		require.ErrorIs(t, err, attempt.ErrActiveAttemptNotFound)
		require.Equal(t, attempt.Attempt{}, actualAttempt)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		repositoryErr := errors.New("repository failed")
		repository := &attemptRepositoryStub{
			getActiveFn: func(
				_ context.Context,
				_ string,
				_ scenario.ScenarioID,
			) (attempt.Attempt, error) {
				return attempt.Attempt{}, repositoryErr
			},
		}

		provider := &scenarioProviderStub{}
		service := attempt.NewService(repository, provider)

		actualAttempt, err := service.Resume(
			context.Background(),
			userID,
			scenarioID,
		)

		require.ErrorIs(t, err, repositoryErr)
		require.Equal(t, attempt.Attempt{}, actualAttempt)
	})
}

func TestServiceRestart(t *testing.T) {
	t.Parallel()

	t.Run("aborts active attempt and creates a new one", func(t *testing.T) {
		t.Parallel()

		currentScenario := testfixture.ValidScenario()
		currentAttempt := attempt.Attempt{
			ID:         "attempt-old",
			UserID:     userID,
			ScenarioID: currentScenario.ID,
			Status:     attempt.StatusInProgress,
		}
		expectedAttempt := attempt.Attempt{
			ID:         "attempt-new",
			UserID:     userID,
			ScenarioID: currentScenario.ID,
			Status:     attempt.StatusInProgress,
		}

		provider := scenarioStub(currentScenario)
		txRepository := &attemptRepositoryStub{
			getActiveFn: func(
				_ context.Context,
				gotUserID string,
				gotScenarioID scenario.ScenarioID,
			) (attempt.Attempt, error) {
				require.Equal(t, userID, gotUserID)
				require.Equal(t, currentScenario.ID, gotScenarioID)
				return currentAttempt, nil
			},
			abortFn: func(
				_ context.Context,
				gotAttemptID attempt.AttemptID,
				gotUserID string,
			) error {
				require.Equal(t, currentAttempt.ID, gotAttemptID)
				require.Equal(t, userID, gotUserID)
				return nil
			},
			createFn: func(
				_ context.Context,
				gotUserID string,
				gotScenarioID scenario.ScenarioID,
				gotStartNodeID scenario.NodeID,
			) (attempt.Attempt, error) {
				require.Equal(t, userID, gotUserID)
				require.Equal(t, currentScenario.ID, gotScenarioID)
				require.Equal(t, currentScenario.StartNodeID, gotStartNodeID)
				return expectedAttempt, nil
			},
		}
		repository := transactionStub(txRepository)
		service := attempt.NewService(repository, provider)

		actualAttempt, err := service.Restart(
			context.Background(),
			userID,
			currentScenario.ID,
		)

		require.NoError(t, err)
		require.Equal(t, expectedAttempt, actualAttempt)
	})

	t.Run("does not start transaction when scenario provider fails", func(t *testing.T) {
		t.Parallel()

		providerErr := errors.New("scenario provider failed")
		provider := &scenarioProviderStub{
			getActiveByIDFn: func(
				context.Context,
				scenario.ScenarioID,
			) (scenario.Scenario, error) {
				return scenario.Scenario{}, providerErr
			},
		}
		repository := &attemptRepositoryStub{}
		service := attempt.NewService(repository, provider)

		actualAttempt, err := service.Restart(
			context.Background(),
			userID,
			scenarioID,
		)

		require.ErrorIs(t, err, providerErr)
		require.Equal(t, attempt.Attempt{}, actualAttempt)
	})

	t.Run("preserves get active attempt error", func(t *testing.T) {
		t.Parallel()

		currentScenario := testfixture.ValidScenario()
		txRepository := &attemptRepositoryStub{
			getActiveFn: func(
				context.Context,
				string,
				scenario.ScenarioID,
			) (attempt.Attempt, error) {
				return attempt.Attempt{}, attempt.ErrActiveAttemptNotFound
			},
		}
		service := attempt.NewService(
			transactionStub(txRepository),
			scenarioStub(currentScenario),
		)

		actualAttempt, err := service.Restart(
			context.Background(),
			userID,
			currentScenario.ID,
		)

		require.ErrorIs(t, err, attempt.ErrActiveAttemptNotFound)
		require.Equal(t, attempt.Attempt{}, actualAttempt)
	})

	t.Run("preserves abort error and does not create attempt", func(t *testing.T) {
		t.Parallel()

		currentScenario := testfixture.ValidScenario()
		abortErr := errors.New("abort failed")
		txRepository := &attemptRepositoryStub{
			getActiveFn: func(
				context.Context,
				string,
				scenario.ScenarioID,
			) (attempt.Attempt, error) {
				return attempt.Attempt{ID: "attempt-old"}, nil
			},
			abortFn: func(
				context.Context,
				attempt.AttemptID,
				string,
			) error {
				return abortErr
			},
		}
		service := attempt.NewService(
			transactionStub(txRepository),
			scenarioStub(currentScenario),
		)

		actualAttempt, err := service.Restart(
			context.Background(),
			userID,
			currentScenario.ID,
		)

		require.ErrorIs(t, err, abortErr)
		require.Equal(t, attempt.Attempt{}, actualAttempt)
	})

	t.Run("preserves create error", func(t *testing.T) {
		t.Parallel()

		currentScenario := testfixture.ValidScenario()
		createErr := errors.New("create failed")
		txRepository := &attemptRepositoryStub{
			getActiveFn: func(
				context.Context,
				string,
				scenario.ScenarioID,
			) (attempt.Attempt, error) {
				return attempt.Attempt{ID: "attempt-old"}, nil
			},
			abortFn: func(
				context.Context,
				attempt.AttemptID,
				string,
			) error {
				return nil
			},
			createFn: func(
				context.Context,
				string,
				scenario.ScenarioID,
				scenario.NodeID,
			) (attempt.Attempt, error) {
				return attempt.Attempt{}, createErr
			},
		}
		service := attempt.NewService(
			transactionStub(txRepository),
			scenarioStub(currentScenario),
		)

		actualAttempt, err := service.Restart(
			context.Background(),
			userID,
			currentScenario.ID,
		)

		require.ErrorIs(t, err, createErr)
		require.Equal(t, attempt.Attempt{}, actualAttempt)
	})
}

func transactionStub(txRepository attempt.AttemptRepository) *attemptRepositoryStub {
	return &attemptRepositoryStub{
		withinTransactionFn: func(
			_ context.Context,
			fn func(attempt.AttemptRepository) error,
		) error {
			return fn(txRepository)
		},
	}
}

func scenarioStub(currentScenario scenario.Scenario) *scenarioProviderStub {
	return &scenarioProviderStub{
		getActiveByIDFn: func(
			context.Context,
			scenario.ScenarioID,
		) (scenario.Scenario, error) {
			return currentScenario, nil
		},
	}
}
