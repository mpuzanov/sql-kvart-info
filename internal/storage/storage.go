package storage

import (
	"context"
	"database/sql"
	"fmt"
	"kvart-info/internal/config"
	"kvart-info/internal/model"
	"kvart-info/pkg/logging"
	"net/url"

	"github.com/jmoiron/sqlx"
)

// Mssql ...
type Mssql struct {
	db  *sqlx.DB
	ctx context.Context
	cfg *config.Config
}

// NewDB Создание подключения к БД
func NewDB(ctx context.Context, cfg *config.Config) (*Mssql, error) {

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
	return &Mssql{ctx: ctx, db: db, cfg: cfg}, nil
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

// GetTotalData получаем сводную информацию из БД
func (s *Mssql) GetTotalData() ([]model.TotalData, error) {
	l := logging.LoggerFromContext(s.ctx)
	l.Info("Executing query", "database", s.cfg.DB.Database)

	var data []model.TotalData
	stmt, err := s.db.PrepareNamedContext(s.ctx, QueryGetTotal)
	if err != nil {
		return nil, fmt.Errorf("failed PrepareNamedContext total: %w", err)
	}
	err = stmt.SelectContext(s.ctx, &data, map[string]interface{}{})
	if err == sql.ErrNoRows {
		return data, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed SelectContext total: %w", err)
	}

	return data, nil
}
