package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/mail"
	"strings"
)

type Handler struct {
	service *Service
	log     *slog.Logger
	// secureCookies — true в проде (HTTPS), false для локальной разработки
	secureCookies bool
}

func NewHandler(service *Service, log *slog.Logger, secureCookies bool) *Handler {
	return &Handler{service: service, log: log, secureCookies: secureCookies}
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "name is required")
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if !isValidEmail(req.Email) {
		respondError(w, http.StatusBadRequest, "invalid email")
		return
	}
	if len(req.Password) < 8 {
		respondError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	u, err := h.service.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, ErrEmailTaken) {
			respondError(w, http.StatusConflict, "email already registered")
			return
		}
		h.log.Error("register failed", "error", err)
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusCreated, u.ToPublic())
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userAgent := r.UserAgent()
	ip := clientIP(r)

	u, sess, err := h.service.Login(r.Context(), req.Email, req.Password, userAgent, ip)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			respondError(w, http.StatusUnauthorized, "invalid email or password")
			return
		}
		h.log.Error("login failed", "error", err)
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	setSessionCookie(w, sess, h.secureCookies)
	respondJSON(w, http.StatusOK, u.ToPublic())
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(SessionCookieName)
	if err == nil {
		if err := h.service.Logout(r.Context(), cookie.Value); err != nil {
			h.log.Error("logout failed", "error", err)
		}
	}

	clearSessionCookie(w, h.secureCookies)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	u, err := h.service.users.GetByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			respondError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		h.log.Error("get me failed", "error", err)
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, u.ToPublic())
}

func setSessionCookie(w http.ResponseWriter, sess Session, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sess.ID,
		Path:     "/",
		Expires:  sess.ExpiresAt,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearSessionCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

type resetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req forgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if !isValidEmail(req.Email) {
		respondError(w, http.StatusBadRequest, "invalid email")
		return
	}

	if err := h.service.RequestPasswordReset(r.Context(), req.Email); err != nil {
		h.log.Error("request password reset failed", "error", err)
		// пользователю всё равно отвечаем 200 — не раскрываем детали
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "если email зарегистрирован, на него отправлена ссылка для восстановления",
	})
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req resetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.NewPassword) < 8 {
		respondError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	err := h.service.ResetPassword(r.Context(), req.Token, req.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, ErrTokenNotFound):
			respondError(w, http.StatusBadRequest, "invalid or expired token")
		case errors.Is(err, ErrTokenExpired):
			respondError(w, http.StatusBadRequest, "token expired")
		case errors.Is(err, ErrTokenAlreadyUsed):
			respondError(w, http.StatusBadRequest, "token already used")
		default:
			h.log.Error("reset password failed", "error", err)
			respondError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "password updated successfully"})
}
