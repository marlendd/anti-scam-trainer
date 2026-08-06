package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrSessionNotFound = errors.New("session not found")
)

type UserRepository interface {
	Create(ctx context.Context, email, passwordHash string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByID(ctx context.Context, id string) (User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, userID, userAgent, ip string, ttl time.Duration) (Session, error)
	GetByID(ctx context.Context, id string) (Session, error)
	Delete(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
}

type PgUserRepository struct {
	db *sql.DB
}

func NewPgUserRepository(db *sql.DB) *PgUserRepository {
	return &PgUserRepository{db: db}
}

func (r *PgUserRepository) Create(ctx context.Context, email, passwordHash string) (User, error) {
	const q = `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, password_hash, created_at, updated_at
	`
	var u User
	err := r.db.QueryRowContext(ctx, q, email, passwordHash).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
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
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE lower(email) = lower($1)
	`
	var u User
	err := r.db.QueryRowContext(ctx, q, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
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
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var u User
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt,
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
	expiresAt := time.Now().Add(ttl)
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
	const q = `DELETE FROM sessions WHERE expires_at < now()`
	_, err := r.db.ExecContext(ctx, q)
	return err
}

func isUniqueViolation(err error) bool {
	return err != nil && errContains(err.Error(), "23505")
}

func errContains(s, substr string) bool {
	return len(s) >= len(substr) && (func() bool {
		for i := 0; i+len(substr) <= len(s); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
		return false
	})()
}
