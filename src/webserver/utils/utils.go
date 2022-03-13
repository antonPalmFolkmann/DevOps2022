package utils

import (
	"errors"
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
