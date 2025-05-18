package mailer

import (
	"context"
	"fmt"
	"net/smtp"
	"weather-forecast/internal/infrastructure/logger"
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

func NewSMTPMailer(from, host, port, username, password string, logger logger.Logger) *SMTPMailer {
	mailer := SMTPMailer{
		from:     from,
		host:     host,
		port:     port,
		username: username,
		password: password,
		logger:   logger,
	}
	mailer.auth = smtp.PlainAuth("", mailer.username, mailer.password, mailer.host)
	return &mailer
}

func (m *SMTPMailer) Send(_ context.Context, subject string, body, email string) {
	msg := []byte(
		fmt.Sprintf("To: %s\r\n", email) +
			fmt.Sprintf("From: %s\r\n", m.from) +
			fmt.Sprintf("Subject: %s\r\n", subject) +
			"\r\n" + body,
	)

	err := smtp.SendMail(m.host+":"+m.port, m.auth, m.from, []string{email}, msg)
	if err != nil {
		m.logger.Warnf("Failed to send email with subject %s to %s. Due to error: %s", subject, email, err.Error())
	} else {
		m.logger.Infof("Email sent successfully to %s", email)
	}
}
