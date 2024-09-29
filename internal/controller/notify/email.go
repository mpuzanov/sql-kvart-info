package notify

import (
	"kvart-info/pkg/email"
)

// Email для отправки по почте
type Email struct {
	cfgEmail    email.Config
	bodyMessage string
	subject     string
	toAddress   string
}

// NewEmail создание объекта для отправки по почте
func NewEmail(cfg email.Config, body, subject, toAddress string) *Email {
	return &Email{
		cfgEmail:    cfg,
		bodyMessage: body,
		subject:     subject,
		toAddress:   toAddress,
	}
}

// Send отправка по почте
func (n Email) Send() (string, error) {
	objEmail := n.cfgEmail
	return objEmail.Send(n.bodyMessage, n.subject, n.toAddress, "")
}
