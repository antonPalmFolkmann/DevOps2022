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

type Simulator struct {
	messageService   services.IMessage
	userService      services.IUser
	simulatorService services.ISimulatorService
}

func NewSimulator(messageService services.IMessage, userService services.IUser, simulatorService services.ISimulatorService) *Simulator {
	return &Simulator{messageService: messageService, userService: userService, simulatorService: simulatorService}
}

func (s *Simulator) LatestHandler(w http.ResponseWriter, r *http.Request) {
	latest := s.simulatorService.ReadLatest()
	respMsg := fmt.Sprintf("{\"latest\": %d}", latest)

	jsonData := []byte(respMsg)
	_, err := w.Write(jsonData)
	if err != nil {
		http.Error(w, "Failed to read latest", http.StatusInternalServerError)
	}
}

func (s *Simulator) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := s.updateLatest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !s.simulatorService.IsAuthorized(w, r) {
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
		_, err = w.Write([]byte(jsonify))
		if err != nil {
			http.Error(w, "Failed to write error message response", http.StatusInternalServerError)
		}
		return
	}

	var regError string
	if r.Method == http.MethodPost {
		if user, _ := s.userService.ReadUserByUsername(requestBody.Username); user.Username != "" {
			regError = "The username is already taken"
			w.WriteHeader(http.StatusBadRequest)
			jsonify := fmt.Sprintf("\"status\": %d, \"error_msg\": %s", 400, regError)
			_, err = w.Write([]byte(jsonify))
			if err != nil {
				http.Error(w, "Failed to write error message response", http.StatusInternalServerError)
			}
			return
		}
		err := s.userService.CreateUser(requestBody.Username, requestBody.Email, requestBody.Password)
		if err != nil {
			// w.WriteHeader(http.StatusInternalServerError)
			log.Println("simulator_controller: An error occured during creation of a user")
		}
		w.WriteHeader(http.StatusNoContent)
		_, err = w.Write([]byte(""))
		if err != nil {
			http.Error(w, "Failed to write response during registration of user", http.StatusInternalServerError)
		}
	}
}

func (s *Simulator) MessagesHandler(w http.ResponseWriter, r *http.Request) {
	err := s.updateLatest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !s.simulatorService.IsAuthorized(w, r) {
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filtered_msgs, err := s.messageService.ReadAllMessages(0, 100)
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
	_, err = w.Write(msgs)
	if err != nil {
		http.Error(w, "Failed to write messages to response", http.StatusInternalServerError)
	}
}

func (s *Simulator) UserPerMessageHandler(w http.ResponseWriter, r *http.Request) {
	err := s.updateLatest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !s.simulatorService.IsAuthorized(w, r) {
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

	err = s.messageService.CreateMessage(requestBody.Username, requestBody.Content)

	if err != nil {
		// w.WriteHeader(http.StatusBadRequest)
		// return
		log.Println("service_simulator: failed to create messages in db")
	}
	w.WriteHeader(http.StatusNoContent)
	_, err = w.Write([]byte(""))
	if err != nil {
		http.Error(w, "Failed to write response posting a message", http.StatusInternalServerError)
	}
}

func (s *Simulator) getUserPerMessage(w http.ResponseWriter, r *http.Request, username string) {
	filtered_msgs, err := s.messageService.ReadAllMessagesByUsername(username)
	if err != nil {
		log.Println("service_simulator: Failed to read messages by author -> Author doesn't exist")
	}
	msgs, err := json.Marshal(filtered_msgs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(msgs)
		if err != nil {
			http.Error(w, "Failed to write response during get user messages", http.StatusInternalServerError)
		}
	}
}

func (s *Simulator) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	err := s.updateLatest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !s.simulatorService.IsAuthorized(w, r) {
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

	err := s.userService.Follow(username, body.Follow)
	if err != nil {
		log.Println("service_simulator: Failed to follow user")
	}
	w.WriteHeader(http.StatusNoContent)
	_, err = w.Write([]byte(""))
	if err != nil {
		http.Error(w, "Failed to write response following user", http.StatusInternalServerError)
	}
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

	err = s.userService.Unfollow(username, requestBody.Unfollow)
	if err != nil {
		log.Println("service_simulator: Failed to unfollow user")
	}
	w.WriteHeader(http.StatusNoContent)
	_, err = w.Write([]byte(""))
	if err != nil {
		http.Error(w, "Failed to write response unfollowing user", http.StatusInternalServerError)
	}
}

func (s *Simulator) getFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	user, err := s.userService.ReadUserByUsername(username)
	if err != nil {
		log.Println("service_simulator: Failed to read all followers from username")
	}
	filteredFollowers := make([]string, 0)
	for _, entry := range user.Follows {
		filteredFollowers = append(filteredFollowers, entry.Username)
	}
	followers, err := json.Marshal(filteredFollowers)
	if err != nil {
		log.Println("service_simulator: Failed to marshall followers")
	}
	_, err = w.Write(followers)
	if err != nil {
		http.Error(w, "Failed to write response getting user followers", http.StatusInternalServerError)
	}
}

func (s *Simulator) updateLatest(r *http.Request) error {
	if !r.URL.Query().Has("latest") {
		return nil
	}

	latest, err := parseLatest(r)
	if err != nil {
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
