package mssql

import (
	"fmt"
	"net/url"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Config .
type Config struct {
	Host         string `yaml:"host" env:"DB_HOST" env-required:"true"`
	Port         int    `yaml:"port" env:"DB_PORT" env-default:"1433" env-description:"sql server port"`
	User         string `yaml:"user" env:"DB_USER" env-required:"true"`
	Password     string `yaml:"password" env:"DB_PASSWORD" env-required:"true"`
	Database     string `yaml:"database" env:"DB_DATABASE" env-required:"true"`
	TimeoutQuery int    `yaml:"timeout_query" env:"TIMEOUT_QUERY" env-default:"300"` // Second
	APPName      string `yaml:"app_name" env:"APP_NAME"`
}

// MSSQL ...
type MSSQL struct {
	DB *sqlx.DB
}

// New Создание подключения к БД
func New(cfg *Config) (*MSSQL, error) {

	if cfg.Host == "" || cfg.Port == 0 || cfg.Database == "" || cfg.User == "" || cfg.Password == "" {
		return nil, fmt.Errorf("не заполнены входные параметры БД")
	}
	dsn := cfg.getDatabaseURL()
	db, err := sqlx.Open("sqlserver", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "sqlx.Open")
	}
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "sqlx.Ping")
	}
	return &MSSQL{DB: db}, nil
}

// Close -.
func (r *MSSQL) Close() {
	if r.DB != nil {
		r.DB.Close()
	}
}

// getDatabaseURL "sqlserver://user:password@host:port?database=database_name"
func (d Config) getDatabaseURL() string {
	v := url.Values{}
	v.Set("database", d.Database)
	if d.APPName != "" {
		v.Add("app name", d.APPName)
	}

	var u = url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(d.User, d.Password),
		Host:     fmt.Sprintf("%s:%d", d.Host, d.Port),
		RawQuery: v.Encode(),
	}
	return u.String()
}
