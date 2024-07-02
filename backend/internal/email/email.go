package email

import (
	"fmt"
	"net/smtp"
)

type Email struct {
	email    string
	password string
	smtpHost string
	smtpPort string
}

// TODO: вынести пароль в переменную окружения

func New(email, password, smtpHost, smtpPort string) *Email {
	return &Email{
		email:    email,
		password: password,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
	}
}

func (e *Email) Send(receiverEmail, msg string) error {
	to := []string{receiverEmail}

	message := []byte(msg)

	auth := smtp.PlainAuth("", e.email, e.password, e.smtpHost)

	err := smtp.SendMail(e.smtpHost+":"+e.smtpPort, auth, e.email, to, message)
	if err != nil {
		fmt.Errorf("can't send email: %w", err)
	}

	return nil
}
