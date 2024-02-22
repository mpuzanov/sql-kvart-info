package storage

import (
	"context"
	"fmt"
	"kvart-info/internal/config"
	"net/url"

	"github.com/jmoiron/sqlx"
)

// Storage ...
type Storage struct {
	db  *sqlx.DB
	ctx context.Context
	cfg *config.Config
}

// NewDB Создание подключения к БД
func NewDB(ctx context.Context, cfg *config.Config) (*Storage, error) {

	if cfg.DB.Host == "" || cfg.DB.Port == 0 || cfg.DB.Database == "" || cfg.DB.User == "" || cfg.DB.Password == "" {
		return nil, fmt.Errorf("не заполнены входные параметры БД")
	}
	dbDSN := NewURLConnectionString(cfg.DB.Host, cfg.DB.Port, cfg.DB.Database, cfg.DB.User, cfg.DB.Password)
	//log.Printf("db_dsn: %s", db_dsn)
	db, err := sqlx.Open("sqlserver", dbDSN)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Storage{ctx: ctx, db: db, cfg: cfg}, nil
}

// NewURLConnectionString "sqlserver://user:password@host:port?database=database_name"
func NewURLConnectionString(host string, port int, database, user, password string) string {
	v := url.Values{}
	v.Set("database", database)
	//v.Add("app name", "MyAppName")

	var u = url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(user, password),
		Host:     fmt.Sprintf("%s:%d", host, port),
		RawQuery: v.Encode(),
	}
	return u.String()
}
