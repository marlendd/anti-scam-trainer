package attempt

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/marlendd/anti-scam-trainer/internal/scenario"
)

type dbExecutor interface {
	ExecContext(
		ctx context.Context,
		query string,
		args ...any,
	) (sql.Result, error)

	QueryRowContext(
		ctx context.Context,
		query string,
		args ...any,
	) *sql.Row
}

type PgRepository struct {
	pool     *sql.DB
	executor dbExecutor
}

func NewPgRepository(db *sql.DB) *PgRepository {
	return &PgRepository{
		pool:     db,
		executor: db,
	}
}

func (pg *PgRepository) Create(
	ctx context.Context,
	userID string,
	scenarioID scenario.ScenarioID,
	startNodeID scenario.NodeID,
) (Attempt, error) {
	const query = `INSERT INTO attempts (
			  user_id,
			  scenario_id,
			  status,
			  current_node_id
			)
			  VALUES ($1, $2, $3, $4)
			  RETURNING
			  id,
			  user_id,
			  scenario_id,
			  status,
			  current_node_id,
			  started_at,
			  updated_at;
			`
	var attempt Attempt
	attempt.CurrentNodeID = new(scenario.NodeID)

	if err := pg.executor.QueryRowContext(ctx,
		query,
		userID,
		scenarioID,
		StatusInProgress,
		startNodeID,
	).Scan(
		&attempt.ID,
		&attempt.UserID,
		&attempt.ScenarioID,
		&attempt.Status,
		attempt.CurrentNodeID,
		&attempt.StartedAt,
		&attempt.UpdatedAt,
	); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" &&
			pgErr.ConstraintName == "attempts_one_in_progress_idx" {
			return Attempt{}, ErrActiveAttemptExists
		}

		return Attempt{}, fmt.Errorf("create attempt: %w", err)
	}

	return attempt, nil
}

func (pg *PgRepository) GetByID(
	ctx context.Context,
	attemptID AttemptID,
	userID string,
) (Attempt, error) {
	const query = `SELECT id,
						  user_id,
						  scenario_id,
						  status,
						  current_node_id,
						  ending_id,
						  score,
						  started_at,
						  updated_at,
						  completed_at
					FROM attempts
					WHERE id = $1 AND user_id = $2
	`
	var attempt Attempt

	var (
		currentNodeID sql.NullString
		endingID      sql.NullString
		score         sql.NullInt64
		completedAt   sql.NullTime
	)

	if err := pg.executor.QueryRowContext(
		ctx,
		query,
		attemptID,
		userID,
	).Scan(
		&attempt.ID,
		&attempt.UserID,
		&attempt.ScenarioID,
		&attempt.Status,
		&currentNodeID,
		&endingID,
		&score,
		&attempt.StartedAt,
		&attempt.UpdatedAt,
		&completedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Attempt{}, ErrAttemptNotFound
		}
		return Attempt{}, fmt.Errorf("get attempt by id: %w", err)
	}

	if currentNodeID.Valid {
		value := scenario.NodeID(currentNodeID.String)
		attempt.CurrentNodeID = &value
	}

	if endingID.Valid {
		value := scenario.EndingID(endingID.String)
		attempt.EndingID = &value
	}

	if score.Valid {
		value := int(score.Int64)
		attempt.Score = &value
	}

	if completedAt.Valid {
		value := completedAt.Time
		attempt.CompletedAt = &value
	}

	return attempt, nil
}

func (pg *PgRepository) GetActive(
	ctx context.Context,
	userID string,
	scenarioID scenario.ScenarioID,
) (Attempt, error) {
	const query = `SELECT 
				id,
				user_id,
				scenario_id,
				status,
				current_node_id,
				started_at,
				updated_at
				FROM attempts
				WHERE user_id = $1 AND scenario_id = $2 AND status = $3
	`

	var attempt Attempt
	attempt.CurrentNodeID = new(scenario.NodeID)

	if err := pg.executor.QueryRowContext(
		ctx,
		query,
		userID,
		scenarioID,
		StatusInProgress,
	).Scan(
		&attempt.ID,
		&attempt.UserID,
		&attempt.ScenarioID,
		&attempt.Status,
		attempt.CurrentNodeID,
		&attempt.StartedAt,
		&attempt.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Attempt{}, ErrActiveAttemptNotFound
		}
		return Attempt{}, fmt.Errorf("get active attempt: %w", err)
	}

	return attempt, nil
}

func (pg *PgRepository) Abort(
	ctx context.Context,
	attemptID AttemptID,
	userID string,
) error {
	const query = `UPDATE attempts
				   SET 
				   		status = $1,
				    	updated_at = now()
				   WHERE 
				   		id = $2 
				   		AND user_id = $3
						AND status = $4
	`

	result, err := pg.executor.ExecContext(
		ctx,
		query,
		StatusAborted,
		attemptID,
		userID,
		StatusInProgress,
	)

	if err != nil {
		return fmt.Errorf("abort attempt: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get abort attempt rows affected: %w", err)
	}

	if affected == 1 {
		return nil
	}

	const statusQuery = `SELECT status 
						 FROM attempts
						 WHERE id = $1 AND user_id = $2
	`

	var status Status

	if err := pg.executor.QueryRowContext(
		ctx,
		statusQuery,
		attemptID,
		userID,
	).Scan(&status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAttemptNotFound
		}
		return fmt.Errorf("get attempts status: %w", err)
	}

	return ErrAttemptNotInProgress
}

func (pg *PgRepository) WithinTransaction(
	ctx context.Context,
	fn func(AttemptRepository) error,
) error {
	tx, err := pg.pool.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin attempt transaction: %w", err)
	}

	txRepository := &PgRepository{
		pool:     pg.pool,
		executor: tx,
	}

	if err := fn(txRepository); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil &&
			!errors.Is(rollbackErr, sql.ErrTxDone) {
			return errors.Join(
				err,
				fmt.Errorf("rollback attempt transaction: %w", rollbackErr),
			)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit attempt transaction: %w", err)
	}

	return nil
}
