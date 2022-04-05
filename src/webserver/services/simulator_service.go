package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ISimulatorService interface {
	IsAuthorized(w http.ResponseWriter, r *http.Request) bool
	ReadLatest() int
	UpdateLatest(latest int)
}

type SimulatorService struct {
	latest int
	log    *logrus.Logger
}

func NewSimulatorService(log *logrus.Logger) *SimulatorService {
	return &SimulatorService{log: log}
}

func (s *SimulatorService) IsAuthorized(w http.ResponseWriter, r *http.Request) bool {
	s.log.Trace("Authorizing the request")

	authorizedReq := r.Header.Get("Authorization")
	if authorizedReq != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
		s.log.Info("Simulator request is not authorized")
		error := "You are not authorized to use this resource!"
		jsonify := fmt.Sprintf("\"status\": %d, \"error_msg\": %s", 403, error)
		_, err := w.Write([]byte(jsonify))
		if err != nil {
			log.Fatalf("Failed to ")
		}
		return false
	}
	s.log.Info("Simulator request is authorized")
	return true
}

func (s *SimulatorService) ReadLatest() int {
	s.log.Trace("Reading latest for the simulator")
	s.log.Debugf("Latest is %d", s.latest)
	return s.latest
}

func (s *SimulatorService) UpdateLatest(latest int) {
	s.log.Trace("Updating latest for the simulator")
	s.latest = latest
	s.log.Debugf("Latest is %d", s.latest)
}
