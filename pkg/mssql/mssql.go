package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"errors"

	"github.com/jmoiron/sqlx"
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
	DB  *sqlx.DB
	Cfg *Config
}

var ErrBadConfigDB = errors.New("не заполнены входные параметры БД")

// New Создание подключения к БД
func New(cfg *Config) (*MSSQL, error) {

	if cfg.Host == "" || cfg.Port == 0 || cfg.Database == "" || cfg.User == "" || cfg.Password == "" {
		return nil, ErrBadConfigDB
	}
	dsn := cfg.getDatabaseURL()
	db, err := sqlx.Open("sqlserver", dsn)
	if err != nil {
		return nil, fmt.Errorf(" sqlx.Open : %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf(" sqlx.Ping : %w", err)
	}
	return &MSSQL{DB: db, Cfg: cfg}, nil
}

// Close -.
func (ms *MSSQL) Close() {
	if ms.DB != nil {
		ms.DB.Close()
	}
}

// SetTimeout ...
func (ms *MSSQL) SetTimeout(timeout int) {
	ms.Cfg.TimeoutQuery = timeout
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

// Get получаем данные из запроса
func (ms *MSSQL) GetSelect(query string, dest any, params map[string]interface{}) error {

	// ограничим время выполнения запроса
	dur := time.Duration(ms.Cfg.TimeoutQuery) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	stmt, err := ms.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		return fmt.Errorf(" [GetSelect] PrepareNamedContext : %w", err)
	}
	err = stmt.SelectContext(ctx, dest, params)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return fmt.Errorf(" [GetSelect] SelectContext : %w", err)
	}

	return nil
}
