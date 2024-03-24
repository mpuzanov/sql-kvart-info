package config

import (
	"errors"
	"fmt"
	"kvart-info/pkg/email"
	"kvart-info/pkg/wslog"
	"os"
	"path/filepath"

	"github.com/mpuzanov/dbwrap"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config ...
type Config struct {
	Env  string        `yaml:"env"`
	DB   dbwrap.Config `yaml:"db"`
	Mail email.Config  `yaml:"mail"`
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
func (cfg Config) LogValue() wslog.Value {
	return wslog.GroupValue(
		wslog.StrAttr("Env", cfg.Env),
		wslog.StrAttr("ToSendEmail", cfg.ToSendEmail),
		wslog.BoolAttr("IsSendEmail", cfg.IsSendEmail),

		wslog.Group(
			"db",
			wslog.StrAttr("host", cfg.DB.Host),
			wslog.IntAttr("port", cfg.DB.Port),
			wslog.StrAttr("user", cfg.DB.User),
			wslog.StrAttr("password", "<REMOVED>"),
			wslog.StrAttr("database", cfg.DB.Database),
		),

		wslog.Group(
			"mail",
			wslog.StrAttr("server", cfg.Mail.Server),
			wslog.IntAttr("port", cfg.Mail.Port),
			wslog.StrAttr("username", cfg.Mail.UserName),
			wslog.StrAttr("password", "<REMOVED>"),
			wslog.BoolAttr("use_tls", cfg.Mail.UseTLS),
			wslog.BoolAttr("use_ssl", cfg.Mail.UseSSL),
		),
	)
}
