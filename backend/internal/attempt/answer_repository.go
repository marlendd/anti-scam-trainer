package attempt

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/marlendd/anti-scam-trainer/internal/scenario"
)

func (pg *PgRepository) WithinAnswerTransaction(
	ctx context.Context,
	fn func(AnswerRepository) error,
) error {
	return pg.withinTransaction(ctx, func(txRepository *PgRepository) error {
		return fn(txRepository)
	})
}

func (pg *PgRepository) GetAttemptForUpdate(
	ctx context.Context,
	attemptID AttemptID,
	userID string,
) (Attempt, error) {
	const query = `SELECT
						id,
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
					FOR UPDATE
	`

	currentAttempt, err := scanAttempt(pg.executor.QueryRowContext(
		ctx,
		query,
		attemptID,
		userID,
	))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Attempt{}, ErrAttemptNotFound
		}

		return Attempt{}, fmt.Errorf("get attempt for update: %w", err)
	}

	return currentAttempt, nil
}

func (pg *PgRepository) GetAnswerByIdempotencyKey(
	ctx context.Context,
	userID string,
	idempotencyKey IdempotencyKey,
) (Answer, error) {
	const query = `SELECT
						a.id,
						a.attempt_id,
						a.node_id,
						a.choice_id,
						a.idempotency_key,
						a.weight,
						a.choice_score,
						a.risk_categories,
						a.consequence,
						a.explanation,
						a.response,
						a.created_at
					FROM answers AS a
					JOIN attempts AS att ON att.id = a.attempt_id
					WHERE a.idempotency_key = $1 AND att.user_id = $2
	`

	answer, err := scanAnswer(pg.executor.QueryRowContext(
		ctx,
		query,
		idempotencyKey,
		userID,
	))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Answer{}, ErrAnswerNotFound
		}

		return Answer{}, fmt.Errorf("get answer by idempotency key: %w", err)
	}

	return answer, nil
}

func (pg *PgRepository) GetAnswerByAttemptNode(
	ctx context.Context,
	attemptID AttemptID,
	userID string,
	nodeID scenario.NodeID,
) (Answer, error) {
	const query = `SELECT
						a.id,
						a.attempt_id,
						a.node_id,
						a.choice_id,
						a.idempotency_key,
						a.weight,
						a.choice_score,
						a.risk_categories,
						a.consequence,
						a.explanation,
						a.response,
						a.created_at
					FROM answers AS a
					JOIN attempts AS att ON att.id = a.attempt_id
					WHERE a.attempt_id = $1
						AND att.user_id = $2
						AND a.node_id = $3
	`

	answer, err := scanAnswer(pg.executor.QueryRowContext(
		ctx,
		query,
		attemptID,
		userID,
		nodeID,
	))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Answer{}, ErrAnswerNotFound
		}

		return Answer{}, fmt.Errorf("get answer by attempt node: %w", err)
	}

	return answer, nil
}

func (pg *PgRepository) GetScoreTotals(
	ctx context.Context,
	attemptID AttemptID,
) (ScoreTotals, error) {
	const query = `SELECT
						COALESCE(SUM(weight * choice_score), 0),
						COALESCE(SUM(weight), 0)
					FROM answers
					WHERE attempt_id = $1
	`

	var (
		weightedScoreSum int64
		weightSum        int64
	)

	if err := pg.executor.QueryRowContext(ctx, query, attemptID).Scan(
		&weightedScoreSum,
		&weightSum,
	); err != nil {
		return ScoreTotals{}, fmt.Errorf("get answer score totals: %w", err)
	}

	return ScoreTotals{
		WeightedScoreSum: int(weightedScoreSum),
		WeightSum:        int(weightSum),
	}, nil
}

func (pg *PgRepository) CreateAnswer(
	ctx context.Context,
	answer Answer,
) (Answer, error) {
	riskCategories := answer.RiskCategories
	if riskCategories == nil {
		riskCategories = []scenario.RiskCategory{}
	}

	riskCategoriesJSON, err := json.Marshal(riskCategories)
	if err != nil {
		return Answer{}, fmt.Errorf("encode answer risk categories: %w", err)
	}

	responseJSON, err := json.Marshal(answer.Response)
	if err != nil {
		return Answer{}, fmt.Errorf("encode answer response: %w", err)
	}

	const query = `INSERT INTO answers (
						attempt_id,
						node_id,
						choice_id,
						idempotency_key,
						weight,
						choice_score,
						risk_categories,
						consequence,
						explanation,
						response
					)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
					RETURNING
						id,
						attempt_id,
						node_id,
						choice_id,
						idempotency_key,
						weight,
						choice_score,
						risk_categories,
						consequence,
						explanation,
						response,
						created_at
	`

	createdAnswer, err := scanAnswer(pg.executor.QueryRowContext(
		ctx,
		query,
		answer.AttemptID,
		answer.NodeID,
		answer.ChoiceID,
		answer.IdempotencyKey,
		answer.Weight,
		answer.ChoiceScore,
		string(riskCategoriesJSON),
		answer.Consequence,
		answer.Explanation,
		string(responseJSON),
	))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "answers_idempotency_key_key":
				return Answer{}, ErrIdempotencyConflict
			case "answers_attempt_id_node_id_key":
				return Answer{}, ErrNodeAlreadyAnswered
			}
		}

		return Answer{}, fmt.Errorf("create answer: %w", err)
	}

	return createdAnswer, nil
}

func (pg *PgRepository) AdvanceAttempt(
	ctx context.Context,
	attemptID AttemptID,
	userID string,
	currentNodeID scenario.NodeID,
	nextNodeID scenario.NodeID,
) error {
	const query = `UPDATE attempts
					SET current_node_id = $1,
						updated_at = now()
					WHERE id = $2
						AND user_id = $3
						AND status = $4
						AND current_node_id = $5
	`

	result, err := pg.executor.ExecContext(
		ctx,
		query,
		nextNodeID,
		attemptID,
		userID,
		StatusInProgress,
		currentNodeID,
	)
	if err != nil {
		return fmt.Errorf("advance attempt: %w", err)
	}

	return requireOneAffectedRow(result, "advance attempt")
}

func (pg *PgRepository) CompleteAttempt(
	ctx context.Context,
	attemptID AttemptID,
	userID string,
	currentNodeID scenario.NodeID,
	endingID scenario.EndingID,
	score int,
) error {
	const query = `UPDATE attempts
					SET status = $1,
						current_node_id = NULL,
						ending_id = $2,
						score = $3,
						updated_at = now(),
						completed_at = now()
					WHERE id = $4
						AND user_id = $5
						AND status = $6
						AND current_node_id = $7
	`

	result, err := pg.executor.ExecContext(
		ctx,
		query,
		StatusCompleted,
		endingID,
		score,
		attemptID,
		userID,
		StatusInProgress,
		currentNodeID,
	)
	if err != nil {
		return fmt.Errorf("complete attempt: %w", err)
	}

	return requireOneAffectedRow(result, "complete attempt")
}

func (pg *PgRepository) GrantFragment(
	ctx context.Context,
	userID string,
	scenarioID scenario.ScenarioID,
	fragmentID scenario.FragmentID,
) (bool, error) {
	const query = `INSERT INTO user_inventory (
						user_id,
						scenario_id,
						fragment_id
					)
					VALUES ($1, $2, $3)
					ON CONFLICT DO NOTHING
	`

	result, err := pg.executor.ExecContext(
		ctx,
		query,
		userID,
		scenarioID,
		fragmentID,
	)
	if err != nil {
		return false, fmt.Errorf("grant reward fragment: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("get granted fragment rows affected: %w", err)
	}

	return affected == 1, nil
}

func scanAnswer(row rowScanner) (Answer, error) {
	var (
		answer             Answer
		riskCategoriesJSON json.RawMessage
		responseJSON       json.RawMessage
	)

	if err := row.Scan(
		&answer.ID,
		&answer.AttemptID,
		&answer.NodeID,
		&answer.ChoiceID,
		&answer.IdempotencyKey,
		&answer.Weight,
		&answer.ChoiceScore,
		&riskCategoriesJSON,
		&answer.Consequence,
		&answer.Explanation,
		&responseJSON,
		&answer.CreatedAt,
	); err != nil {
		return Answer{}, err
	}

	if err := json.Unmarshal(riskCategoriesJSON, &answer.RiskCategories); err != nil {
		return Answer{}, fmt.Errorf("decode answer risk categories: %w", err)
	}

	if err := json.Unmarshal(responseJSON, &answer.Response); err != nil {
		return Answer{}, fmt.Errorf("decode answer response: %w", err)
	}

	return answer, nil
}

func requireOneAffectedRow(result sql.Result, operation string) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get %s rows affected: %w", operation, err)
	}

	if affected != 1 {
		return ErrInvalidAttemptState
	}

	return nil
}
