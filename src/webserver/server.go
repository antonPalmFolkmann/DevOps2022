package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func init() {

}

func main() {
	r := mux.NewRouter()
	log.Fatalln(http.ListenAndServe(":8080", r))
}
