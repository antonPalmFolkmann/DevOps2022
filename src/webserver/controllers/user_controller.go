package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
)

type IUser interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Timeline(w http.ResponseWriter, r *http.Request)
	Follow(w http.ResponseWriter, r *http.Request)
	Unfollow(w http.ResponseWriter, r *http.Request)
	SetupRoutes(r *mux.Router)
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
		http.Error(w, "Could not parse the JSON body", http.StatusInternalServerError)
		return
	}

	var data RegisterReq
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "The JSON body is malformed", http.StatusBadRequest)
		return
	}

	if _, err := u.users.ReadUserByUsername(data.Username); !errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}

	pwHash := u.users.Hash(data.Password)
	err = u.users.CreateUser(data.Username, data.Email, pwHash)
	if err != nil {
		http.Error(w, "Could not create a new user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp := RegisterResp{Error: ""}
	jsonify, _ := json.Marshal(&resp)
	w.Write(jsonify)
}

func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	session, _ := u.store.Get(r, "session-name")

	if isAuthenticated, _ := session.Values["isAuthenticated"].(bool); isAuthenticated {
		http.Error(w, "Already logged in", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not parse the JSON body", http.StatusInternalServerError)
		return
	}

	var data UserReq
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "The JSON body is malformed", http.StatusBadRequest)
		return
	}

	user, err := u.users.ReadUserByUsername(data.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Cannot login to a user that does not exist", http.StatusNotFound)
		return
	}

	if u.users.Hash(data.Password) != user.PwHash {
		http.Error(w, "Password is incorrect", http.StatusForbidden)
		return
	}

	session.Values["username"] = user.Username
	session.Values["isAuthenticated"] = true
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "There was an error while saving your request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := LoginResp{
		Username: user.Username,
		Email:    user.Email,
		Avatar:   "not yet implemented",
		Follows:  []string{}, // FIXME: Need this as part of the user services interface
	}
	jsonify, _ := json.Marshal(&resp)
	w.Write(jsonify)

}

func (u *User) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := u.store.Get(r, "session-name")
	if isAuthenticated, found := session.Values["isAuthenticated"].(bool); !isAuthenticated || !found {
		http.Error(w, "Must be logged in to log out", http.StatusBadRequest)
		return
	}

	session.Values["isAuthenticated"] = false
	session.Save(r, w)

	w.WriteHeader(http.StatusOK)
}

func (u *User) Timeline(w http.ResponseWriter, r *http.Request) {
	session, _ := u.store.Get(r, "session-name")
	if isAuthenticated, found := session.Values["isAuthenticated"].(bool); !isAuthenticated || !found {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	username := session.Values["username"].(string)

	msgs, err := u.messages.ReadAllMessagesForUsername(username)
	if err != nil {
		http.Error(w, "There was an error while reading the messages", http.StatusInternalServerError)
		return
	}

	jsonify, _ := json.Marshal(msgs)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonify)
}

func (u *User) Follow(w http.ResponseWriter, r *http.Request) {
	session, _ := u.store.Get(r, "session-name")
	if isAuthenticated, found := session.Values["isAuthenticated"].(bool); !isAuthenticated || !found {
		http.Error(w, "You must be logged in to follow", http.StatusForbidden)
		return
	}

	username, err := parseUsername(r)
	if err != nil {
		http.Error(w, "There is no username to follow", http.StatusNotFound)
		return
	}

	toFollowID, err := u.users.ReadUserIdByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "The user being followed does not exist", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "There was an error while reading user information", http.StatusInternalServerError)
		return
	}

	userID, _ := u.users.ReadUserIdByUsername(session.Values["username"].(string))

	err = u.users.Follow(userID, toFollowID)
	if err != nil {
		http.Error(w, "There was an error while performing the follow operation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (u *User) Unfollow(w http.ResponseWriter, r *http.Request) {
	session, _ := u.store.Get(r, "session-name")
	if isAuthenticated, found := session.Values["isAuthenticated"].(bool); !isAuthenticated || !found {
		http.Error(w, "You must be logged in to unfollow", http.StatusForbidden)
		return
	}

	username, err := parseUsername(r)
	if err != nil {
		http.Error(w, "There is no username to follow", http.StatusNotFound)
		return
	}

	toUnfollowID, err := u.users.ReadUserIdByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "The user being followed does not exist", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "There was an error while reading user information", http.StatusInternalServerError)
		return
	}

	userID, _ := u.users.ReadUserIdByUsername(session.Values["username"].(string))

	err = u.users.Unfollow(userID, toUnfollowID)
	if err != nil {
		http.Error(w, "There was an error while performing the unfollow operation", http.StatusInternalServerError)
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

func (u *User) SetupRoutes(r *mux.Router) {
	r.HandleFunc("/register", u.Register)
	r.HandleFunc("/login", u.Login)
	r.HandleFunc("/logout", u.Logout)
	r.HandleFunc("/msgs/{username}", u.Timeline)
	r.HandleFunc("/fllw/{username}", u.Follow)
	r.HandleFunc("/unfllw/{username}", u.Unfollow)
}
