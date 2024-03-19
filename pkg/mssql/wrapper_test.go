package mssql

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/stretchr/testify/suite"
)

// User Пользователь БД
type User struct {
	Name      string
	Email     string
	CreatedAt time.Time `db:"created_at"`
}

// Users список пользователей
type Users []User

// DBSuite структура для набора тестов с БД
type TestDBSuite struct {
	suite.Suite
	db *MSSQL // коннект к БД master
}

var (
	dbName    = "users_test"
	tableName = fmt.Sprintf("%s.dbo.users", dbName)
)

func TestTestDBSuite(t *testing.T) {
	suite.Run(t, &TestDBSuite{})
}

func (ts *TestDBSuite) SetupSuite() {

	config := NewConfig().WithPassword("123")
	db, err := New(config)
	if err != nil {
		ts.T().Fatalf("cannot connect to [master] : %v", err)
	}
	ts.db = db
	setupDatabase(ts)
}

func (its *TestDBSuite) TearDownSuite() {
	tearDownDatabase(its)
}

func setupDatabase(its *TestDBSuite) {
	its.T().Log("setting up database")

	_, err := its.db.DBX.Exec(fmt.Sprintf(`CREATE DATABASE %s`, dbName))
	if err != nil {
		its.FailNowf("unable to create database", err.Error())
	}
	its.T().Logf("База [%s] создана\n", dbName)

	_, err = its.db.DBX.Exec(fmt.Sprintf(`CREATE TABLE %s (
		[name] varchar(50) PRIMARY KEY,
		email varchar(100) UNIQUE NOT NULL,
		created_at datetime NOT NULL DEFAULT current_timestamp
	)`, tableName))
	if err != nil {
		its.FailNowf("unable to create table", err.Error())
	}
	its.T().Logf("Таблица [%s] создана\n", tableName)

}

func tearDownDatabase(its *TestDBSuite) {
	its.T().Log("tearing down database")

	_, err := its.db.DBX.Exec(fmt.Sprintf(`DROP TABLE %s`, tableName))
	if err != nil {
		its.FailNowf("unable to drop table", err.Error())
	}

	_, err = its.db.DBX.Exec(fmt.Sprintf(`DROP DATABASE %s`, dbName))
	if err != nil {
		its.FailNowf("unable to drop database", err.Error())
	}

	err = its.db.Close()
	if err != nil {
		its.FailNowf("unable to close database", err.Error())
	}
}

func (ts *TestDBSuite) TestData1() {

	dataInsert := map[string]interface{}{
		"Name":  "admin",
		"Email": "email@example.com",
	}
	ts.T().Logf("dataInsert: %#v", dataInsert)

	ts.Suite.Run("insert Test", func() {
		query := fmt.Sprintf(`INSERT INTO %s (Name, Email) VALUES (:Name, :Email)`, tableName)
		count, err := ts.db.NamedExec(query, dataInsert)
		ts.Require().NoError(err)
		ts.Equal(int64(1), count)
	})

	//===========================================================
	ts.Suite.Run("select Test", func() {
		query := fmt.Sprintf(`select name, email, created_at from %s`, tableName)
		var users []User
		err := ts.db.Select(&users, query)
		ts.Require().NoError(err)
		ts.Len(users, 1)
		ts.Equal("admin", users[0].Name)
		ts.T().Logf("%+v", users)

		// ===========================================================
		ts.T().Log("select Test where")
		query = fmt.Sprintf(`select name, email, created_at from %s where name=:Name`, tableName)
		var users2 []User
		err = ts.db.NamedSelect(&users2, query, map[string]interface{}{"Name": "admin"})
		ts.Require().NoError(err)
		ts.Len(users, 1)
		ts.Equal("email@example.com", users[0].Email)
		ts.T().Logf("%+v", users)
	})

	//===========================================================
	ts.Suite.Run("update Test", func() {
		query := fmt.Sprintf(`UPDATE %s SET email=:Email where name=:Name`, tableName)
		dataUpdate := map[string]interface{}{
			"Name":  "admin",
			"Email": "email_update@example.com",
		}
		ts.T().Logf("dataUpdate: %+v", dataUpdate)
		count, err := ts.db.NamedExec(query, dataUpdate)
		ts.Require().NoError(err)
		ts.Equal(int64(1), count)
	})

	//===========================================================
	ts.Suite.Run("delete Test", func() {
		query := fmt.Sprintf(`delete from %s where name=:Name`, tableName)
		dataDelete := map[string]interface{}{
			"Name": "admin",
		}
		ts.T().Logf("dataDelete: %+v", dataDelete)
		count, err := ts.db.NamedExec(query, dataDelete)
		ts.Require().NoError(err)
		ts.Equal(int64(1), count)
	})
}

func (ts *TestDBSuite) TestData2() {

	// batch insert with maps
	dtIns := []map[string]interface{}{
		{"Name": "admin", "Email": "email@example.com"},
		{"Name": "manager", "Email": "manager@gmail.com"},
		{"Name": "analyst", "Email": "analyst@mail.ru"},
	}
	ts.T().Logf("данные для вставки: %+v", dtIns)
	query := fmt.Sprintf(`INSERT INTO %s (Name, Email) VALUES (:Name, :Email)`, tableName)
	_, err := ts.db.NamedExec(query, dtIns)
	ts.NoError(err)

	//==================================================
	query = fmt.Sprintf(`select name, email from %s`, tableName)
	resultMap, err := ts.db.SelectMaps(query)
	ts.Require().NoError(err)
	ts.Len(resultMap, 3)
	ts.T().Logf("получим в map: %+v", resultMap)

	//==================================================
	var resultSlice []User
	err = ts.db.Select(&resultSlice, query)
	ts.Require().NoError(err)
	ts.Len(resultSlice, 3)
	ts.T().Logf("получим в slice: %+v", resultSlice)
}
