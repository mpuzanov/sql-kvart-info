package config

import (
	"errors"
	"fmt"
	"kvart-info/pkg/email"
	"log/slog"
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
func (cfg Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("Env", cfg.Env),
		slog.String("ToSendEmail", cfg.ToSendEmail),
		slog.Bool("IsSendEmail", cfg.IsSendEmail),

		slog.Group(
			"db",
			slog.String("host", cfg.DB.Host),
			slog.Int("port", cfg.DB.Port),
			slog.String("user", cfg.DB.User),
			slog.String("password", "<REMOVED>"),
			slog.String("database", cfg.DB.Database),
		),

		slog.Group(
			"mail",
			slog.String("server", cfg.Mail.Server),
			slog.Int("port", cfg.Mail.Port),
			slog.String("username", cfg.Mail.UserName),
			slog.String("password", "<REMOVED>"),
			slog.Bool("use_tls", cfg.Mail.UseTLS),
			slog.Bool("use_ssl", cfg.Mail.UseSSL),
		),
	)
}
