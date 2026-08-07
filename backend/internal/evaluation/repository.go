package evaluation

import (
	"context"
	"database/sql"
)

type Repository interface {
	GetAnswersByAttempt(ctx context.Context, attemptID string) ([]AnswerData, error)
	GetStatsByRole(ctx context.Context) ([]RoleStats, error)
	SaveReward(ctx context.Context, userID, fragmentID string) error
	GetUserFragments(ctx context.Context, userID string) ([]PuzzleFragment, error)
	GetTotalAvailableFragments(ctx context.Context) (int, error)
	GetStatsByCategory(ctx context.Context) ([]CategoryStat, error)
	GetTotalCompletedCount(ctx context.Context) (int, error)
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

func (r *PgRepository) SaveReward(ctx context.Context, userID, fragmentID string) error {
	const q = `
		INSERT INTO user_inventory (user_id, fragment_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, fragment_id) DO NOTHING`

	_, err := r.db.ExecContext(ctx, q, userID, fragmentID)
	return err
}

func (r *PgRepository) GetUserFragments(ctx context.Context, userID string) ([]PuzzleFragment, error) {
	const q = `SELECT fragment_id, earned_at FROM user_inventory WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fragments []PuzzleFragment
	for rows.Next() {
		var f PuzzleFragment
		if err := rows.Scan(&f.FragmentID, &f.EarnedAt); err != nil {
			return nil, err
		}
		fragments = append(fragments, f)
	}
	return fragments, rows.Err()
}

func (r *PgRepository) GetTotalAvailableFragments(ctx context.Context) (int, error) {
	const q = `SELECT COUNT(DISTINCT reward_fragment_id) FROM scenario_versions WHERE reward_fragment_id IS NOT NULL`
	var count int
	err := r.db.QueryRowContext(ctx, q).Scan(&count)
	return count, err
}
func (r *PgRepository) GetStatsByCategory(ctx context.Context) ([]CategoryStat, error) {
	const q = `
		SELECT 
			category, 
			COUNT(DISTINCT scenario_id) as count
		FROM (
			SELECT 
				jsonb_array_elements_text(sv.content->'risk_categories') as category,
				a.scenario_id
			FROM attempts a
			JOIN scenario_versions sv ON a.scenario_id = sv.id
			WHERE a.status = 'completed'
		) as expanded_categories
		GROUP BY category
		ORDER BY count DESC`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []CategoryStat
	for rows.Next() {
		var s CategoryStat
		if err := rows.Scan(&s.Category, &s.Count); err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, rows.Err()
}

func (r *PgRepository) GetTotalCompletedCount(ctx context.Context) (int, error) {
	const q = `SELECT COUNT(DISTINCT scenario_id) FROM attempts WHERE status = 'completed'`
	var total int
	err := r.db.QueryRowContext(ctx, q).Scan(&total)
	return total, err
}
