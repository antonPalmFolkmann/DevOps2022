package services

import (
	"fmt"
	"log"
	"net/http"
)

type ISimulatorService interface {
	IsAuthorized(w http.ResponseWriter, r *http.Request) bool
	ReadLatest() int
	UpdateLatest(latest int)
}

type SimulatorService struct {
	latest int
}

func NewSimulatorService() *SimulatorService {
	return &SimulatorService{}
}

func (s *SimulatorService) IsAuthorized(w http.ResponseWriter, r *http.Request) bool {
	authorizedReq := r.Header.Get("Authorization")
	if authorizedReq != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
		error := "You are not authorized to use this resource!"
		jsonify := fmt.Sprintf("\"status\": %d, \"error_msg\": %s", 403, error)
		_, err := w.Write([]byte(jsonify))
		if err != nil {
			log.Fatalf("Failed to ")
		}
		return false
	}
	return true
}

func (s *SimulatorService) ReadLatest() int {
	return s.latest
}

func (s *SimulatorService) UpdateLatest(latest int) {
	s.latest = latest
}
