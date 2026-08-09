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

func (pg *PgRepository) ListActiveByRole(
	ctx context.Context,
	userID string,
	role Role,
) ([]CatalogItem, error) {
	const query = `
		SELECT active.id,
		       active.logical_id,
		       active.version,
		       active.role,
		       active.title,
		       active.description,
		       CASE
		           WHEN bool_or(a.status = 'in_progress') THEN 'in_progress'
		           WHEN bool_or(a.status = 'completed') THEN 'completed'
		           ELSE 'not_started'
		       END AS progress_status
		FROM scenario_versions AS active
		LEFT JOIN scenario_versions AS history
		       ON history.logical_id = active.logical_id
		LEFT JOIN attempts AS a
		       ON a.scenario_id = history.id
		      AND a.user_id = $1
		WHERE active.is_active = TRUE
		  AND active.role = $2
		GROUP BY active.id
		ORDER BY active.created_at, active.id
	`

	rows, err := pg.db.QueryContext(ctx, query, userID, role)
	if err != nil {
		return nil, fmt.Errorf("query active scenario catalog: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	items := make([]CatalogItem, 0)
	for rows.Next() {
		var item CatalogItem
		if err := rows.Scan(
			&item.ID,
			&item.LogicalID,
			&item.Version,
			&item.Role,
			&item.Title,
			&item.Description,
			&item.Status,
		); err != nil {
			return nil, fmt.Errorf("scan scenario catalog item: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate scenario catalog: %w", err)
	}

	return items, nil
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
						  COALESCE(reward_fragment_id, ''),
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
		&result.RewardFragmentID,
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
		return Content{}, fmt.Errorf("%w: %w", ErrInvalidScenarioContent, err)
	}

	return content, nil
}

func (pg *PgRepository) getScenarioFromContent(scenario Scenario, content Content) (Scenario, error) {
	scenario.StartNodeID = content.StartNodeID
	scenario.SuccessfulEndingIDs = content.SuccessfulEndingIDs
	scenario.Nodes = content.Nodes
	scenario.Endings = content.Endings

	if err := Validate(scenario); err != nil {
		return Scenario{}, fmt.Errorf("%w: %w", ErrInvalidScenarioContent, err)
	}

	return scenario, nil
}
