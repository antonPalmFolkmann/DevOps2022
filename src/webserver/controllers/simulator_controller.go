package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/antonPalmFolkmann/DevOps2022/storage"
)

type IService interface {
	IsAuthorized(token string) bool
	ReadLatest() int
	UpdateLatest(latest int)
}

type Simulator struct {
	messageService   storage.IMessage
	userSservice     storage.IUser
	simulatorService IService
	followerService  storage.IFollows
}

func NewSimulator(messageService storage.IMessage, userService storage.IUser, simulatorService IService, followerService storage.IFollows) *Simulator {
	return &Simulator{messageService: messageService, userSservice: userService, simulatorService: simulatorService, followerService: followerService}
}

func (s *Simulator) LatestHandler(w http.ResponseWriter, r *http.Request) {
	latest := s.simulatorService.ReadLatest()
	respMsg := fmt.Sprintf("{\"latest\": %d}", latest)

	jsonData := []byte(respMsg)
	w.Write(jsonData)
}

func (s *Simulator) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := s.updateLatest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !isAuthorized(w, r) {
		return
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	var requestBody registerRequestBody
	err = json.Unmarshal(body, &requestBody)
	//Error handling if the struct doesn't get the necessary paramters for initialization
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	var regError string
	if r.Method == http.MethodPost {
		//TODO Implement service -> query db and store the user registering
	}

	if regError != "" {
		w.WriteHeader(400)
		jsonify := fmt.Sprintf("\"status\": %d, \"error_msg\": %s", 400, regError)
		w.Write([]byte(jsonify))
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Simulator) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	err := s.updateLatest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func (s *Simulator) updateLatest(r *http.Request) error {
	if !r.URL.Query().Has("latest") {
		return nil
	}

	latest, err := parseLatest(r)
	if err != nil {
		return errors.New("Latest was not an integer")
	}

	s.simulatorService.UpdateLatest(*latest)
	return nil
}

func parseLatest(r *http.Request) (*int, error) {
	asInt, err := strconv.Atoi(r.URL.Query().Get("latest"))
	if err != nil {
		return nil, errors.New("latest is not an integer")
	}

	return &asInt, nil
}

func isAuthorized(w http.ResponseWriter, r *http.Request) bool {
	authorizedReq := r.Header.Get("Authorization")
	if authorizedReq != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
		error := "You are not authorized to use this resource!"
		jsonify := fmt.Sprintf("\"status\": %d, \"error_msg\": %s", 403, error)
		w.Write([]byte(jsonify))
		return false
	}
	return true
}

type registerRequestBody struct {
	Latest   int    `json:"latest"`
	PostType string `json:"post_type"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"user_pwd"`
}

type tweetRequestBody struct {
	Latest   int    `json:"latest"`
	PostType string `json:"post_type"`
	Username string `json:"username"`
	Content  string `json:"content"`
}

type followRequestBody struct {
	Latest   int    `json:"latest"`
	PostType string `json:"post_type"`
	Username string `json:"username"`
	Follow   string `json:"user_to_follow"`
}

type unfollowRequestBody struct {
	Latest   int    `json:"latest"`
	PostType string `json:"post_type"`
	Username string `json:"username"`
	Unfollow string `json:"user_to_unfollow"`
}