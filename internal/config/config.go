package config

import (
	"errors"
	"kvart-info/pkg/email"
	"log"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config ...
type Config struct {
	Env  string           `yaml:"env"`
	DB   DBConfig         `yaml:"db"`
	Mail email.MailConfig `yaml:"mail"`
	// ToSendEmail адрес (куда отправлять письмо)
	ToSendEmail string `yaml:"to_send_email" env:"TO_SEND_EMAIL"`
	// IsSendEmail признак для отправки по почте
	IsSendEmail bool `yaml:"is_send_email" env:"IS_SEND_EMAIL"`
}

// DBConfig ...
type DBConfig struct {
	Host     string `yaml:"host" env:"HOST" env-description:"sql server host"`
	Port     int    `yaml:"port" env:"PORT" env-default:"1433" env-description:"sql server port"`
	User     string `yaml:"user" env:"USER"`
	Password string `yaml:"password" env:"PASSWORD"`
	Database string `yaml:"database" env:"DATABASE"`
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
