package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/antonPalmFolkmann/DevOps2022/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
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
	log      *logrus.Logger
}

func NewUserController(users services.IUser, messages services.IMessage, store sessions.Store, log *logrus.Logger) *User {
	return &User{users: users, messages: messages, store: store, log: log}
}

func (u *User) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		u.log.Warnf("Failed to parse JSON body with error %s", err.Error())
		http.Error(w, "Could not parse the JSON body", http.StatusInternalServerError)
		return
	}

	var data RegisterReq
	err = json.Unmarshal(body, &data)
	if err != nil {
		u.log.Warnf("JSON body malformed with error %s", err.Error())
		http.Error(w, "The JSON body is malformed", http.StatusBadRequest)
		return
	}

	if u.users.IsUsernameTaken(data.Username) {
		u.log.Warnf("Username already taken, error")
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}

	if err := u.users.CreateUser(data.Username, data.Email, data.Password); err != nil {
		u.log.Warnf("Failed to create user with error %s", err.Error())
		http.Error(w, "Could not create a new user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Access-Control-Allow-Origin", "*")
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

	if !u.users.IsPasswordCorrect(data.Username, data.Password) {
		http.Error(w, "Password is incorrect", http.StatusForbidden)
		return
	}

	user, _ := u.users.ReadUserByUsername(data.Username)

	session.Values["username"] = user.Username
	session.Values["isAuthenticated"] = true
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "There was an error while saving session data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	resp := LoginResp{
		Username: user.Username,
		Email:    user.Email,
		Avatar:   "not yet implemented",
		Follows:  followersToUsernames(user.Follows),
	}
	jsonify, _ := json.Marshal(&resp)
	_, err = w.Write(jsonify)
	if err != nil {
		http.Error(w, "There was an error while writing the response", http.StatusInternalServerError)
	}

}

func followersToUsernames(followers []*storage.User) []string {
	usernames := make([]string, 0)
	for _, f := range followers {
		usernames = append(usernames, f.Username)
	}
	return usernames
}

func (u *User) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := u.store.Get(r, "session-name")
	if isAuthenticated, found := session.Values["isAuthenticated"].(bool); !isAuthenticated || !found {
		http.Error(w, "Must be logged in to log out", http.StatusBadRequest)
		return
	}

	session.Values["isAuthenticated"] = false
	delete(session.Values, "username")
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func (u *User) Timeline(w http.ResponseWriter, r *http.Request) {
	session, _ := u.store.Get(r, "session-name")
	if isAuthenticated, found := session.Values["isAuthenticated"].(bool); !isAuthenticated || !found {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	username := session.Values["username"].(string)

	msgs, err := u.messages.ReadAllMessagesOfFollowedUsers(username)
	if err != nil {
		http.Error(w, "There was an error while reading the messages", http.StatusInternalServerError)
		return
	}

	jsonify, _ := json.Marshal(msgs)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, err = w.Write(jsonify)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (u *User) Follow(w http.ResponseWriter, r *http.Request) {
	session, _ := u.store.Get(r, "session-name")
	if isAuthenticated, found := session.Values["isAuthenticated"].(bool); !isAuthenticated || !found {
		http.Error(w, "You must be logged in to follow", http.StatusForbidden)
		return
	}

	username := session.Values["username"].(string)

	whomname, err := utils.ParseUsername(r)
	if err != nil {
		http.Error(w, "There is no username to follow", http.StatusNotFound)
		return
	}

	if err := u.users.Follow(username, whomname); errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "That user does not exist", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "There was an error while completing the follow operation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func (u *User) Unfollow(w http.ResponseWriter, r *http.Request) {
	session, _ := u.store.Get(r, "session-name")
	if isAuthenticated, found := session.Values["isAuthenticated"].(bool); !isAuthenticated || !found {
		http.Error(w, "You must be logged in to unfollow", http.StatusForbidden)
		return
	}

	username := session.Values["username"].(string)

	whomname, err := utils.ParseUsername(r)
	if err != nil {
		http.Error(w, "There is no username to unfollow", http.StatusNotFound)
		return
	}

	if err := u.users.Unfollow(username, whomname); errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "That user does not exist", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "There was an error while completing the unfollow operation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func (u *User) SetupRoutes(r *mux.Router) {
	r.HandleFunc("/api/register", u.Register)
	r.HandleFunc("/api/login", u.Login)
	r.HandleFunc("/api/logout", u.Logout)
	r.HandleFunc("/api/", u.Timeline)
	r.HandleFunc("/api/fllw/{username}", u.Follow)
	r.HandleFunc("/api/unfllw/{username}", u.Unfollow)
}
