package controllers

import (
	"errors"
	"fmt"
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
