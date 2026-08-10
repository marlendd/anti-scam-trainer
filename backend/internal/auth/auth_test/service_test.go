package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/marlendd/anti-scam-trainer/internal/auth"
	"github.com/marlendd/anti-scam-trainer/internal/platform/mailer"
	"github.com/stretchr/testify/require"
)

// ---------- mocks ----------

type userRepositoryStub struct {
	createFn         func(ctx context.Context, name, email, passwordHash string) (auth.User, error)
	getByEmailFn     func(ctx context.Context, email string) (auth.User, error)
	getByIDFn        func(ctx context.Context, id string) (auth.User, error)
	updatePasswordFn func(ctx context.Context, userID, passwordHash string) error
}

func (s *userRepositoryStub) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	if s.updatePasswordFn == nil {
		return nil
	}
	return s.updatePasswordFn(ctx, userID, passwordHash)
}

func (s *userRepositoryStub) Create(ctx context.Context, name, email, passwordHash string) (auth.User, error) {
	if s.createFn == nil {
		panic("unexpected Create call")
	}
	return s.createFn(ctx, name, email, passwordHash)
}

func (s *userRepositoryStub) GetByEmail(ctx context.Context, email string) (auth.User, error) {
	if s.getByEmailFn == nil {
		panic("unexpected GetByEmail call")
	}
	return s.getByEmailFn(ctx, email)
}

func (s *userRepositoryStub) GetByID(ctx context.Context, id string) (auth.User, error) {
	if s.getByIDFn == nil {
		panic("unexpected GetByID call")
	}
	return s.getByIDFn(ctx, id)
}

type sessionRepositoryStub struct {
	createFn         func(ctx context.Context, userID, userAgent, ip string, ttl time.Duration) (auth.Session, error)
	getByIDFn        func(ctx context.Context, id string) (auth.Session, error)
	deleteFn         func(ctx context.Context, id string) error
	deleteCalledWith []string
}

func (s *sessionRepositoryStub) Create(ctx context.Context, userID, userAgent, ip string, ttl time.Duration) (auth.Session, error) {
	if s.createFn == nil {
		panic("unexpected Create call")
	}
	return s.createFn(ctx, userID, userAgent, ip, ttl)
}

func (s *sessionRepositoryStub) GetByID(ctx context.Context, id string) (auth.Session, error) {
	if s.getByIDFn == nil {
		panic("unexpected GetByID call")
	}
	return s.getByIDFn(ctx, id)
}

func (s *sessionRepositoryStub) Delete(ctx context.Context, id string) error {
	s.deleteCalledWith = append(s.deleteCalledWith, id)
	if s.deleteFn == nil {
		return nil
	}
	return s.deleteFn(ctx, id)
}

func (s *sessionRepositoryStub) DeleteExpired(ctx context.Context) error {
	return nil
}

// ---------- Register ----------

func TestService_Register(t *testing.T) {
	t.Parallel()

	t.Run("successful registration hashes password before storing", func(t *testing.T) {
		t.Parallel()

		var capturedHash string

		users := &userRepositoryStub{
			createFn: func(_ context.Context, name, email, passwordHash string) (auth.User, error) {
				require.Equal(t, "Test User", name)
				require.Equal(t, "test@example.com", email)

				capturedHash = passwordHash

				return auth.User{
					ID:           "user-1",
					Name:         name,
					Email:        email,
					PasswordHash: passwordHash,
				}, nil
			},
		}

		sessions := &sessionRepositoryStub{}

		service := newTestService(users, sessions)

		u, err := service.Register(
			context.Background(),
			"Test User",
			"test@example.com",
			"password123",
		)

		require.NoError(t, err)
		require.Equal(t, "user-1", u.ID)
		require.Equal(t, "Test User", u.Name)
		require.NotEqual(t, "password123", capturedHash, "password must be hashed, not stored in plaintext")
		require.NoError(t, auth.CheckPassword(capturedHash, "password123"))
	})

	t.Run("returns ErrEmailTaken when repository reports duplicate", func(t *testing.T) {
		t.Parallel()

		users := &userRepositoryStub{
			createFn: func(context.Context, string, string, string) (auth.User, error) {
				return auth.User{}, auth.ErrUserExists
			},
		}

		sessions := &sessionRepositoryStub{}

		service := newTestService(users, sessions)

		_, err := service.Register(
			context.Background(),
			"Test User",
			"taken@example.com",
			"password123",
		)

		require.ErrorIs(t, err, auth.ErrEmailTaken)
	})

	t.Run("propagates unexpected repository error", func(t *testing.T) {
		t.Parallel()

		repoErr := errors.New("db connection lost")

		users := &userRepositoryStub{
			createFn: func(context.Context, string, string, string) (auth.User, error) {
				return auth.User{}, repoErr
			},
		}

		sessions := &sessionRepositoryStub{}

		service := newTestService(users, sessions)

		_, err := service.Register(
			context.Background(),
			"Test User",
			"test@example.com",
			"password123",
		)

		require.ErrorIs(t, err, repoErr)
	})
}

// ---------- Login ----------

func TestService_Login(t *testing.T) {
	t.Parallel()

	t.Run("successful login creates a session", func(t *testing.T) {
		t.Parallel()

		hash, err := auth.HashPassword("password123")
		require.NoError(t, err)

		existingUser := auth.User{
			ID:           "user-1",
			Email:        "test@example.com",
			PasswordHash: hash,
		}

		users := &userRepositoryStub{
			getByEmailFn: func(_ context.Context, email string) (auth.User, error) {
				require.Equal(t, "test@example.com", email)
				return existingUser, nil
			},
		}

		var capturedUserID string

		sessions := &sessionRepositoryStub{
			createFn: func(_ context.Context, userID, userAgent, ip string, ttl time.Duration) (auth.Session, error) {
				capturedUserID = userID

				require.Equal(t, auth.SessionTTL, ttl)

				return auth.Session{
					ID:        "session-1",
					UserID:    userID,
					ExpiresAt: time.Now().Add(ttl),
				}, nil
			},
		}

		service := newTestService(users, sessions)

		u, sess, err := service.Login(
			context.Background(),
			"test@example.com",
			"password123",
			"test-agent",
			"127.0.0.1",
		)

		require.NoError(t, err)
		require.Equal(t, existingUser.ID, u.ID)
		require.Equal(t, "session-1", sess.ID)
		require.Equal(t, existingUser.ID, capturedUserID)
	})

	t.Run("returns ErrInvalidCredentials when user not found", func(t *testing.T) {
		t.Parallel()

		users := &userRepositoryStub{
			getByEmailFn: func(context.Context, string) (auth.User, error) {
				return auth.User{}, auth.ErrUserNotFound
			},
		}

		sessions := &sessionRepositoryStub{}

		service := newTestService(users, sessions)

		_, _, err := service.Login(
			context.Background(),
			"nobody@example.com",
			"password123",
			"",
			"",
		)

		require.ErrorIs(t, err, auth.ErrInvalidCredentials)
	})

	t.Run("returns ErrInvalidCredentials on wrong password without creating session", func(t *testing.T) {
		t.Parallel()

		hash, err := auth.HashPassword("correct-password")
		require.NoError(t, err)

		users := &userRepositoryStub{
			getByEmailFn: func(context.Context, string) (auth.User, error) {
				return auth.User{
					ID:           "user-1",
					PasswordHash: hash,
				}, nil
			},
		}

		sessions := &sessionRepositoryStub{
			createFn: func(context.Context, string, string, string, time.Duration) (auth.Session, error) {
				t.Fatal("session must not be created on failed login")
				return auth.Session{}, nil
			},
		}

		service := newTestService(users, sessions)

		_, _, err = service.Login(
			context.Background(),
			"test@example.com",
			"wrong-password",
			"",
			"",
		)

		require.ErrorIs(t, err, auth.ErrInvalidCredentials)
	})

	t.Run("propagates unexpected repository error", func(t *testing.T) {
		t.Parallel()

		repoErr := errors.New("db timeout")

		users := &userRepositoryStub{
			getByEmailFn: func(context.Context, string) (auth.User, error) {
				return auth.User{}, repoErr
			},
		}

		sessions := &sessionRepositoryStub{}

		service := newTestService(users, sessions)

		_, _, err := service.Login(
			context.Background(),
			"test@example.com",
			"password123",
			"",
			"",
		)

		require.ErrorIs(t, err, repoErr)
	})
}

// ---------- Logout ----------

func TestService_Logout(t *testing.T) {
	t.Parallel()

	sessions := &sessionRepositoryStub{}
	service := newTestService(&userRepositoryStub{}, sessions)

	err := service.Logout(context.Background(), "session-123")

	require.NoError(t, err)
	require.Equal(t, []string{"session-123"}, sessions.deleteCalledWith)
}

// ---------- GetUserBySession ----------

func TestService_GetUserBySession(t *testing.T) {
	t.Parallel()

	t.Run("returns user for active session", func(t *testing.T) {
		t.Parallel()

		expectedUser := auth.User{
			ID:    "user-1",
			Email: "test@example.com",
		}

		sessions := &sessionRepositoryStub{
			getByIDFn: func(_ context.Context, id string) (auth.Session, error) {
				require.Equal(t, "session-1", id)

				return auth.Session{
					ID:        "session-1",
					UserID:    "user-1",
					ExpiresAt: time.Now().Add(time.Hour),
				}, nil
			},
		}

		users := &userRepositoryStub{
			getByIDFn: func(_ context.Context, id string) (auth.User, error) {
				require.Equal(t, "user-1", id)
				return expectedUser, nil
			},
		}

		service := newTestService(users, sessions)

		u, err := service.GetUserBySession(context.Background(), "session-1")

		require.NoError(t, err)
		require.Equal(t, expectedUser, u)
		require.Empty(t, sessions.deleteCalledWith, "active session must not be deleted")
	})

	t.Run("expired session is deleted and returns ErrSessionNotFound", func(t *testing.T) {
		t.Parallel()

		sessions := &sessionRepositoryStub{
			getByIDFn: func(context.Context, string) (auth.Session, error) {
				return auth.Session{
					ID:        "session-1",
					UserID:    "user-1",
					ExpiresAt: time.Now().Add(-time.Hour),
				}, nil
			},
		}

		users := &userRepositoryStub{}

		service := newTestService(users, sessions)

		_, err := service.GetUserBySession(context.Background(), "session-1")

		require.ErrorIs(t, err, auth.ErrSessionNotFound)
		require.Equal(
			t,
			[]string{"session-1"},
			sessions.deleteCalledWith,
			"expired session must be cleaned up",
		)
	})

	t.Run("returns ErrSessionNotFound when session does not exist", func(t *testing.T) {
		t.Parallel()

		sessions := &sessionRepositoryStub{
			getByIDFn: func(context.Context, string) (auth.Session, error) {
				return auth.Session{}, auth.ErrSessionNotFound
			},
		}

		service := newTestService(&userRepositoryStub{}, sessions)

		_, err := service.GetUserBySession(
			context.Background(),
			"missing-session",
		)

		require.ErrorIs(t, err, auth.ErrSessionNotFound)
	})
}

type passwordResetRepositoryStub struct{}

func (s *passwordResetRepositoryStub) Create(
	context.Context,
	string,
	string,
	time.Duration,
) (auth.PasswordResetToken, error) {
	return auth.PasswordResetToken{}, nil
}

func (s *passwordResetRepositoryStub) GetByTokenHash(
	context.Context,
	string,
) (auth.PasswordResetToken, error) {
	return auth.PasswordResetToken{}, nil
}

func (s *passwordResetRepositoryStub) MarkUsed(
	context.Context,
	string,
) error {
	return nil
}

func newTestService(
	users auth.UserRepository,
	sessions auth.SessionRepository,
) *auth.Service {
	m := mailer.New(mailer.Config{
		APIKey: "test-api-key",
		From:   "test@test.com",
	})

	return auth.NewService(
		users,
		sessions,
		&passwordResetRepositoryStub{},
		m,
		"http://localhost:3000",
	)
}
