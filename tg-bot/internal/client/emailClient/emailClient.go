package emailClient

import (
	"fmt"
	"net/smtp"
)

type EmailClient struct {
	email    string
	password string
	smtpHost string
	smtpPort string
}

func New(email, password, smtpHost, smtpPort string) *EmailClient {
	return &EmailClient{
		email:    email,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

func (e *EmailClient) Send(receiverEmail, msg string) error {
	to := []string{receiverEmail}

	message := []byte(msg)

	auth := smtp.PlainAuth("", e.email, e.password, e.smtpHost)

	err := smtp.SendMail(e.smtpHost+":"+e.smtpPort, auth, e.email, to, message)
	if err != nil {
		return fmt.Errorf("can't send email: %w", err)
	}

	return nil
}
