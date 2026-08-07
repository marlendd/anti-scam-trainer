diff --git a/backend/cmd/api/main.go b/backend/cmd/api/main.go
index 69f0d62..dc5f0e5 100644
--- a/backend/cmd/api/main.go
+++ b/backend/cmd/api/main.go
@@ -13,6 +13,7 @@ import (
 	"time"
 
 	"github.com/marlendd/anti-scam-trainer/internal/auth"
+	"github.com/marlendd/anti-scam-trainer/internal/evaluation"
 	"github.com/marlendd/anti-scam-trainer/internal/platform/config"
 	"github.com/marlendd/anti-scam-trainer/internal/platform/mailer"
 	"github.com/marlendd/anti-scam-trainer/internal/platform/postgres"
@@ -66,6 +67,10 @@ func run(cfg *config.Config, log *slog.Logger) error {
 	authService := auth.NewService(userRepo, sessionRepo, passwordResetRepo, m, cfg.AppBaseURL)
 	authHandler := auth.NewHandler(authService, log, cfg.SecureCookies)
 	requireAuth := auth.RequireAuth(authService, log)
+	// evaluation
+	evalRepo := evaluation.NewPgRepository(db)
+	evalService := evaluation.NewService(evalRepo)
+	evalHandler := evaluation.NewHandler(evalService, log)
 
 	mux := http.NewServeMux()
 
@@ -83,7 +88,8 @@ func run(cfg *config.Config, log *slog.Logger) error {
 	mux.HandleFunc("POST /api/v1/auth/logout", authHandler.Logout)
 	mux.HandleFunc("POST /api/v1/auth/forgot-password", authHandler.ForgotPassword)
 	mux.HandleFunc("POST /api/v1/auth/reset-password", authHandler.ResetPassword)
-
+	mux.HandleFunc("GET  /api/v1/attempts/{id}/result", evalHandler.GetStatsOfAttempt)
+	mux.HandleFunc("GET  /api/v1/profile/progress", evalHandler.GetGlobalStatsHandler)
 	// protected routes
 	mux.Handle("GET /api/v1/users/me", requireAuth(http.HandlerFunc(authHandler.Me)))
 
diff --git a/backend/internal/evaluation/handlers.go b/backend/internal/evaluation/handlers.go
new file mode 100644
index 0000000..844069d
--- /dev/null
+++ b/backend/internal/evaluation/handlers.go
@@ -0,0 +1,55 @@
+package evaluation
+
+import (
+	"encoding/json"
+	"log/slog"
+	"net/http"
+)
+
+type Handler struct {
+	service *Service
+	log     *slog.Logger
+}
+
+func NewHandler(service *Service, log *slog.Logger) *Handler {
+	return &Handler{service: service, log: log}
+}
+
+func (h *Handler) GetStatsOfAttempt(w http.ResponseWriter, r *http.Request) {
+	attemptID := r.PathValue("id")
+
+	res, err := h.service.GetAttemptResults(r.Context(), attemptID)
+	if err != nil {
+		h.log.Error("failed to get attempt stats", "attempt_id", attemptID, "error", err)
+
+		h.respondError(w, http.StatusInternalServerError, "database error")
+		return
+	}
+
+	w.Header().Set("Content-Type", "application/json")
+
+	h.respondJSON(w, http.StatusOK, map[string]int{"score": res})
+}
+
+func (h *Handler) GetGlobalStatsHandler(w http.ResponseWriter, r *http.Request) {
+	ctx := r.Context()
+
+	stats, err := h.service.GetGlobalStats(ctx)
+	if err != nil {
+		h.log.Error("failed to get stats", "error", err)
+		h.respondError(w, http.StatusInternalServerError, "internal server error")
+		return
+	}
+
+	h.respondJSON(w, http.StatusOK, stats)
+}
+
+func (h *Handler) respondJSON(w http.ResponseWriter, status int, payload any) {
+	w.Header().Set("Content-Type", "application/json")
+	w.WriteHeader(status)
+	_ = json.NewEncoder(w).Encode(payload)
+}
+
+func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
+	h.respondJSON(w, status, map[string]string{"error": message})
+}
diff --git a/backend/internal/evaluation/models.go b/backend/internal/evaluation/models.go
new file mode 100644
index 0000000..f51acea
--- /dev/null
+++ b/backend/internal/evaluation/models.go
@@ -0,0 +1,13 @@
+package evaluation
+
+type AnswerData struct {
+	Weight      int16 `json:"weight"`
+	ChoiceScore int16 `json:"choice_score"`
+}
+
+type RoleStats struct {
+	Role            string `json:"role"`
+	CompletedCount  int64  `json:"completed_count"`
+	InProgressCount int64  `json:"in_progress_count"`
+	TotalStarted    int64  `json:"total_started"`
+}
diff --git a/backend/internal/evaluation/repository.go b/backend/internal/evaluation/repository.go
new file mode 100644
index 0000000..0962e7c
--- /dev/null
+++ b/backend/internal/evaluation/repository.go
@@ -0,0 +1,67 @@
+package evaluation
+
+import (
+	"context"
+	"database/sql"
+)
+
+type Repository interface {
+	GetAnswersByAttempt(ctx context.Context, attemptID string) ([]AnswerData, error)
+	GetStatsByRole(ctx context.Context) ([]RoleStats, error)
+}
+
+type PgRepository struct {
+	db *sql.DB
+}
+
+func NewPgRepository(db *sql.DB) *PgRepository {
+	return &PgRepository{db: db}
+}
+
+func (r *PgRepository) GetAnswersByAttempt(ctx context.Context, attemptID string) ([]AnswerData, error) {
+	const q = `SELECT weight, choice_score FROM answers WHERE attempt_id = $1`
+
+	rows, err := r.db.QueryContext(ctx, q, attemptID)
+	if err != nil {
+		return nil, err
+	}
+	defer rows.Close()
+
+	var res []AnswerData
+	for rows.Next() {
+		var a AnswerData
+		if err := rows.Scan(&a.Weight, &a.ChoiceScore); err != nil {
+			return nil, err
+		}
+		res = append(res, a)
+	}
+	return res, rows.Err()
+}
+
+func (r *PgRepository) GetStatsByRole(ctx context.Context) ([]RoleStats, error) {
+	const q = `
+		SELECT 
+			sv.role,
+			COUNT(*) FILTER (WHERE a.status = 'completed'),
+			COUNT(*) FILTER (WHERE a.status = 'in_progress'),
+			COUNT(*)
+		FROM attempts a
+		JOIN scenario_versions sv ON a.scenario_id = sv.id
+		GROUP BY sv.role`
+
+	rows, err := r.db.QueryContext(ctx, q)
+	if err != nil {
+		return nil, err
+	}
+	defer rows.Close()
+
+	var stats []RoleStats
+	for rows.Next() {
+		var s RoleStats
+		if err := rows.Scan(&s.Role, &s.CompletedCount, &s.InProgressCount, &s.TotalStarted); err != nil {
+			return nil, err
+		}
+		stats = append(stats, s)
+	}
+	return stats, rows.Err()
+}
diff --git a/backend/internal/evaluation/service.go b/backend/internal/evaluation/service.go
new file mode 100644
index 0000000..47e2192
--- /dev/null
+++ b/backend/internal/evaluation/service.go
@@ -0,0 +1,37 @@
+package evaluation
+
+import "context"
+
+type Service struct {
+	repo Repository
+}
+
+func NewService(repo Repository) *Service {
+	return &Service{repo: repo}
+}
+
+func (s *Service) GetAttemptResults(ctx context.Context, attemptID string) (int, error) {
+	dbRes, err := s.repo.GetAnswersByAttempt(ctx, attemptID)
+	if err != nil {
+		return -1, err
+	}
+
+	if len(dbRes) == 0 {
+		return 0, nil
+	}
+
+	var sum int
+	var weightSum int
+	for _, val := range dbRes {
+		sum += int(val.Weight * val.ChoiceScore)
+		weightSum += int(val.Weight)
+	}
+
+	res := sum / weightSum
+
+	return res, nil
+}
+
+func (s *Service) GetGlobalStats(ctx context.Context) ([]RoleStats, error) {
+	return s.repo.GetStatsByRole(ctx)
+}
diff --git a/backend/tests/integration/evaluatuion_test.go b/backend/tests/integration/evaluatuion_test.go
new file mode 100644
index 0000000..6f24262
--- /dev/null
+++ b/backend/tests/integration/evaluatuion_test.go
@@ -0,0 +1,71 @@
+package evaluation_test
+
+import (
+	"context"
+	"database/sql"
+	"os"
+	"testing"
+
+	_ "github.com/lib/pq"
+	"github.com/marlendd/anti-scam-trainer/internal/evaluation"
+	"github.com/stretchr/testify/require"
+)
+
+func setupTestDB(t *testing.T) *sql.DB {
+	dbURL := os.Getenv("DATABASE_URL")
+	if dbURL == "" {
+		dbURL = "postgres://postgres:password@localhost:5433/postgres?sslmode=disable"
+	}
+
+	db, err := sql.Open("postgres", dbURL)
+	require.NoError(t, err)
+
+	return db
+}
+
+func TestEvaluation_Integration(t *testing.T) {
+	db := setupTestDB(t)
+	defer db.Close()
+
+	repo := evaluation.NewPgRepository(db)
+	svc := evaluation.NewService(repo)
+	ctx := context.Background()
+
+	_, err := db.Exec("TRUNCATE users, scenario_versions, attempts, answers CASCADE")
+	require.NoError(t, err)
+
+	t.Run("Seed and Calculate Score", func(t *testing.T) {
+		seedSQL := `
+			INSERT INTO users (id, email, password_hash) VALUES ('00000000-0000-0000-0000-000000000001', 'test@test.com', 'hash');
+			INSERT INTO scenario_versions (id, logical_id, version, role, title, description, content) 
+			VALUES ('00000000-0000-0000-0000-000000000002', gen_random_uuid(), 1, 'buyer', 'title', 'desc', '{}'::jsonb);
+			INSERT INTO attempts (id, user_id, scenario_id, status, current_node_id) 
+			VALUES ('120b7935-62bf-4fd8-828a-6bbe7ef7a19a', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', 'in_progress', 'start_node');
+			
+			-- Тут оценка 50, значит риск ОБЯЗАТЕЛЕН
+			INSERT INTO answers (attempt_id, node_id, choice_id, idempotency_key, weight, choice_score, risk_categories, consequence, explanation, response)
+			VALUES ('120b7935-62bf-4fd8-828a-6bbe7ef7a19a', 'node1', 'c1', gen_random_uuid(), 2, 50, '["suspicious_link"]'::jsonb, 'cons', 'expl', '{}');
+			
+			-- Тут оценка 100, риск может быть пустым
+			INSERT INTO answers (attempt_id, node_id, choice_id, idempotency_key, weight, choice_score, risk_categories, consequence, explanation, response)
+			VALUES ('120b7935-62bf-4fd8-828a-6bbe7ef7a19a', 'node2', 'c2', gen_random_uuid(), 1, 100, '[]'::jsonb, 'cons', 'expl', '{}');
+		`
+		_, err := db.Exec(seedSQL)
+		require.NoError(t, err)
+
+		score, err := svc.GetAttemptResults(ctx, "120b7935-62bf-4fd8-828a-6bbe7ef7a19a")
+		require.NoError(t, err)
+		require.Equal(t, 66, score)
+	})
+
+	t.Run("Verify Global Progress Stats", func(t *testing.T) {
+		stats, err := repo.GetStatsByRole(ctx)
+		require.NoError(t, err)
+		require.NotEmpty(t, stats)
+		require.Len(t, stats, 1)
+
+		require.Equal(t, "buyer", stats[0].Role)
+		require.Equal(t, int64(1), stats[0].InProgressCount)
+		require.Equal(t, int64(0), stats[0].CompletedCount)
+	})
+}
diff --git a/test-sql-data-set/data.sql b/test-sql-data-set/data.sql
new file mode 100644
index 0000000..e9b67d4
--- /dev/null
+++ b/test-sql-data-set/data.sql
@@ -0,0 +1,87 @@
+INSERT INTO users (
+	email,
+	password_hash
+) VALUES (
+	'piter@gmail.com',
+	'fsdf98789s0-sdf789798s0df'
+);
+
+INSERT INTO scenario_versions(
+	logical_id,
+	version,
+	role,
+	title,
+	description,
+	is_active,
+	content
+) VALUES (
+	gen_random_uuid(),
+	1,
+	'buyer',
+	'buying videocard',                   
+	'scam when buying',                     
+	true,
+	'{}'::jsonb  
+);
+
+INSERT INTO attempts (
+	user_id,
+	scenario_id,
+	status,
+	current_node_id
+) VALUES (
+	(SELECT id FROM users ORDER BY created_at DESC LIMIT 1),
+	(SELECT id FROM scenario_versions ORDER BY created_at DESC LIMIT 1),
+	'in_progress',
+	'start_node'
+);
+
+INSERT INTO answers (
+	attempt_id,
+	node_id,
+	choice_id,
+	idempotency_key,
+	weight,
+	choice_score,
+	risk_categories,
+	consequence,
+	explanation,
+	response
+) VALUES (
+	(SELECT id FROM attempts ORDER BY started_at DESC LIMIT 1), 
+	'start_dialog',
+	'get_suspicious link',
+	gen_random_uuid(), 
+	2,
+	50,
+	'["Not check url"]'::jsonb,
+	'lost a prime',
+	'You shouldnt click on suspicious links.',
+	'{"emotion": "happy"}'::jsonb
+);
+
+INSERT INTO answers (
+	attempt_id,
+	node_id,
+	choice_id,
+	idempotency_key,
+	weight,
+	choice_score,
+	risk_categories,
+	consequence,
+	explanation,
+	response
+) VALUES (
+	(SELECT id FROM attempts ORDER BY started_at DESC LIMIT 1), 
+	'check_url_node',           
+	'ignore_link',              
+	gen_random_uuid(),          
+	1,                          
+	100,                        
+	'[]'::jsonb,                
+	'safely avoided the scam',
+	'User successfully recognized a suspicious URL.',
+	'{"emotion": "relieved"}'::jsonb
+);
+
+SELECT * FROM answers;
\ No newline at end of file
