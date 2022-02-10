package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	mat "github.com/mattn/go-sqlite3"

	"github.com/gorilla/mux"
)

func YourHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

func main() {
	db, err := sql.Open("sqlite3", "./minitwit.db")

	if err != nil {
		r := mux.NewRouter()

		// Routes consist of a path and a handler function.
		r.HandleFunc("/", YourHandler)
		db.Begin()

		// Bind to a port and pass our router in
		log.Fatal(http.ListenAndServe(":8080", r))
	}
	matString := mat.ErrAbort.Error()
	fmt.Printf(matString)

}
