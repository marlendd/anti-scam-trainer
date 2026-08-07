package evaluation

import (
	"context"
	"database/sql"
)

type Repository interface {
	GetAnswersByAttempt(ctx context.Context, attemptID string) ([]AnswerData, error)
	GetStatsByRole(ctx context.Context) ([]RoleStats, error)
}

type PgRepository struct {
	db *sql.DB
}

func NewPgRepository(db *sql.DB) *PgRepository {
	return &PgRepository{db: db}
}

func (r *PgRepository) GetAnswersByAttempt(ctx context.Context, attemptID string) ([]AnswerData, error) {
	const q = `SELECT weight, choice_score FROM answers WHERE attempt_id = $1`

	rows, err := r.db.QueryContext(ctx, q, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []AnswerData
	for rows.Next() {
		var a AnswerData
		if err := rows.Scan(&a.Weight, &a.ChoiceScore); err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	return res, rows.Err()
}

func (r *PgRepository) GetStatsByRole(ctx context.Context) ([]RoleStats, error) {
	const q = `
		SELECT 
			sv.role,
			COUNT(*) FILTER (WHERE a.status = 'completed'),
			COUNT(*) FILTER (WHERE a.status = 'in_progress'),
			COUNT(*)
		FROM attempts a
		JOIN scenario_versions sv ON a.scenario_id = sv.id
		GROUP BY sv.role`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []RoleStats
	for rows.Next() {
		var s RoleStats
		if err := rows.Scan(&s.Role, &s.CompletedCount, &s.InProgressCount, &s.TotalStarted); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}
