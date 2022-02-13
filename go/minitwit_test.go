package main

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func SetUp() {
	db, _ = sql.Open("sqlite3", ":memory:")
	// Use httptest package instead of minitwit.app.test_client()
	DATABASE = ":memory:"
	InitDb()
}

func TestSomething(t *testing.T) {
	SetUp()

	assert.True(t, true, "True is true!")
}
