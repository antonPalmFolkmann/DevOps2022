package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
	assert.True(t, true, "True is true!")
}

func Test_gravatar_url_given_email_and_size_returns_url(t *testing.T) {
	actual := gravatar_url("email", 80)
	expected := "http://www.gravatar.com/avatar/656d61696c?d=identicon&s=80"
	
	assert.Equal(t, actual, expected, "Should be the same hash")
}
