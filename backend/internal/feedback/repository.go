package feedback

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
)

type Repository interface {
	GetAttemptDataForFeedback(ctx context.Context, userID, attemptID string) (*PromptData, error)
}

type PgRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewPgRepository(db *sql.DB, log *slog.Logger) *PgRepository {
	return &PgRepository{db: db, log: log}
}

func (r *PgRepository) GetAttemptDataForFeedback(ctx context.Context, userID, attemptID string) (*PromptData, error) {
	const q = `
		SELECT 
			sv.role, sv.title, sv.description, a.score,
			ans.choice_score, ans.risk_categories, ans.consequence, ans.explanation, ans.response
		FROM attempts a
		JOIN scenario_versions sv ON a.scenario_id = sv.id
		JOIN answers ans ON ans.attempt_id = a.id
		WHERE a.id = $1 AND a.user_id = $2
		ORDER BY ans.created_at ASC`

	rows, err := r.db.QueryContext(ctx, q, attemptID, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			r.log.Error("failed to close rows", "error", err)
		}
	}()

	var promptData PromptData
	var initialized bool
	stepNumber := 1

	for rows.Next() {
		var (
			role, title, description      string
			score                         sql.NullInt32
			choiceScore                   int
			consequence, explanation      string
			riskCategoriesBytes, resBytes []byte
		)

		err := rows.Scan(
			&role, &title, &description, &score,
			&choiceScore, &riskCategoriesBytes, &consequence, &explanation, &resBytes,
		)
		if err != nil {
			return nil, err
		}

		if !initialized {
			promptData.Role = role
			promptData.ScenarioTitle = title
			promptData.ScenarioDescription = description
			if score.Valid {
				promptData.TotalScore = int(score.Int32)
			}
			initialized = true
		}

		var risks []string
		_ = json.Unmarshal(riskCategoriesBytes, &risks)

		var rawRes AnswerResponseRaw
		_ = json.Unmarshal(resBytes, &rawRes)

		risksStr := "Нет"
		if len(risks) > 0 {
			risksStr = ""
			for i, risk := range risks {
				if i > 0 {
					risksStr += ", "
				}
				risksStr += risk
			}
		}

		promptData.Answers = append(promptData.Answers, AnswerData{
			StepNumber:     stepNumber,
			NodeQuestion:   rawRes.NodeQuestion,
			ChoiceText:     rawRes.ChoiceText,
			ChoiceScore:    choiceScore,
			RiskCategories: risksStr,
			Consequence:    consequence,
			Explanation:    explanation,
		})
		stepNumber++
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if !initialized {
		return nil, sql.ErrNoRows
	}

	return &promptData, nil
}
