package mssql

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"
)

func sqlErr(err error, query string, args ...interface{}) error {
	return fmt.Errorf(`run query "%s" with args %+v: %w`, query, args, err)
}

func namedQuery(query string, arg interface{}) (nq string, args []interface{}, err error) {
	nq, args, err = sqlx.Named(query, arg)
	if err != nil {
		return "", nil, sqlErr(err, query, args...)
	}
	return nq, args, nil
}

// Exec Выполнение запроса DML
func (ms *MSSQL) Exec(query string, args ...interface{}) (int64, error) {

	// ограничим время выполнения запроса
	dur := time.Duration(ms.Cfg.TimeoutQuery) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	result, err := ms.DBX.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, sqlErr(err, query, args...)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// NamedExec Выполнение запроса DML
func (ms *MSSQL) NamedExec(query string, arg interface{}) (int64, error) {

	nq, args, err := namedQuery(query, arg)
	if err != nil {
		return 0, err
	}

	return ms.Exec(ms.DBX.Rebind(nq), args...)
}

// Select получаем данные из запроса в слайс структур
//
// var users []User
//
// err := ts.db.Select(&users, "select * from users")
func (ms *MSSQL) Select(dest interface{}, query string, args ...interface{}) error {

	// ограничим время выполнения запроса
	dur := time.Duration(ms.Cfg.TimeoutQuery) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	if err := sqlx.SelectContext(ctx, ms.DBX, dest, query, args...); err != nil {
		return sqlErr(err, query, args...)
	}

	return nil
}

// NamedSelect получаем данные из запроса в слайс структур
//
// var users []User
//
// err := ts.db.NamedSelect(&users, "select * from users where name=:Name", map[string]interface{}{"Name": "admin"})
func (ms *MSSQL) NamedSelect(dest interface{}, query string, arg interface{}) error {

	nq, args, err := namedQuery(query, arg)
	if err != nil {
		return err
	}

	return ms.Select(dest, ms.DBX.Rebind(nq), args...)
}

// SelectMaps ...
func (ms *MSSQL) SelectMaps(query string, args ...interface{}) (ret []map[string]interface{}, err error) {

	dur := time.Duration(ms.Cfg.TimeoutQuery) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	rows, err := ms.DBX.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, sqlErr(err, query, args...)
	}

	defer func() {
		err = multierr.Combine(err, rows.Close())
	}()

	ret = []map[string]interface{}{}
	numCols := -1
	for rows.Next() {
		var m map[string]interface{}
		if numCols < 0 {
			m = map[string]interface{}{}
		} else {
			m = make(map[string]interface{}, numCols)
		}

		if err = rows.MapScan(m); err != nil {
			return nil, sqlErr(err, query, args...)
		}
		ret = append(ret, m)
		numCols = len(m)
	}

	if err = rows.Err(); err != nil {
		return nil, sqlErr(err, query, args...)
	}

	return ret, nil
}

// NamedSelectMaps ...
func (ms *MSSQL) NamedSelectMaps(query string, arg interface{}) (ret []map[string]interface{}, err error) {
	nq, args, err := namedQuery(query, arg)
	if err != nil {
		return nil, err
	}

	return ms.SelectMaps(ms.DBX.Rebind(nq), args...)
}
