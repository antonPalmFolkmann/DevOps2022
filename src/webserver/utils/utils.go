package utils

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func FormatDatetime(timestamp int64) string {
	timeUnix := time.Unix(timestamp, 0)
	return timeUnix.Format("2006-01-02 15:04")
}

func ParseUsername(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	if username, found := vars["username"]; !found {
		return "", errors.New("there is no username")
	} else {
		return username, nil
	}
}

func Check_if_test_fail(err error) {
	if err != nil {
		log.Fatalf("An error occured during test: %s", err.Error())
	}
}
