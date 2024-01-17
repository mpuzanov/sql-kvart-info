package email

import (
	"errors"
	"kvart-info/internal/config"
	"os"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
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
func (s *AppEmail) Send(bodyMessage, subject, toAddress string, file string, statusEmail chan string) {

	m := gomail.NewMessage()
	m.SetHeader("From", s.cfg.Mail.UserName)
	m.SetHeader("To", strings.Split(toAddress, ";")...)
	m.SetHeader("Subject", subject)
	var contentType = "text/html"
	if s.cfg.Mail.ContentType != "" {
		contentType = s.cfg.Mail.ContentType
	}
	m.SetBody(contentType, bodyMessage)

	if file != "" {
		if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
			statusEmail <- "File " + file + "for send email not found"
			return
		}
		m.Attach(file)
	}

	d := gomail.NewDialer(s.cfg.Mail.Server, s.cfg.Mail.Port, s.cfg.Mail.UserName, s.cfg.Mail.Password)

	var err error
	for i := 0; i < 3; i++ {
		if err = d.DialAndSend(m); err == nil {
			statusEmail <- "Email sent to " + toAddress + " completed"
			return
		}
		time.Sleep(5 * time.Second) // Retry after 5 seconds
	}
	statusEmail <- "Failed to send email to " + toAddress + ": " + err.Error()

}
