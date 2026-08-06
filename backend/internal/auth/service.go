package auth

import (
	"context"
	"errors"
	"net/http"
	"time"
)

const SessionTTL = 30 * 24 * time.Hour // 30 дней

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already registered")
)

type Service struct {
	users    UserRepository
	sessions SessionRepository
}

func NewService(users UserRepository, sessions SessionRepository) *Service {
	return &Service{users: users, sessions: sessions}
}

func (s *Service) Register(ctx context.Context, email, password string) (User, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return User{}, err
	}

	u, err := s.users.Create(ctx, email, hash)
	if err != nil {
		if errors.Is(err, ErrUserExists) {
			return User{}, ErrEmailTaken
		}
		return User{}, err
	}
	return u, nil
}

func (s *Service) Login(ctx context.Context, email, password, userAgent, ip string) (User, Session, error) {
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return User{}, Session{}, ErrInvalidCredentials
		}
		return User{}, Session{}, err
	}

	if err := CheckPassword(u.PasswordHash, password); err != nil {
		return User{}, Session{}, ErrInvalidCredentials
	}

	sess, err := s.sessions.Create(ctx, u.ID, userAgent, ip, SessionTTL)
	if err != nil {
		return User{}, Session{}, err
	}

	return u, sess, nil
}

func (s *Service) Logout(ctx context.Context, sessionID string) error {
	return s.sessions.Delete(ctx, sessionID)
}

func (s *Service) GetUserBySession(ctx context.Context, sessionID string) (User, error) {
	sess, err := s.sessions.GetByID(ctx, sessionID)
	if err != nil {
		return User{}, err
	}

	if time.Now().After(sess.ExpiresAt) {
		_ = s.sessions.Delete(ctx, sessionID)
		return User{}, ErrSessionNotFound
	}

	return s.users.GetByID(ctx, sess.UserID)
}

func clientIP(r *http.Request) string {
	ip := r.RemoteAddr
	return ip
}
