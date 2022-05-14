package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/antonPalmFolkmann/DevOps2022/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
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
	log      *logrus.Logger
}

func NewMessage(store sessions.Store, messages services.IMessage, users services.IUser, log *logrus.Logger) *Message {
	return &Message{store: store, messages: messages, users: users, log: log}
}

func (m *Message) AllMessages(w http.ResponseWriter, r *http.Request) {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}
	m.log.Debugf("Read offset to be %d", offset)

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 100
	}
	m.log.Debugf("Read limit to be %d", limit)

	msgs, err := m.messages.ReadAllMessages(limit, offset)
	if err != nil {
		http.Error(w, "There was an error while reading messages", http.StatusInternalServerError)
		return
	}

	jsonify, _ := json.Marshal(&msgs)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, err = w.Write(jsonify)
	if err != nil {
		http.Error(w, "An error occured during writing all messages", http.StatusInternalServerError)
		return
	}
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	resp := MsgsPerUsernameResp{
		Username: user.Username,
		Email:    user.Email,
		Avatar:   "not yet implemented...",
		Msgs:     msgs,
	}
	jsonify, _ := json.Marshal(resp)
	_, err = w.Write(jsonify)
	if err != nil {
		http.Error(w, "An error occured during writing all messages by username", http.StatusInternalServerError)
		return
	}
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

	err = m.messages.CreateMessage(data.AuthorName, data.Text)
	if err != nil {
		http.Error(w, "Failed to add message", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func (m *Message) SetupRoutes(r *mux.Router) {
	r.HandleFunc("/api/msgs/{username}", m.UserMessages)
	r.HandleFunc("/api/public", m.AllMessages)
	r.HandleFunc("/api/add_message", m.AddMessage)
}
