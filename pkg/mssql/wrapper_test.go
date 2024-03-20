package mssql

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/stretchr/testify/suite"
)

// Person человек
type Person struct {
	LastName    string     `db:"last_name"`
	Birthdate   *time.Time `db:"birthdate"`
	Salary      *float64
	IsOwnerFlat *bool `db:"is_owner_flat"` // признак владельца помещения
	Email       string
	CreatedAt   time.Time `db:"created_at"`
}

// DBSuite структура для набора тестов с БД
type TestDBSuite struct {
	suite.Suite
	db *MSSQL // коннект к БД master
}

var (
	dbName    = "go_db_test"
	tableName = fmt.Sprintf("%s.dbo.people", dbName)
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

	_, err := its.db.DBX.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS %s; CREATE DATABASE %s`, dbName, dbName))
	if err != nil {
		its.FailNowf("unable to create database", err.Error())
	}
	its.T().Logf("База [%s] создана\n", dbName)

	query := fmt.Sprintf(`CREATE TABLE %s (
		last_name varchar(50) PRIMARY KEY,
		birthdate datetime,
		salary decimal(15,2),
		is_owner_flat bit,
		email varchar(100) UNIQUE,
		created_at datetime NOT NULL DEFAULT current_timestamp
	)`, tableName)

	_, err = its.db.DBX.Exec(query)
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
		"LastName": "Иванов",
		"Email":    "ivan@example.com",
		//"is_owner_flat": true,
	}
	ts.T().Logf("dataInsert: %#v", dataInsert)

	ts.Suite.Run("insert Test", func() {
		query := fmt.Sprintf(`INSERT INTO %s (last_name, Email) VALUES (:LastName, :Email)`, tableName)
		count, err := ts.db.NamedExec(query, dataInsert)
		ts.Require().NoError(err)
		ts.Equal(int64(1), count)
	})

	//===========================================================
	ts.Suite.Run("select Test", func() {
		query := fmt.Sprintf(`select * from %s`, tableName)
		var people []Person
		err := ts.db.Select(&people, query)
		ts.Require().NoError(err)
		ts.Len(people, 1)
		ts.Equal("Иванов", people[0].LastName)
		ts.T().Logf("%+v", people)

		// ===========================================================
		ts.T().Log("select Test where")
		query = fmt.Sprintf(`select last_name, email, created_at from %s where last_name=:Name`, tableName)
		var people2 []Person
		err = ts.db.NamedSelect(&people2, query, map[string]interface{}{"Name": "Иванов"})
		ts.Require().NoError(err)
		ts.Len(people2, 1)
		ts.Equal("ivan@example.com", people2[0].Email)
		ts.T().Logf("%+v", people2)
	})

	//===========================================================
	ts.Suite.Run("update Test", func() {
		query := fmt.Sprintf(`UPDATE %s SET email=:Email where last_name=:Name`, tableName)
		dataUpdate := map[string]interface{}{
			"Name":  "Иванов",
			"Email": "email_update@example.com",
		}
		ts.T().Logf("dataUpdate: %+v", dataUpdate)
		count, err := ts.db.NamedExec(query, dataUpdate)
		ts.Require().NoError(err)
		ts.Equal(int64(1), count)
	})

	//===========================================================
	ts.Suite.Run("delete Test", func() {
		query := fmt.Sprintf(`delete from %s where last_name=:Name`, tableName)
		dataDelete := map[string]interface{}{
			"Name": "Иванов",
		}
		ts.T().Logf("dataDelete: %+v", dataDelete)
		count, err := ts.db.NamedExec(query, dataDelete)
		ts.Require().NoError(err)
		ts.Equal(int64(1), count, "удалено")
	})
}

func (ts *TestDBSuite) TestData2() {

	// batch insert with maps
	dtIns := []map[string]interface{}{
		{"LastName": "Сидоров", "Email": "sidr@example.com", "Birthdate": time.Date(2000, 2, 21, 0, 0, 0, 0, time.UTC)},
		{"LastName": "Кузнецов", "Email": "kuz@gmail.com", "Birthdate": nil},
		{"LastName": "Петров", "Email": "petr@mail.ru", "Birthdate": nil},
	}
	ts.T().Logf("данные для вставки: %+v", dtIns)
	query := fmt.Sprintf(`INSERT INTO %s (last_name, Email, Birthdate) VALUES (:LastName, :Email, :Birthdate)`, tableName)
	_, err := ts.db.NamedExec(query, dtIns)
	ts.NoError(err)

	//==================================================
	query = fmt.Sprintf(`select * from %s`, tableName)
	resultMap, err := ts.db.SelectMaps(query)
	ts.Require().NoError(err)
	ts.Len(resultMap, 3)
	ts.T().Logf("получим в map: %+v", resultMap)

	//==================================================
	var resultSlice []Person
	err = ts.db.Select(&resultSlice, query)
	ts.Require().NoError(err)
	ts.Len(resultSlice, 3)
	ts.T().Logf("получим в slice: %+v", resultSlice)
}
