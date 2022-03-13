package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var data RegisterReq
	err = json.Unmarshal(body, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := u.users.ReadUserByUsername(data.Username); !errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusConflict)
		resp := RegisterResp{Error: "Username already taken"}
		jsonify, _ := json.Marshal(&resp)
		w.Write(jsonify)
		return
	}

	pwHash := u.users.Hash(data.Password)
	err = u.users.CreateUser(data.Username, data.Email, pwHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var data UserReq
	err = json.Unmarshal(body, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

	session.Values["username"] = user.Username
	session.Values["isAuthenticated"] = true
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	username := vars["username"]

	id, err := u.users.ReadUserIdByUsername(username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msgs, err := u.messages.ReadAllMessagesByAuthorId(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	filteredMsgs := make([]MsgResp, 0)
	for _, m := range msgs {
		author, _ := u.users.ReadUserById(m.UserID)
		filteredM := MsgResp{
			AuthorName: author.Username,
			Text:       m.Text,
			PubDate:    formatDatetime(int64(m.PubDate)),
			Flagged:    m.Flagged,
		}
		filteredMsgs = append(filteredMsgs, filteredM)
	}

	user, err := u.users.ReadUserByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

func formatDatetime(timestamp int64) string {
	timeUnix := time.Unix(timestamp, 0)
	return timeUnix.Format("2006-01-02 15:04")
}

func (u *User) Follow(w http.ResponseWriter, r *http.Request) {
	session, _ := u.store.Get(r, "session-name")
	if isAuthenticated, found := session.Values["isAuthenticated"].(bool); !isAuthenticated || !found {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	username, err := parseUsername(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusBadRequest)
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
		w.WriteHeader(http.StatusInternalServerError)
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
