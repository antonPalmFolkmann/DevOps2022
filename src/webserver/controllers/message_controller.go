package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/antonPalmFolkmann/DevOps2022/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
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
	vars := mux.Vars(r)
	if _, found := vars["username"]; !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	username := vars["username"]

	id, err := m.users.ReadUserIdByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "There are no messages for that user because they don't exist", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "There was an error while trying to read the user id", http.StatusInternalServerError)
		return
	}

	msgs, err := m.messages.ReadAllMessagesByAuthorId(id)
	if err != nil {
		http.Error(w, "There was an error while trying to read the messages", http.StatusInternalServerError)
		return
	}

	filteredMsgs := make([]MsgResp, 0)
	for _, msg := range msgs {
		author, _ := m.users.ReadUserById(msg.UserID)
		filteredM := MsgResp{
			AuthorName: author.Username,
			Text:       msg.Text,
			PubDate:    utils.FormatDatetime(int64(msg.PubDate)),
			Flagged:    msg.Flagged,
		}
		filteredMsgs = append(filteredMsgs, filteredM)
	}

	user, err := m.users.ReadUserByUsername(username)
	if err != nil {
		http.Error(w, "There was an error while trying to read the user information", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	resp := MsgsPerUsernameResp{
		Username: user.Username,
		Email:    user.Email,
		Avatar:   "not yet implemented...",
		Msgs:     filteredMsgs,
	}
	jsonify, _ := json.Marshal(resp)
	w.Write(jsonify)
}

func (m *Message) SetupRoutes(r *mux.Router) {
	r.HandleFunc("/msgs/{username}", m.UserMessages)
	r.HandleFunc("/public", m.AllMessages)
}
