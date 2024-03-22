package dbwrap

import (
	"fmt"

	"errors"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

// DBSQL ...
type DBSQL struct {
	DBX *sqlx.DB
	Cfg *Config
}

// ErrBadConfigDB ошибка
var ErrBadConfigDB = errors.New("не заполнены параметры подключения к БД")

// New Создание подключения к БД
func New(cfg *Config) (*DBSQL, error) {

	if cfg.Host == "" || cfg.Port == 0 || cfg.Database == "" || cfg.User == "" || cfg.Password == "" {
		return nil, ErrBadConfigDB
	}
	driverName := "sqlserver"
	dsn := cfg.getDatabaseURL(driverName)
	db, err := sqlx.Connect(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf(" sqlx.Connect : %w", err)
	}
	return &DBSQL{DBX: db, Cfg: cfg}, nil
}

// Close -.
func (d *DBSQL) Close() error {
	return d.DBX.Close()
}

// SetTimeout установка таймаута для выполнения запроса
func (d *DBSQL) SetTimeout(timeout uint) {
	d.Cfg.TimeoutQuery = int(timeout)
}
