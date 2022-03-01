package main

import (
	"log"
	"net/http"

	"github.com/antonPalmFolkmann/DevOps2022/minitwit"
	"github.com/antonPalmFolkmann/DevOps2022/simulator"
	"github.com/gorilla/mux"
)

func main() {
	go func() {
		r := mux.NewRouter()
		simulator.SetupRoutes(r)
		log.Fatalln(http.ListenAndServe("localhost:8081", r))
	}()

	// Setup minitwit "website"
	r := mux.NewRouter()
	minitwit.SetupRoutes(r)
	log.Fatalln(http.ListenAndServe("localhost:8080", r))
}