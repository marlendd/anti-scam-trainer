package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/marlendd/anti-scam-trainer/internal/platform/mailer"
)

const SessionTTL = 30 * 24 * time.Hour // 30 дней

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already registered")
	ErrTokenExpired       = errors.New("reset token expired")
	ErrTokenAlreadyUsed   = errors.New("reset token already used")
)

func (s *Service) Register(ctx context.Context, name, email, password string) (User, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return User{}, err
	}

	u, err := s.users.Create(ctx, name, email, hash)
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

const PasswordResetTTL = 1 * time.Hour

type Service struct {
	users         UserRepository
	sessions      SessionRepository
	passwordReset PasswordResetRepository
	mailer        *mailer.Mailer
	appBaseURL    string
}

func NewService(users UserRepository, sessions SessionRepository, passwordReset PasswordResetRepository, m *mailer.Mailer, appBaseURL string) *Service {
	return &Service{
		users:         users,
		sessions:      sessions,
		passwordReset: passwordReset,
		mailer:        m,
		appBaseURL:    appBaseURL,
	}
}

func (s *Service) RequestPasswordReset(ctx context.Context, email string) error {
	u, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil // молча выходим, как будто всё ок
		}
		return err
	}

	rawToken, err := generateRandomToken()
	if err != nil {
		return err
	}
	tokenHash := hashToken(rawToken)

	if _, err := s.passwordReset.Create(ctx, u.ID, tokenHash, PasswordResetTTL); err != nil {
		return err
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.appBaseURL, rawToken)

	if err := s.mailer.SendPasswordReset(ctx, u.Email, resetLink); err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

// ResetPassword проверяет токен и устанавливает новый пароль.
func (s *Service) ResetPassword(ctx context.Context, rawToken, newPassword string) error {
	tokenHash := hashToken(rawToken)

	t, err := s.passwordReset.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return err // ErrTokenNotFound наружу
	}

	if t.UsedAt != nil {
		return ErrTokenAlreadyUsed
	}
	if time.Now().After(t.ExpiresAt) {
		return ErrTokenExpired
	}

	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.users.UpdatePassword(ctx, t.UserID, hash); err != nil {
		return err
	}

	if err := s.passwordReset.MarkUsed(ctx, t.ID); err != nil {
		return err
	}

	return nil
}

func generateRandomToken() (string, error) {
	b := make([]byte, 32) // 256 бит энтропии
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
