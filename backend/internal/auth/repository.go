package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrSessionNotFound = errors.New("session not found")
	ErrTokenNotFound   = errors.New("reset token not found")
)

type UserRepository interface {
	Create(ctx context.Context, name, email, passwordHash string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByID(ctx context.Context, id string) (User, error)
	UpdatePassword(ctx context.Context, userID, passwordHash string) error
}

type SessionRepository interface {
	Create(ctx context.Context, userID, userAgent, ip string, ttl time.Duration) (Session, error)
	GetByID(ctx context.Context, id string) (Session, error)
	Delete(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
}

func (r *PgUserRepository) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	const q = `UPDATE users SET password_hash = $1, updated_at = now() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, q, passwordHash, userID)
	return err
}

type PgUserRepository struct {
	db *sql.DB
}

func NewPgUserRepository(db *sql.DB) *PgUserRepository {
	return &PgUserRepository{db: db}
}

func (r *PgUserRepository) Create(ctx context.Context, name, email, passwordHash string) (User, error) {
	const q = `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, password_hash, created_at, updated_at
	`
	var u User
	err := r.db.QueryRowContext(ctx, q, name, email, passwordHash).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return User{}, ErrUserExists
		}
		return User{}, err
	}
	return u, nil
}

func (r *PgUserRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	const q = `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE lower(email) = lower($1)
	`
	var u User
	err := r.db.QueryRowContext(ctx, q, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (r *PgUserRepository) GetByID(ctx context.Context, id string) (User, error) {
	const q = `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var u User
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, err
	}
	return u, nil
}

type PgSessionRepository struct {
	db *sql.DB
}

func NewPgSessionRepository(db *sql.DB) *PgSessionRepository {
	return &PgSessionRepository{db: db}
}

func (r *PgSessionRepository) Create(ctx context.Context, userID, userAgent, ip string, ttl time.Duration) (Session, error) {
	const q = `
		INSERT INTO sessions (user_id, user_agent, ip, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, user_agent, ip, expires_at, created_at
	`
	var s Session
	expiresAt := time.Now().UTC().Add(ttl)
	err := r.db.QueryRowContext(ctx, q, userID, userAgent, ip, expiresAt).Scan(
		&s.ID, &s.UserID, &s.UserAgent, &s.IP, &s.ExpiresAt, &s.CreatedAt,
	)
	if err != nil {
		return Session{}, err
	}
	return s, nil
}

func (r *PgSessionRepository) GetByID(ctx context.Context, id string) (Session, error) {
	const q = `
		SELECT id, user_id, user_agent, ip, expires_at, created_at
		FROM sessions
		WHERE id = $1
	`
	var s Session
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&s.ID, &s.UserID, &s.UserAgent, &s.IP, &s.ExpiresAt, &s.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Session{}, ErrSessionNotFound
	}
	if err != nil {
		return Session{}, err
	}
	return s, nil
}

func (r *PgSessionRepository) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM sessions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}

func (r *PgSessionRepository) DeleteExpired(ctx context.Context) error {
	const q = `DELETE FROM sessions WHERE expires_at <= now()`
	_, err := r.db.ExecContext(ctx, q)
	return err
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

type PasswordResetRepository interface {
	Create(ctx context.Context, userID, tokenHash string, ttl time.Duration) (PasswordResetToken, error)
	GetByTokenHash(ctx context.Context, tokenHash string) (PasswordResetToken, error)
	MarkUsed(ctx context.Context, id string) error
}

type PgPasswordResetRepository struct {
	db *sql.DB
}

func NewPgPasswordResetRepository(db *sql.DB) *PgPasswordResetRepository {
	return &PgPasswordResetRepository{db: db}
}

func (r *PgPasswordResetRepository) Create(ctx context.Context, userID, tokenHash string, ttl time.Duration) (PasswordResetToken, error) {
	const q = `
		INSERT INTO password_reset_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, token_hash, expires_at, used_at, created_at
	`
	var t PasswordResetToken
	err := r.db.QueryRowContext(ctx, q, userID, tokenHash, time.Now().UTC().Add(ttl)).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.UsedAt, &t.CreatedAt,
	)
	return t, err
}

func (r *PgPasswordResetRepository) GetByTokenHash(ctx context.Context, tokenHash string) (PasswordResetToken, error) {
	const q = `
		SELECT id, user_id, token_hash, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE token_hash = $1
	`
	var t PasswordResetToken
	err := r.db.QueryRowContext(ctx, q, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.UsedAt, &t.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return PasswordResetToken{}, ErrTokenNotFound
	}
	return t, err
}

func (r *PgPasswordResetRepository) MarkUsed(ctx context.Context, id string) error {
	const q = `
		UPDATE password_reset_tokens
		SET used_at = now()
		WHERE id = $1
		  AND used_at IS NULL
		  AND expires_at > now()
	`

	result, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrTokenNotFound
	}

	return nil
}
