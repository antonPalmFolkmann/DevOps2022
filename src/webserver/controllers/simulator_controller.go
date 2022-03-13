package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	services "github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/gorilla/mux"
)

type ISimulatorService interface {
	IsAuthorized(token string) bool
	ReadLatest() int
	UpdateLatest(latest int)
}

type Simulator struct {
	messageService   services.IMessage
	userService      services.IUser
	simulatorService ISimulatorService
}

func NewSimulator(messageService services.IMessage, userService services.IUser, simulatorService ISimulatorService) *Simulator {
	return &Simulator{messageService: messageService, userService: userService, simulatorService: simulatorService}
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

	if !IsAuthorized(w, r) {
		return
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	var requestBody registerRequestBody
	err = json.Unmarshal(body, &requestBody)
	//Error handling if the struct doesn't get the necessary paramters for initialization
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonify := fmt.Sprintf("\"status\": %d, \"error_msg\": %s", 400, err.Error())
		w.Write([]byte(jsonify))
		return
	}

	var regError string
	if r.Method == http.MethodPost {
		//TODO Implement service -> query db and store the user registering
		if user, _ := s.userService.ReadUserByUsername(requestBody.Username); user.Username != "" {
			regError = "The username is already taken"
			w.WriteHeader(http.StatusBadRequest)
			jsonify := fmt.Sprintf("\"status\": %d, \"error_msg\": %s", 400, regError)
			w.Write([]byte(jsonify))
			return
		}
		err := s.userService.CreateUser(requestBody.Username, requestBody.Email, s.userService.Hash(requestBody.Password))
		if err != nil {
			// w.WriteHeader(http.StatusInternalServerError)
			log.Println("simulator_controller: An error occured during creation of a user")
		}
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(""))
	}
}

func (s *Simulator) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	err := s.updateLatest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !IsAuthorized(w, r) {
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filtered_msgs, err := s.messageService.ReadAllMessages()
	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		log.Println("simulator_controller: An error occured during reading all messages")
	}
	msgs, err := json.Marshal(filtered_msgs)
	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		log.Println("simulator_controller: An error occured during marshalling of messages")
	}
	w.WriteHeader(http.StatusOK)
	w.Write(msgs)
}

func (s *Simulator) UserPerMessageHandler(w http.ResponseWriter, r *http.Request) {
	err := s.updateLatest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !IsAuthorized(w, r) {
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

	if r.Method == http.MethodGet {
		s.postUserPerMessage(w, r, username)
	} else if r.Method == http.MethodPost {
		s.getUserPerMessage(w, r, username)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (s *Simulator) postUserPerMessage(w http.ResponseWriter, r *http.Request, username string) {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	var requestBody tweetRequestBody
	err := json.Unmarshal(body, &requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userid, err := s.userService.ReadUserIdByUsername(requestBody.Username)
	if err != nil {
		// w.WriteHeader(http.StatusNotFound)
		// w.Write([]byte(""))
		// return
		log.Println("service_simulator: user doesn't exist -> trying to post message")
	}
	err = s.messageService.CreateMessage(userid, requestBody.Content, false)

	if err != nil {
		// w.WriteHeader(http.StatusBadRequest)
		// return
		log.Println("service_simulator: failed to create messages in db")
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}

func (s *Simulator) getUserPerMessage(w http.ResponseWriter, r *http.Request, username string) {
	userid, err := s.userService.ReadUserIdByUsername(username)
	if err != nil {
		// w.WriteHeader(http.StatusNotFound)
		// w.Write([]byte(""))
		// return
		log.Println("service_simulator: user doesn't exist -> trying to post message")
	}
	filtered_msgs, err := s.messageService.ReadAllMessagesByAuthorId(userid)
	if err != nil {
		log.Println("service_simulator: Failed to read messages by author -> Author doesn't exist")
	}
	msgs, err := json.Marshal(filtered_msgs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(msgs)
	}
}

func (s *Simulator) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	err := s.updateLatest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !IsAuthorized(w, r) {
		return
	}

	if r.Method == http.MethodPost {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		var requestBody followRequestBody
		err := json.Unmarshal(body, &requestBody)
		if err != nil {
			s.unfollowUser(w, r)
		} else {
			s.followUser(w, r, requestBody)
		}
	}

	if r.Method == http.MethodGet {
		s.getFollowers(w, r)
	}
}

func (s *Simulator) followUser(w http.ResponseWriter, r *http.Request, body followRequestBody) {
	vars := mux.Vars(r)
	username := vars["username"]

	userid, err := s.userService.ReadUserIdByUsername(username)
	if err != nil {
		log.Println("service_simulator: user doesn't exist -> trying to follow user")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	useridToFollow, err := s.userService.ReadUserIdByUsername(body.Follow)
	if err != nil {
		log.Println("service_simulator: user doesn't exist -> trying to follow user")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = s.userService.Follow(userid, useridToFollow)
	if err != nil {
		log.Println("service_simulator: Failed to follow user")
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}

func (s *Simulator) unfollowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	var requestBody unfollowRequestBody
	err := json.Unmarshal(body, &requestBody)
	if err != nil {
		log.Println("Service_simulator: Request body is invalid")
	}

	userid, err := s.userService.ReadUserIdByUsername(username)
	if err != nil {
		log.Println("service_simulator: User doesn't exist -> trying to unfollow user")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	useridToUnfollow, err := s.userService.ReadUserIdByUsername(requestBody.Unfollow)
	if err != nil {
		log.Println("service_simulator: User doesn't exist -> trying to unfollow user")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = s.userService.Unfollow(userid, useridToUnfollow)
	if err != nil {
		log.Println("service_simulator: Failed to unfollow user")
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}

func (s *Simulator) getFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	filtered_followers, err := s.userService.ReadFollowersForUsername(username)
	if err != nil {
		log.Println("service_simulator: Failed to read all followers from username")
	}
	followers, err := json.Marshal(filtered_followers)
	if err != nil {
		log.Println("service_simulator: Failed to marshall followers")
	}
	w.Write(followers)
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

func IsAuthorized(w http.ResponseWriter, r *http.Request) bool {
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
