package mssql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getDatabaseUrl(t *testing.T) {
	var tests = []struct {
		in   Config
		want string
	}{
		{Config{}, "sqlserver://:@:0?database="},
		{Config{Host: "host", Port: 1433, Database: "database_name", User: "user", Password: "password"},
			"sqlserver://user:password@host:1433?database=database_name"},
		{Config{Host: "host", Port: 1433, Database: "database_name", User: "user", Password: "password", APPName: "APP"},
			"sqlserver://user:password@host:1433?app+name=APP&database=database_name"},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, test.in.getDatabaseURL())
	}
}
