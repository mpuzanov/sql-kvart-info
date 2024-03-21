package mssql

import (
	"fmt"

	"errors"

	"github.com/jmoiron/sqlx"
)

// MSSQL ...
type MSSQL struct {
	DBX *sqlx.DB
	Cfg *Config
}

// ErrBadConfigDB ошибка
var ErrBadConfigDB = errors.New("не заполнены параметры подключения к БД")

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
	return &MSSQL{DBX: db, Cfg: cfg}, nil
}

// Close -.
func (ms *MSSQL) Close() error {
	if ms.DBX != nil {
		return ms.DBX.Close()
	}
	return nil
}

// SetTimeout установка таймаута для выполнения запроса
func (ms *MSSQL) SetTimeout(timeout uint) {
	ms.Cfg.TimeoutQuery = int(timeout)
}
