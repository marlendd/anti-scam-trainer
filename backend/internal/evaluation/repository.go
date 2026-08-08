package evaluation

import (
	"context"
	"database/sql"
)

type Repository interface {
	GetAnswersByAttempt(ctx context.Context, userID, attemptID string) ([]AnswerData, error)
	GetUserStatsByRole(ctx context.Context, userID string) ([]RoleStats, error)
	SaveReward(ctx context.Context, userID, fragmentID string) error
	GetUserFragments(ctx context.Context, userID string) ([]PuzzleFragment, error)
	GetTotalAvailableFragments(ctx context.Context) (int, error)
	GetUserStatsByCategory(ctx context.Context, userID string) ([]CategoryStat, error)
	GetUserTotalCompletedCount(ctx context.Context, userID string) (int, error)
	GetLeaderboard(ctx context.Context, limit, offset int) ([]LeaderboardEntry, error)
}

type PgRepository struct {
	db *sql.DB
}

func NewPgRepository(db *sql.DB) *PgRepository {
	return &PgRepository{db: db}
}

func (r *PgRepository) GetAnswersByAttempt(ctx context.Context, userID, attemptID string) ([]AnswerData, error) {
	const q = `
		SELECT ans.weight, ans.choice_score 
		FROM answers ans
		JOIN attempts att ON ans.attempt_id = att.id
		WHERE ans.attempt_id = $1 AND att.user_id = $2`

	rows, err := r.db.QueryContext(ctx, q, attemptID, userID)
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

func (r *PgRepository) GetUserStatsByRole(ctx context.Context, userID string) ([]RoleStats, error) {
	const q = `
		SELECT 
			sv.role,
			COUNT(a.id) FILTER (WHERE a.status = 'completed'),
			COUNT(a.id) FILTER (WHERE a.status = 'in_progress'),
			COUNT(a.id)
		FROM scenario_versions sv
		LEFT JOIN attempts a ON a.scenario_id = sv.id AND a.user_id = $1
		GROUP BY sv.role`

	rows, err := r.db.QueryContext(ctx, q, userID)
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

func (r *PgRepository) GetUserStatsByCategory(ctx context.Context, userID string) ([]CategoryStat, error) {
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
			WHERE a.status = 'completed' AND a.user_id = $1
		) as expanded_categories
		GROUP BY category
		ORDER BY count DESC`

	rows, err := r.db.QueryContext(ctx, q, userID)
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

func (r *PgRepository) GetUserTotalCompletedCount(ctx context.Context, userID string) (int, error) {
	const q = `SELECT COUNT(DISTINCT scenario_id) FROM attempts WHERE status = 'completed' AND user_id = $1`
	var total int
	err := r.db.QueryRowContext(ctx, q, userID).Scan(&total)
	return total, err
}

func (r *PgRepository) GetLeaderboard(ctx context.Context, limit, offset int) ([]LeaderboardEntry, error) {
	const q = `
		WITH current_scores AS (
			SELECT a.user_id, sv.logical_id, MAX(a.score) as max_score
			FROM attempts a
			JOIN scenario_versions sv ON a.scenario_id = sv.id
			WHERE a.status = 'completed'
			GROUP BY a.user_id, sv.logical_id
		),
		current_totals AS (
			SELECT user_id, SUM(max_score) as total_score
			FROM current_scores
			GROUP BY user_id
		),
		current_ranks AS (
			SELECT user_id, total_score,
				   DENSE_RANK() OVER (ORDER BY total_score DESC) as rnk
			FROM current_totals
		),
		prev_scores AS (
			SELECT a.user_id, sv.logical_id, MAX(a.score) as max_score
			FROM attempts a
			JOIN scenario_versions sv ON a.scenario_id = sv.id
			WHERE a.status = 'completed' AND a.completed_at < NOW() - INTERVAL '1 day'
			GROUP BY a.user_id, sv.logical_id
		),
		prev_totals AS (
			SELECT user_id, SUM(max_score) as total_score
			FROM prev_scores
			GROUP BY user_id
		),
		prev_ranks AS (
			SELECT user_id,
				   DENSE_RANK() OVER (ORDER BY total_score DESC) as rnk
			FROM prev_totals
		),
		user_fragments AS (
			SELECT user_id, COUNT(DISTINCT fragment_id) as fragments_count
			FROM user_inventory
			GROUP BY user_id
		)
		SELECT 
			c.rnk as current_rank,
			u.email,
			COALESCE(uf.fragments_count, 0) as fragments_count,
			c.total_score,
			p.rnk as prev_rank
		FROM current_ranks c
		JOIN users u ON c.user_id = u.id
		LEFT JOIN prev_ranks p ON c.user_id = p.user_id
		LEFT JOIN user_fragments uf ON c.user_id = uf.user_id
		ORDER BY c.rnk ASC
		LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []LeaderboardEntry
	for rows.Next() {
		var e LeaderboardEntry
		var prevRank sql.NullInt32

		if err := rows.Scan(&e.Rank, &e.Player, &e.Fragments, &e.Score, &prevRank); err != nil {
			return nil, err
		}

		if prevRank.Valid {
			change := int(prevRank.Int32) - e.Rank
			e.RankChange = &change
		} else {
			e.RankChange = nil
		}

		entries = append(entries, e)
	}

	return entries, rows.Err()
}
