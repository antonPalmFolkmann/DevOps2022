package main

import (
	"log"
	"time"

	"github.com/antonPalmFolkmann/DevOps2022/storage"
)

func main() {
	time.Sleep(5 * time.Second)
	db := storage.ConnectPsql()
	defer db.Close()

	storage.Migrate(db)

	var user storage.User
	db.First(&user, 1)
	log.Println(user)

	log.Println()

	// go func() {
	// 	r := mux.NewRouter()
	// 	simulator.SetupRoutes(r)
	// 	log.Fatalln(http.ListenAndServe(":8081", r))
	// }()

	// // Setup minitwit "website"
	// r := mux.NewRouter()
	// minitwit.SetupRoutes(r)
	// log.Fatalln(http.ListenAndServe(":8080", r))
}
