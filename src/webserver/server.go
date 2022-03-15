package main

import (
	"log"
	"time"

	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/antonPalmFolkmann/DevOps2022/storage"
)

func main() {
	time.Sleep(5 * time.Second)
	db := storage.ConnectPsql()
	defer db.Close()

	storage.Migrate(db)

	UserService := *services.NewUserService(db)
	user, err := UserService.ReadUserByUsername("frick")

	if err != nil {
		log.Panicf("Error: " + err.Error())
	}

	log.Println(user)

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
