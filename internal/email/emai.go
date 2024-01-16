package email

import (
	"gopkg.in/gomail.v2"
	"kvart-info/internal/config"
	"log/slog"
	"strings"
)

// AppEmail ...
type AppEmail struct {
	cfg *config.Config
}

// New создание объекта для отправки по почте
func New(cfg *config.Config) *AppEmail {
	return &AppEmail{cfg: cfg}
}

// Send Отправка письма по email
func (s *AppEmail) Send(bodyMessage, subject, address string, file string) error {
	slog.Info("Send", slog.String("email", address), slog.String("file", file))

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.Mail.UserName)
	m.SetHeader("To", strings.Split(address, ";")...)
	m.SetHeader("Subject", subject)
	var contentType = "text/html"
	if s.cfg.Mail.ContentType != "" {
		contentType = s.cfg.Mail.ContentType
	}
	m.SetBody(contentType, bodyMessage)

	if file != "" {
		m.Attach(file)
	}

	d := gomail.NewDialer(s.cfg.Mail.Server, s.cfg.Mail.Port, s.cfg.Mail.UserName, s.cfg.Mail.Password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	slog.Info("Sending email completed")
	return nil
}
