package email

import (
	"gopkg.in/gomail.v2"
	"kvart-info/internal/config"
	"log/slog"
	"strings"
)

type AppEmail struct {
	cfg *config.Config
}

func New(cfg *config.Config) *AppEmail {
	return &AppEmail{cfg: cfg}
}

// SendText Отправка сообщения по почте
func (s *AppEmail) SendText(bodyMessage, subject, address string) error {
	slog.Info("Send text", "email", address)

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.Mail.UserName)
	m.SetHeader("To", strings.Split(address, ";")...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", bodyMessage)

	d := gomail.NewDialer(s.cfg.Mail.Server, s.cfg.Mail.Port, s.cfg.Mail.UserName, s.cfg.Mail.Password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	slog.Info("Sending completed")
	return nil
}

// SendFile Отправка по почте файла
func (s *AppEmail) SendFile(bodyMessage, subject, address string, file string) error {
	if file != "" {
		slog.Info("Send file", slog.String("file", file), slog.String("email", address))
	}
	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.Mail.UserName)
	m.SetHeader("To", strings.Split(address, ";")...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", bodyMessage)
	m.Attach(file)

	d := gomail.NewDialer(s.cfg.Mail.Server, s.cfg.Mail.Port, s.cfg.Mail.UserName, s.cfg.Mail.Password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	slog.Info("Sending completed")
	return nil
}
