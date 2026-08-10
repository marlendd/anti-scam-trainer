package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const resendAPIURL = "https://api.resend.com/emails"

type Config struct {
	APIKey string
	From   string
}

type Mailer struct {
	cfg    Config
	client *http.Client
}

func New(cfg Config) *Mailer {
	return &Mailer{
		cfg: cfg,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
	Text    string   `json:"text"`
}

func (m *Mailer) SendPasswordReset(ctx context.Context, toEmail, resetLink string) error {
	subject := "Восстановление пароля — Anti-Scam Trainer"

	text := fmt.Sprintf(
		"Здравствуйте!\n\n"+
			"Вы запросили восстановление пароля. Перейдите по ссылке ниже, чтобы задать новый пароль:\n\n"+
			"%s\n\n"+
			"Ссылка действительна 1 час. Если вы не запрашивали восстановление — просто проигнорируйте это письмо.",
		resetLink,
	)

	html := fmt.Sprintf(
		`<p>Здравствуйте!</p>
<p>Вы запросили восстановление пароля. Перейдите по ссылке ниже, чтобы задать новый пароль:</p>
<p><a href="%s">%s</a></p>
<p>Ссылка действительна 1 час. Если вы не запрашивали восстановление — просто проигнорируйте это письмо.</p>`,
		resetLink, resetLink,
	)

	payload := resendRequest{
		From:    m.cfg.From,
		To:      []string{toEmail},
		Subject: subject,
		HTML:    html,
		Text:    text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal resend payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendAPIURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to build resend request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+m.cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("resend request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resend returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
