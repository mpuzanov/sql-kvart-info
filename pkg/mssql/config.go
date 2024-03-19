package mssql

import (
	"fmt"
	"net/url"
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

// NewConfig создание конфига по умолчанию
func NewConfig() *Config {
	return &Config{Host: "localhost",
		Port:         1433,
		User:         "sa",
		Database:     "master",
		TimeoutQuery: 300}
}

// WithPassword установка пароля
func (c *Config) WithPassword(pwd string) *Config {
	c.Password = pwd
	return c
}

// WithDB установка БД
func (c *Config) WithDB(dbname string) *Config {
	c.Database = dbname
	return c
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
