package email

import (
	"errors"
	"os"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

// MailConfig ...
type MailConfig struct {
	Server           string `yaml:"server" env:"mail_server"`
	Port             int    `yaml:"port" env:"mail_port"`
	UseTLS           bool   `yaml:"use_tls" env:"mail_use_tls" env-default:"false"`
	UseSSL           bool   `yaml:"use_ssl" env:"mail_use_ssl" env-default:"true"`
	UserName         string `yaml:"username" env:"mail_username"`
	Password         string `yaml:"password" env:"mail_password"`
	Timeout          int    `yaml:"timeout" env:"mail_timeout" env-default:"10"`
	CountRetryFail   int    `yaml:"count_retry_fail" env:"mail_count_retry_fail"  env-default:"3"`
	TimeoutRetryFail int    `yaml:"timeout_count_retry_fail" env:"mail_count_retry_fail"  env-default:"5"`
}

// AppEmail ...
type AppEmail MailConfig

// New создание объекта для отправки по почте
func New(cfg MailConfig) *AppEmail {
	v := AppEmail(cfg)
	return &v
}

// Send Отправка письма по email
func (cfg *AppEmail) Send(bodyMessage, subject, toAddress string, file string) (string, error) {
	var err error

	m := gomail.NewMessage()
	m.SetHeader("From", cfg.UserName)
	m.SetHeader("To", strings.Split(toAddress, ";")...)
	m.SetHeader("Subject", subject)

	var contentType = "text/plain"
	if strings.Contains(bodyMessage, "<body>") {
		contentType = "text/html"
	}
	m.SetBody(contentType, bodyMessage)

	if file != "" {
		if _, err = os.Stat(file); errors.Is(err, os.ErrNotExist) {
			return "file " + file + "for send email not found", err
		}
		m.Attach(file)
	}

	d := gomail.NewDialer(cfg.Server, cfg.Port, cfg.UserName, cfg.Password)

	for i := 0; i < cfg.CountRetryFail; i++ {
		if err = d.DialAndSend(m); err == nil {
			return "email sent to " + toAddress + " completed", nil
		}
		time.Sleep(time.Duration(cfg.TimeoutRetryFail) * time.Second) // Retry after 5 seconds
	}
	return "failed to send email to " + toAddress + ": " + err.Error(), err

}
