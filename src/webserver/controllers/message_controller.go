package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type IMessage interface {
	AllMessages(w http.ResponseWriter, r *http.Request)
	UserMessages(w http.ResponseWriter, r *http.Request)
	SetupRoutes(r *mux.Router)
}

type Message struct {
	store    sessions.Store
	messages services.IMessage
	users    services.IUser
}

func NewMessage(store sessions.Store, messages services.IMessage, users services.IUser) *Message {
	return &Message{store: store, messages: messages, users: users}
}

func (m *Message) AllMessages(w http.ResponseWriter, r *http.Request) {
	msgs, err := m.messages.ReadAllMessages()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonify, _ := json.Marshal(&msgs)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonify)
}

func (m *Message) UserMessages(w http.ResponseWriter, r *http.Request) {
	session, _ := m.store.Get(r, "session-name")
	if isAuthenticated, found := session.Values["isAuthenticated"].(bool); !isAuthenticated || !found {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	username := session.Values["username"].(string)

	msgs, err := m.messages.ReadAllMessagesForUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonify, _ := json.Marshal(msgs)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonify)
}

func (m *Message) SetupRoutes(r *mux.Router) {
	r.HandleFunc("/", m.UserMessages)
	r.HandleFunc("/public", m.AllMessages)
}
