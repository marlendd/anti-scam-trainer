package scenario

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

type scenarioRow struct {
	Scenario
	IsActive bool
	Content  json.RawMessage
}

type PgRepository struct {
	db *sql.DB
}

func NewPgRepository(db *sql.DB) PgRepository {
	return PgRepository{
		db: db,
	}
}

func (pg *PgRepository) GetByID(ctx context.Context, scenarioID ScenarioID) (Scenario, error) {
	scenarioRow, err := pg.getScenarioData(ctx, scenarioID)
	if err != nil {
		return Scenario{}, err
	}

	content, err := pg.decodeContent(scenarioRow.Content)
	if err != nil {
		return Scenario{}, err
	}

	scenario, err := pg.getScenarioFromContent(scenarioRow.Scenario, content)
	if err != nil {
		return Scenario{}, err
	}

	return scenario, nil
}

func (pg *PgRepository) GetActiveByID(ctx context.Context, scenarioID ScenarioID) (Scenario, error) {
	scenarioRow, err := pg.getScenarioData(ctx, scenarioID)
	if err != nil {
		return Scenario{}, err
	}
	if !scenarioRow.IsActive {
		return Scenario{}, ErrScenarioInactive
	}

	content, err := pg.decodeContent(scenarioRow.Content)
	if err != nil {
		return Scenario{}, err
	}

	scenario, err := pg.getScenarioFromContent(scenarioRow.Scenario, content)
	if err != nil {
		return Scenario{}, err
	}

	return scenario, nil
}

func (pg *PgRepository) getScenarioData(ctx context.Context, scenarioID ScenarioID) (scenarioRow, error) {
	const query = `SELECT id,
						  logical_id,
						  version,
						  role,
						  title,
						  description,
						  is_active,
						  content
					FROM scenario_versions
					WHERE id = $1
			    `
	var result scenarioRow
	if err := pg.db.QueryRowContext(
		ctx,
		query,
		scenarioID,
	).Scan(
		&result.ID,
		&result.LogicalID,
		&result.Version,
		&result.Role,
		&result.Title,
		&result.Description,
		&result.IsActive,
		&result.Content,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return scenarioRow{}, ErrScenarioNotFound
		}
		return scenarioRow{}, fmt.Errorf("get scenario row: %w", err)
	}

	return result, nil
}

func (pg *PgRepository) decodeContent(rawContent json.RawMessage) (Content, error) {
	var content Content
	if err := json.Unmarshal(rawContent, &content); err != nil {
		return Content{}, fmt.Errorf("%w: %v", ErrInvalidScenarioContent, err)
	}

	return content, nil
}

func (pg *PgRepository) getScenarioFromContent(scenario Scenario, content Content) (Scenario, error) {
	scenario.StartNodeID = content.StartNodeID
	scenario.Nodes = content.Nodes
	scenario.Endings = content.Endings

	if err := Validate(scenario); err != nil {
		return Scenario{}, err
	}

	return scenario, nil
}
