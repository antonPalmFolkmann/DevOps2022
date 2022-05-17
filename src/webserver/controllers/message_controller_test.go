package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/antonPalmFolkmann/DevOps2022/controllers"
	"github.com/stretchr/testify/assert"
)

func TestAddMessageReturnsStatusCreated(t *testing.T) {
	r := setUp()
	session := login(r)

	msg := &controllers.AddMsgsReq{
		AuthorName: "rnsk",
		Text:       "msg text",
	}

	jsonMsg, _ := json.Marshal(msg)
	req, _ := http.NewRequest("POST", "/api/add_message", bytes.NewBuffer(jsonMsg))
	req.AddCookie(session)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Result().StatusCode)
}

func TestAddMessageReturnsStatusForbidden(t *testing.T) {
	r := setUp()

	msg := &controllers.AddMsgsReq{
		AuthorName: "rnsk",
		Text:       "msg text",
	}

	jsonMsg, _ := json.Marshal(msg)
	req, _ := http.NewRequest("POST", "/api/add_message", bytes.NewBuffer(jsonMsg))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusForbidden, resp.Result().StatusCode)
}

func TestAllMessagesReturnsStatusOk(t *testing.T) {
	r := setUp()

	req, _ := http.NewRequest("GET", "/api/public", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
}

func TestUserMessageGivenExistingUserReturnsStatusOK(t *testing.T) {
	r := setUp()

	req, _ := http.NewRequest("GET", "/api/msgs/rnsk", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Result().StatusCode)
}

func TestUserMessageGivenNotExistingUserReturnsStatusOK(t *testing.T) {
	r := setUp()

	req, _ := http.NewRequest("GET", "/api/msgs/ibaby", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Result().StatusCode)
}
