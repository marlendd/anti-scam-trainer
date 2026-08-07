package mailer

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type Mailer struct {
	cfg Config
}

func New(cfg Config) *Mailer {
	return &Mailer{cfg: cfg}
}

func (m *Mailer) SendPasswordReset(toEmail, resetLink string) error {
	subject := "Восстановление пароля — Anti-Scam Trainer"
	body := fmt.Sprintf(
		"Здравствуйте!\r\n\r\n"+
			"Вы запросили восстановление пароля. Перейдите по ссылке ниже, чтобы задать новый пароль:\r\n\r\n"+
			"%s\r\n\r\n"+
			"Ссылка действительна 1 час. Если вы не запрашивали восстановление — просто проигнорируйте это письмо.",
		resetLink,
	)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		m.cfg.From, toEmail, subject, body,
	)

	return m.sendViaImplicitTLS(toEmail, []byte(msg))
}

// sendViaImplicitTLS отправляет письмо через SMTPS (порт 465) —
// в отличие от STARTTLS (587), тут TLS устанавливается сразу при подключении,
// именно так требует Яндекс для порта 465.
func (m *Mailer) sendViaImplicitTLS(toEmail string, msg []byte) error {
	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)

	tlsConfig := &tls.Config{
		ServerName: m.cfg.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("tls dial failed: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, m.cfg.Host)
	if err != nil {
		return fmt.Errorf("smtp client init failed: %w", err)
	}
	defer client.Close()

	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth failed: %w", err)
	}

	if err := client.Mail(m.cfg.From); err != nil {
		return fmt.Errorf("smtp MAIL FROM failed: %w", err)
	}
	if err := client.Rcpt(toEmail); err != nil {
		return fmt.Errorf("smtp RCPT TO failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA failed: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("smtp write body failed: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close writer failed: %w", err)
	}

	return client.Quit()
}
