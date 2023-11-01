package config

import (
	"encoding/json"
	"os"
)

// Config ...
type Config struct {
	Log  LogConfig  `json:"logging"`
	Mail MailConfig `json:"mail"`
	DB   DBConfig   `json:"db"`
}

// LogConfig ...
type LogConfig struct {
	Level   string `json:"level"`
	UseJSON bool   `json:"use_json"`
}

// DBConfig ...
type DBConfig struct {
	Host     string `json:"host" env:"HOST" env-description:"sql server host"`
	Port     int    `json:"port" env:"PORT" envDefault:"1433" env-description:"sql server port"`
	User     string `json:"user" env:"USER"`
	Password string `json:"password" env:"PASSWORD"`
	Database string `json:"database" env:"DATABASE"`
}

// MailConfig ...
type MailConfig struct {
	Server      string `json:"server" env:"MAIL_SERVER"`
	Port        int    `json:"port" env:"MAIL_PORT"`
	UseTLS      bool   `json:"useTLS" env:"MAIL_USE_TLS"`
	UseSSL      bool   `json:"useSSL" env:"MAIL_USE_SSL"`
	UserName    string `json:"userName" env:"MAIL_USERNAME"`
	Password    string `json:"password" env:"MAIL_PASSWORD"`
	ToSendEmail string `json:"toSendEmail" env:"TO_SEND_EMAIL"`
	IsSendEmail bool   `json:"isSendEmail" env:"IS_SEND_EMAIL"`
}

// LoadConfig ...
func LoadConfig(pathFile string) (*Config, error) {
	file, err := os.Open(pathFile)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	config := new(Config)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
