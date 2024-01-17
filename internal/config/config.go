package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config ...
type Config struct {
	Env  string     `yaml:"env"`
	DB   DBConfig   `yaml:"db"`
	Mail MailConfig `yaml:"mail"`	
}

// DBConfig ...
type DBConfig struct {
	Host     string `yaml:"host" env:"HOST" env-description:"sql server host"`
	Port     int    `yaml:"port" env:"PORT" env-default:"1433" env-description:"sql server port"`
	User     string `yaml:"user" env:"USER"`
	Password string `yaml:"password" env:"PASSWORD"`
	Database string `yaml:"database" env:"DATABASE"`
}

// MailConfig ...
type MailConfig struct {
	Server   string `yaml:"server" env:"MAIL_SERVER"`
	Port     int    `yaml:"port" env:"MAIL_PORT"`
	UseTLS   bool   `yaml:"use_tls" env:"MAIL_USE_TLS"`
	UseSSL   bool   `yaml:"use_ssl" env:"MAIL_USE_SSL"`
	UserName string `yaml:"username" env:"MAIL_USERNAME"`
	Password string `yaml:"password" env:"MAIL_PASSWORD"`
	// ToSendEmail адрес (куда отправлять письмо)
	ToSendEmail string `yaml:"toSendEmail" env:"TO_SEND_EMAIL"`
	// IsSendEmail признак для отправки по почте
	IsSendEmail bool   `yaml:"isSendEmail" env:"IS_SEND_EMAIL"`
	ContentType string `yaml:"contentType" env-default:"text/html"` //text/html  text/plain
	Timeout     int    `yaml:"timeout" env-default:"10"`
}

// MustConfig загрузка файла конфига
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
