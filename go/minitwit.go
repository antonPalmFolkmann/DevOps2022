package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func YourHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

const dbFile = "../minitwit.db"

func main() {
	_, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Failed to connect to the database with error: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", YourHandler)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8080", r))
}
