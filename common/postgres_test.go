// Copyright © 2023 OSINTAMI. This is not yours.
package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresOpen(t *testing.T) {
	cfg := &PostgresConfig{
		PgHost:     "",
		PgPort:     "",
		PgUser:     "",
		PgPassword: "",
		PgDB:       ""}

	db, err := OpenDB(cfg, ".")
	assert.NotNil(t, db)
	assert.NotNil(t, err)
}

func TestPostgresOpenBadPath(t *testing.T) {
	cfg := &PostgresConfig{
		PgHost:     "",
		PgPort:     "",
		PgUser:     "",
		PgPassword: "",
		PgDB:       ""}

	db, err := OpenDB(cfg, "...")
	assert.Nil(t, db)
	assert.NotNil(t, err)
}

func TestPostgresMock(t *testing.T) {
	db, mock := CreateMockDatabase()
	assert.NotNil(t, db)
	assert.NotNil(t, mock)
}
