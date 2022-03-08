package main

import (
	"log"
	"net/http"

	"github.com/antonPalmFolkmann/DevOps2022/database"
	"github.com/antonPalmFolkmann/DevOps2022/routes"
	"github.com/antonPalmFolkmann/DevOps2022/services"
)

func main() {
	
	/*
	go func() {
		r := mux.NewRouter()
		simulator.SetupRoutes(r)
		log.Fatalln(http.ListenAndServe(":8081", r))
	}()

	_ = database.GetConnection()

	// Setup minitwit "website"
	r := mux.NewRouter()
	minitwit.SetupRoutes(r)
	log.Fatalln(http.ListenAndServe(":8080", r))
	*/

	var db = database.GetConnection()
	services.SetDB(db)
	var appRouter = routes.CreateRouter()
	
	log.Println("Listening on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", appRouter))
}
