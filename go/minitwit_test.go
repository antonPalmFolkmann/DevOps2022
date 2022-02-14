package main

import (
	"bytes"
	"database/sql"
	"log"
	"net/http"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setUp() {
	db, _ = sql.Open("sqlite3", ":memory:")
	// Use httptest package instead of minitwit.app.test_client()
	DATABASE = ":memory:"
	InitDb()
}

// Helper functions Login, Logout, and RegisterAndLogin

func login(username string, password string) {
	var jsonData = []byte(`{
		"username": username,
		"password": password
	}`)
	_, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("fatal")
	}
}

func logout() {
	_, err := http.NewRequest("POST", "/logout", bytes.NewBuffer([]byte{}))
	if err != nil {
		log.Fatal("fatal")
	}
}

func register(username string, password string, password2 string, email string) {
	if password2 == "" {
		password2 = password
	}
	if email == "" {
		email = username + "@example.com"
	}

	var jsonData = []byte(`{
		"username": username,
		"password": password,
		"password2": password2,
		"email": email
	}`)
	_, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("fatal")
	}
}

func addMessage(t *testing.T, text string) {
	rv, err := http.NewRequest("POST", "/add_message", bytes.NewBuffer([]byte(text)))
	if err != nil {
		log.Fatal("fatal")
	}
	assert.Equal(t, "Your message was recorded", rv)
}

func RegisterAndLogin(username string, password string) {
	register(username, password, "", "")
	login(username, password)
}

func TestSomething(t *testing.T) {
	setUp()

	assert.True(t, true, "True is true!")
}
