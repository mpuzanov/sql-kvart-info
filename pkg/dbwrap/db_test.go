package dbwrap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	config := NewConfig()

	_, err := New(config)
	assert.ErrorIs(t, err, ErrBadConfigDB)

	config.WithPassword("12345")
	_, err = New(config)
	assert.ErrorContains(t, err, "sqlx.Connect")

	config.WithPassword("123")
	db, err := New(config)
	assert.NoError(t, err)

	db.SetTimeout(100)
	assert.Equal(t, 100, db.Cfg.TimeoutQuery)
}
