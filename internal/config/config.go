package config

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"path/filepath"
)

// Config ...
type Config struct {
	Env    string     `json:"env"`
	Format string     `json:"format"`
	Mail   MailConfig `json:"mail"`
	DB     DBConfig   `json:"db"`
}

// DBConfig ...
type DBConfig struct {
	Host     string `json:"host" env:"HOST" env-description:"sql server host"`
	Port     int    `json:"port" env:"PORT" env-default:"1433" env-description:"sql server port"`
	User     string `json:"user" env:"USER"`
	Password string `json:"password" env:"PASSWORD"`
	Database string `json:"database" env:"DATABASE"`
}

// MailConfig ...
type MailConfig struct {
	Server   string `json:"server" env:"MAIL_SERVER"`
	Port     int    `json:"port" env:"MAIL_PORT"`
	UseTLS   bool   `json:"useTLS" env:"MAIL_USE_TLS"`
	UseSSL   bool   `json:"useSSL" env:"MAIL_USE_SSL"`
	UserName string `json:"userName" env:"MAIL_USERNAME"`
	Password string `json:"password" env:"MAIL_PASSWORD"`
	// ToSendEmail адрес (куда отправлять письмо)
	ToSendEmail string `json:"toSendEmail" env:"TO_SEND_EMAIL"`
	// IsSendEmail признак для отправки по почте
	IsSendEmail bool   `json:"isSendEmail" env:"IS_SEND_EMAIL"`
	ContentType string `json:"contentType" env-default:"text/html"` //text/html  text/plain
}

// MustConfig ...
func MustConfig(fileConf string) *Config {
	if _, err := os.Stat(fileConf); errors.Is(err, os.ErrNotExist) {
		workingDir, _ := os.Executable()
		fileConf = filepath.Join(filepath.Dir(workingDir), filepath.Base(fileConf))
	}

	var cfg Config

	// переменные в .env файле в любом случае в приоритете
	if err := cleanenv.ReadConfig(fileConf, &cfg); err != nil {
		log.Fatal(err)
	}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal(err)
	}

	return &cfg
}
