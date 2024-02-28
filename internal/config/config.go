package config

import (
	"errors"
	"fmt"
	"kvart-info/pkg/email"
	"kvart-info/pkg/logging"
	"kvart-info/pkg/mssql"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config ...
type Config struct {
	Env  string       `yaml:"env"`
	DB   mssql.Config `yaml:"db"`
	Mail email.Config `yaml:"mail"`
	// ToSendEmail адрес (куда отправлять письмо)
	ToSendEmail string `yaml:"to_send_email" env:"TO_SEND_EMAIL"`
	// IsSendEmail признак для отправки по почте
	IsSendEmail bool `yaml:"is_send_email" env:"IS_SEND_EMAIL"`
}

// NewConfig returns app config.
func NewConfig(fileConf string) (*Config, error) {
	cfg := &Config{}

	if _, err := os.Stat(fileConf); errors.Is(err, os.ErrNotExist) {
		workingDir, _ := os.Executable()
		fileConf = filepath.Join(filepath.Dir(workingDir), filepath.Base(fileConf))
	}

	err := cleanenv.ReadConfig(fileConf, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// LogValue для удобного логирования структуры
func (cfg Config) LogValue() logging.Value {
	return logging.GroupValue(
		logging.StringAttr("Env", cfg.Env),
		logging.StringAttr("ToSendEmail", cfg.ToSendEmail),
		logging.BoolAttr("IsSendEmail", cfg.IsSendEmail),

		logging.Group(
			"db",
			logging.StringAttr("host", cfg.DB.Host),
			logging.IntAttr("port", cfg.DB.Port),
			logging.StringAttr("user", cfg.DB.User),
			logging.StringAttr("password", "<REMOVED>"),
			logging.StringAttr("database", cfg.DB.Database),
		),

		logging.Group(
			"mail",
			logging.StringAttr("server", cfg.Mail.Server),
			logging.IntAttr("port", cfg.Mail.Port),
			logging.StringAttr("username", cfg.Mail.UserName),
			logging.StringAttr("password", "<REMOVED>"),
			logging.BoolAttr("use_tls", cfg.Mail.UseTLS),
			logging.BoolAttr("use_ssl", cfg.Mail.UseSSL),
		),
	)
}
