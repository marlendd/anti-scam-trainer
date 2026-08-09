package scenario

import (
	"context"
	"fmt"
)

type ProgressStatus string

const (
	ProgressNotStarted ProgressStatus = "not_started"
	ProgressInProgress ProgressStatus = "in_progress"
	ProgressCompleted  ProgressStatus = "completed"
)

type CatalogItem struct {
	ID          ScenarioID
	LogicalID   LogicalScenarioID
	Version     int
	Role        Role
	Title       string
	Description string
	Product     Product
	Status      ProgressStatus
}

type CatalogRepository interface {
	ListActiveByRole(
		ctx context.Context,
		userID string,
		role Role,
	) ([]CatalogItem, error)
}

type CatalogService struct {
	repository CatalogRepository
}

func NewCatalogService(repository CatalogRepository) *CatalogService {
	return &CatalogService{repository: repository}
}

func (service *CatalogService) ListActiveByRole(
	ctx context.Context,
	userID string,
	role Role,
) ([]CatalogItem, error) {
	if role != RoleBuyer && role != RoleSeller {
		return nil, ErrInvalidRole
	}

	items, err := service.repository.ListActiveByRole(ctx, userID, role)
	if err != nil {
		return nil, fmt.Errorf("list active scenarios by role: %w", err)
	}

	return items, nil
}
