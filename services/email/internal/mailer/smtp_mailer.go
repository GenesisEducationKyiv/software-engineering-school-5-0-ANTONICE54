package mailer

import (
	"context"
	"email-service/internal/config"
	"fmt"
	"net/smtp"
	"weather-forecast/pkg/logger"
)

type SMTPMailer struct {
	from     string
	host     string
	port     string
	username string
	password string
	auth     smtp.Auth
	logger   logger.Logger
}

func NewSMTPMailer(cfg *config.Config, logger logger.Logger) *SMTPMailer {
	mailer := SMTPMailer{
		from:     cfg.MailerFrom,
		host:     cfg.MailerHost,
		port:     cfg.MailerPort,
		username: cfg.MailerUsername,
		password: cfg.MailerPassword,
		logger:   logger,
	}
	mailer.auth = smtp.PlainAuth("", mailer.username, mailer.password, mailer.host)
	return &mailer
}

func (m *SMTPMailer) Send(_ context.Context, subject string, body, email string) error {
	msg := []byte(
		fmt.Sprintf("To: %s\r\n", email) +
			fmt.Sprintf("From: %s\r\n", m.from) +
			fmt.Sprintf("Subject: %s\r\n", subject) +
			"\r\n" + body,
	)

	return smtp.SendMail(m.host+":"+m.port, m.auth, m.from, []string{email}, msg)

}
