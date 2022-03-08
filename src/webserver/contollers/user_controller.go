package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/gorilla/mux"
)

type IUserController interface {
	MessagesPerUserHandler(w http.ResponseWriter, r *http.Request)
	FollowHandler(w http.ResponseWriter, r *http.Request)
	UnfollowHandler(w http.ResponseWriter, r *http.Request)
	SetupRoutes(r *mux.Router)
}

type UserController struct {
	userService     services.IUserService
	messagesService services.IMessageService
}

func NewUserController(userService services.IUserService, messagesService services.IMessageService) *UserController {
	return &UserController{userService: userService}
}

func (u *UserController) SetupRoutes(r *mux.Router) {
	r.HandleFunc("/msgs/{username}", u.MessagesPerUserHandler)
}

func (u *UserController) MessagesPerUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		u.getMessagesByUser(w, r)
	} else if r.Method == http.MethodPost {
		u.postMessage(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (u *UserController) getMessagesByUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	limit, err := parseNo(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	user, err := u.userService.ReadUserIdByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	msgs, err := u.messagesService.ReadAllMessagesByAuthorId(user, *limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(&msgs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (u *UserController) postMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var payload PostMessageRequest
	err = json.Unmarshal(body, &payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// TODO: Maybe we do need DTOs/method headers that mimic DTOs

	log.Println(username)
}

func parseNo(r *http.Request) (*int, error) {
	asInt, err := strconv.Atoi(r.URL.Query().Get("no"))
	if err != nil {
		return nil, errors.New("no query parameter is not a number")
	}
	return &asInt, nil
}
