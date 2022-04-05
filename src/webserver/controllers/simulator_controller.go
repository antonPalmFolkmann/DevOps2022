package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	services "github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Simulator struct {
	messageService   services.IMessage
	userService      services.IUser
	simulatorService services.ISimulatorService
	log              *logrus.Logger
}

func NewSimulator(messageService services.IMessage, userService services.IUser, simulatorService services.ISimulatorService, log *logrus.Logger) *Simulator {
	return &Simulator{messageService: messageService, userService: userService, simulatorService: simulatorService, log: log}
}

func (s *Simulator) LatestHandler(w http.ResponseWriter, r *http.Request) {
	s.log.Trace("Hit latest endpoint")

	latest := s.simulatorService.ReadLatest()
	respMsg := fmt.Sprintf("{\"latest\": %d}", latest)

	jsonData := []byte(respMsg)
	w.Write(jsonData)
}

func (s *Simulator) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	s.log.Trace("Hit simulator register endpoint")

	err := s.updateLatest(r)
	if err != nil {
		s.log.Warnf("Failed to update latest with errror: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !s.simulatorService.IsAuthorized(w, r) {
		s.log.Warnf("A request to the simulator is not authorized")
		return
	}

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	var requestBody registerRequestBody
	err = json.Unmarshal(body, &requestBody)
	//Error handling if the struct doesn't get the necessary paramters for initialization
	if err != nil {
		s.log.Warnf("Failed to unmarshal request body into data object with error: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		jsonify := fmt.Sprintf("\"status\": %d, \"error_msg\": %s", 400, err.Error())
		w.Write([]byte(jsonify))
		return
	}

	var regError string
	if r.Method == http.MethodPost {
		if user, _ := s.userService.ReadUserByUsername(requestBody.Username); user.Username != "" {
			regError = "The username is already taken"
			w.WriteHeader(http.StatusBadRequest)
			jsonify := fmt.Sprintf("\"status\": %d, \"error_msg\": %s", 400, regError)
			w.Write([]byte(jsonify))
			return
		}
		err := s.userService.CreateUser(requestBody.Username, requestBody.Email, requestBody.Password)
		if err != nil {
			s.log.Warnf("An error occured during creation of a user with error: %s", err.Error())
		}
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(""))
	}
}

func (s *Simulator) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	s.log.Trace("Messages endpoint hit")

	err := s.updateLatest(r)
	if err != nil {
		s.log.Warnf("Failed to update latest with error: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !s.simulatorService.IsAuthorized(w, r) {
		s.log.Warnf("A request to the simulator is not authorized")
		return
	}

	if r.Method != http.MethodGet {
		s.log.Warnf("The messages endpoint received a non-GET request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filtered_msgs, err := s.messageService.ReadAllMessages(0, 100)
	if err != nil {
		s.log.Warnf("Failed reading all messages with error: %s", err.Error())
	}
	msgs, err := json.Marshal(filtered_msgs)
	if err != nil {
		s.log.Warnf("Failed marshalling messages to JSON with error: %s", err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write(msgs)
}

func (s *Simulator) UserPerMessageHandler(w http.ResponseWriter, r *http.Request) {
	s.log.Trace("Hit the messages per user endpoint")

	err := s.updateLatest(r)
	if err != nil {
		s.log.Warnf("Failed to update latest with error: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !s.simulatorService.IsAuthorized(w, r) {
		s.log.Warnf("A request to the simulator is not authorized")
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

	if r.Method == http.MethodGet {
		s.log.Trace("Handling a GET request for messages per user endpoint")
		s.postUserPerMessage(w, r, username)
	} else if r.Method == http.MethodPost {
		s.log.Trace("Handling a POST request for messages per user endpoint")
		s.getUserPerMessage(w, r, username)
	} else {
		s.log.Warn("Received a non-GET/POST request for message per user endpoint")
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
		s.log.Warnf("Failed to unmarshal request body into data object with error: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.messageService.CreateMessage(requestBody.Username, requestBody.Content)

	if err != nil {
		s.log.Warnf("Failed to create messages in db with error: %s", err.Error())
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}

func (s *Simulator) getUserPerMessage(w http.ResponseWriter, r *http.Request, username string) {
	filtered_msgs, err := s.messageService.ReadAllMessagesByUsername(username)
	if err != nil {
		s.log.Warnf("Failed to read messages by author with error: %s", err.Error())
	}

	msgs, err := json.Marshal(filtered_msgs)
	if err != nil {
		s.log.Warnf("Failed to marshall messages to JSON object with error: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(msgs)
	}
}

func (s *Simulator) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	s.log.Trace("Hit simulator follows endpoint")
	err := s.updateLatest(r)
	if err != nil {
		s.log.Warnf("Failed to update latest with error: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !s.simulatorService.IsAuthorized(w, r) {
		s.log.Warnf("A request to the simulator is not authorized")
		return
	}

	if r.Method == http.MethodPost {
		s.log.Trace("Received a POST request to the simulator follows endpoint")
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		var requestBody followRequestBody
		err := json.Unmarshal(body, &requestBody)
		if err != nil {
			s.log.Info("No follows JSON key was found, assuming unfollows exist")
			s.unfollowUser(w, r)
		} else {
			s.log.Trace("A follows JSON key was found")
			s.followUser(w, r, requestBody)
		}
	}

	if r.Method == http.MethodGet {
		s.log.Trace("Received a GET request to the simulator follows endpoint")
		s.getFollowers(w, r)
	}
}

func (s *Simulator) followUser(w http.ResponseWriter, r *http.Request, body followRequestBody) {
	s.log.Trace("Simulator is following a user")
	vars := mux.Vars(r)
	username := vars["username"]

	err := s.userService.Follow(username, body.Follow)
	if err != nil {
		s.log.Warnf("Failed to follow user with error: %s", err.Error())
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}

func (s *Simulator) unfollowUser(w http.ResponseWriter, r *http.Request) {
	s.log.Trace("Simulator is unfollowing a user")
	vars := mux.Vars(r)
	username := vars["username"]

	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	var requestBody unfollowRequestBody
	err := json.Unmarshal(body, &requestBody)
	if err != nil {
		s.log.Warnf("Requese body is invalid with error: %s", err.Error())
	}

	err = s.userService.Unfollow(username, requestBody.Unfollow)
	if err != nil {
		s.log.Warnf("Failed to unfollow user with error: %s", err.Error())
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}

func (s *Simulator) getFollowers(w http.ResponseWriter, r *http.Request) {
	s.log.Trace("Simulator is reading followers for a user")
	vars := mux.Vars(r)
	username := vars["username"]

	user, err := s.userService.ReadUserByUsername(username)
	if err != nil {
		s.log.Warnf("Failed to read all followers from username with error: %s", err.Error())
	}
	filteredFollowers := make([]string, 0)
	for _, entry := range user.Follows {
		filteredFollowers = append(filteredFollowers, entry.Username)
	}
	followers, err := json.Marshal(user)
	if err != nil {
		s.log.Warnf("Failed to marshall followers into JSON object with error: %s", err.Error())
	}
	w.Write(followers)
}

func (s *Simulator) updateLatest(r *http.Request) error {
	if !r.URL.Query().Has("latest") {
		s.log.Warnf("Trying to update latest but there is no latest query parameter in the URL")
		return nil
	}

	latest, err := parseLatest(r)
	if err != nil {
		s.log.Warnf("Tried to update latest but the new value is not an integer")
		return errors.New("latest was not an integer")
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

func (s *Simulator) SetupRoutes(r *mux.Router) {
	r.HandleFunc("/fllws/{username}", s.FollowUserHandler)
	r.HandleFunc("/register", s.RegisterHandler)
	r.HandleFunc("/msgs", s.MessagesHandler)
	r.HandleFunc("/msgs/{username}", s.UserPerMessageHandler)
	r.HandleFunc("/latest", s.LatestHandler)
}
