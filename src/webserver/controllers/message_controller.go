package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"log/syslog"
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
	msgs, err := m.messages.ReadAllMessages(0, 100)
	if err != nil {
		http.Error(w, "There was an error while reading messages", http.StatusInternalServerError)
		return
	}
	logger, err := syslog.New(syslog.LOG_INFO, "MESSAGE CONTROLLER: ")
	if err != nil {
		log.Println(err)
	} else {
		logger.Info("READ ALL MESSAGES")
	}
	jsonify, _ := json.Marshal(&msgs)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonify)
}

func (m *Message) UserMessages(w http.ResponseWriter, r *http.Request) {
	username, err := utils.ParseUsername(r)
	if err != nil {
		http.Error(w, "There is no username", http.StatusNotFound)
		return
	}

	user, err := m.users.ReadUserByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "There was an error while reading user information", http.StatusInternalServerError)
		return
	}

	msgs, err := m.messages.ReadAllMessagesByUsername(username)
	if err != nil {
		http.Error(w, "There was an error while trying to read the messages", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := MsgsPerUsernameResp{
		Username: user.Username,
		Email:    user.Email,
		Avatar:   "not yet implemented...",
		Msgs:     msgs,
	}
	jsonify, _ := json.Marshal(resp)
	w.Write(jsonify)
}

func (m *Message) AddMessage(w http.ResponseWriter, r *http.Request) {
	session, _ := m.store.Get(r, "session-name")

	if isAuthenticated, found := session.Values["isAuthenticated"].(bool); !isAuthenticated || !found {
		http.Error(w, "You must be logged in to add a message", http.StatusForbidden)
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not parse the JSON body", http.StatusInternalServerError)
		return
	}

	var data AddMsgsReq
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "The JSON body is malformed", http.StatusBadRequest)
		return
	}

	m.messages.CreateMessage(data.AuthorName, data.Text)
	w.WriteHeader(http.StatusCreated)
}

func (m *Message) SetupRoutes(r *mux.Router) {
	r.HandleFunc("/msgs/{username}", m.UserMessages)
	r.HandleFunc("/public", m.AllMessages)
	r.HandleFunc("/add_message", m.AddMessage)
}
