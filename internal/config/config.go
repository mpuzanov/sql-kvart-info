package config

import (
	"errors"
	"fmt"
	"kvart-info/pkg/email"
	"os"
	"path/filepath"

	"github.com/mpuzanov/wslog"

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
		wslog.String("Env", cfg.Env),
		wslog.String("ToSendEmail", cfg.ToSendEmail),
		wslog.Bool("IsSendEmail", cfg.IsSendEmail),

		wslog.Group(
			"db",
			wslog.String("host", cfg.DB.Host),
			wslog.Int("port", cfg.DB.Port),
			wslog.String("user", cfg.DB.User),
			wslog.String("password", "<REMOVED>"),
			wslog.String("database", cfg.DB.Database),
		),

		wslog.Group(
			"mail",
			wslog.String("server", cfg.Mail.Server),
			wslog.Int("port", cfg.Mail.Port),
			wslog.String("username", cfg.Mail.UserName),
			wslog.String("password", "<REMOVED>"),
			wslog.Bool("use_tls", cfg.Mail.UseTLS),
			wslog.Bool("use_ssl", cfg.Mail.UseSSL),
		),
	)
}
