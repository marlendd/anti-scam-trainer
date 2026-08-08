package scenario_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/marlendd/anti-scam-trainer/internal/scenario"
)

type catalogRepositoryStub struct {
	items     []scenario.CatalogItem
	err       error
	gotUserID string
	gotRole   scenario.Role
	callCount int
}

func (stub *catalogRepositoryStub) ListActiveByRole(
	_ context.Context,
	userID string,
	role scenario.Role,
) ([]scenario.CatalogItem, error) {
	stub.callCount++
	stub.gotUserID = userID
	stub.gotRole = role
	return stub.items, stub.err
}

func TestCatalogService_ListActiveByRole(t *testing.T) {
	t.Parallel()

	expected := []scenario.CatalogItem{{ID: "scenario-1"}}
	repository := &catalogRepositoryStub{items: expected}
	service := scenario.NewCatalogService(repository)

	items, err := service.ListActiveByRole(
		context.Background(),
		"user-1",
		scenario.RoleBuyer,
	)

	require.NoError(t, err)
	require.Equal(t, expected, items)
	require.Equal(t, 1, repository.callCount)
	require.Equal(t, "user-1", repository.gotUserID)
	require.Equal(t, scenario.RoleBuyer, repository.gotRole)
}

func TestCatalogService_ListActiveByRoleRejectsInvalidRole(t *testing.T) {
	t.Parallel()

	repository := &catalogRepositoryStub{}
	service := scenario.NewCatalogService(repository)

	_, err := service.ListActiveByRole(context.Background(), "user-1", "admin")

	require.ErrorIs(t, err, scenario.ErrInvalidRole)
	require.Zero(t, repository.callCount)
}

func TestCatalogService_ListActiveByRoleWrapsRepositoryError(t *testing.T) {
	t.Parallel()

	repositoryErr := errors.New("database unavailable")
	service := scenario.NewCatalogService(&catalogRepositoryStub{err: repositoryErr})

	_, err := service.ListActiveByRole(
		context.Background(),
		"user-1",
		scenario.RoleSeller,
	)

	require.ErrorIs(t, err, repositoryErr)
}
