package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type IUser interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Timeline(w http.ResponseWriter, r *http.Request)
	Follow(w http.ResponseWriter, r *http.Request)
	Unfollow(w http.ResponseWriter, r *http.Request)
}

type UserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterReq struct {
	UserReq
	Email string `json:"email"`
}

type User struct {
	store    sessions.Store
	users    services.IUser
	messages services.IMessage
}

func NewUserController(users services.IUser, messages services.IMessage, store sessions.Store) *User {
	return &User{users: users, messages: messages, store: store}
}

func (u *User) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var data RegisterReq
	err = json.Unmarshal(body, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pwHash := u.users.Hash(data.Password)
	err = u.users.CreateUser(data.Username, data.Email, pwHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var data UserReq
	err = json.Unmarshal(body, &data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := u.users.ReadUserByUsername(data.Username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if u.users.Hash(data.Password) != user.PwHash {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	session, _ := u.store.Get(r, "session-name")
	session.Values["username"] = data.Username
	session.Values["isAuthenticated"] = true
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (u *User) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := u.store.Get(r, "session-name")
	if isAuthenticated, found := session.Values["isAuthenticated"].(bool); !isAuthenticated || !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session.Values["isAuthenticated"] = false
	session.Save(r, w)

	w.WriteHeader(http.StatusOK)
}

func (u *User) Timeline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if _, found := vars["username"]; !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := vars["username"]

	id, err := u.users.ReadUserIdByUsername(user)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	msgs, err := u.messages.ReadAllMessagesByAuthorId(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonify, err := json.Marshal(&msgs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonify)
	w.WriteHeader(http.StatusNotImplemented)
}

func (u *User) Follow(w http.ResponseWriter, r *http.Request) {
	session, _ := u.store.Get(r, "session-name")
	if isAuthenticated, found := session.Values["isAuthenticated"].(bool); !isAuthenticated || !found {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	username, err := parseUsername(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	toFollowID, err := u.users.ReadUserIdByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userID, _ := u.users.ReadUserIdByUsername(session.Values["username"].(string))

	err = u.users.Follow(userID, toFollowID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (u *User) Unfollow(w http.ResponseWriter, r *http.Request) {
	session, _ := u.store.Get(r, "session-name")
	if isAuthenticated, found := session.Values["isAuthenticated"].(bool); !isAuthenticated || !found {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	username, err := parseUsername(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	toUnfollowID, err := u.users.ReadUserIdByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userID, _ := u.users.ReadUserIdByUsername(session.Values["username"].(string))

	err = u.users.Unfollow(userID, toUnfollowID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func parseUsername(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	if username, found := vars["username"]; !found {
		return "", errors.New("there is no username")
	} else {
		return username, nil
	}
}
