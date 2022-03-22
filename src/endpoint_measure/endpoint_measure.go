package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	runs = 100
)

var (
	endpoints = map[string]func(){
		"public":     public,
		"register":   register,
		"login":      login,
		"logout":     logout,
		"addMessage": addMessage,
		"msgs":       msgs,
		"fllw":       fllw,
		"unfllw":     unfllw,
	}
	baseUrl = "http://localhost:8080"
)

func main() {
	for e, fn := range endpoints {
		log.Println(e, ":")
		testEndpoint(fn)
	}
}

func testEndpoint(endpoint func()) {
	var totalDur float64

	for i := 0; i < runs; i++ {
		start := time.Now()
		endpoint()
		totalDur += time.Since(start).Seconds()
	}

	log.Print(totalDur / float64(runs))
}

func public() {
	url := baseUrl + "/public"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
}

func register() {
	url := baseUrl + "/register"

	payload := strings.NewReader("{ \"username\": \"sendit\", \"password\": \"sender\", \"email\": \"send@send.com\" }")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
}

func login() {
	url := baseUrl + "/login"

	payload := strings.NewReader("{ \"username\": \"sendit\", \"password\": \"sender\" }")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
}

func logout() {
	url := baseUrl + "/logout"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
}

func addMessage() {
	url := baseUrl + "/add_message"

	payload := strings.NewReader("{ \"authorName\": \"sendit\", \"text\": \"RNSK\" }")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
}

func msgs() {
	url := baseUrl + "/msgs/frick"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
}

func fllw() {
	url := baseUrl + "/fllw/frick"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
}

func unfllw() {
	url := baseUrl + "/unfllw/frick"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
}
