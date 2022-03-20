package main

import (
	"log"
	"net/http"

	"github.com/antonPalmFolkmann/DevOps2022/controllers"
	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/gorilla/mux"
)

func init() {

}

func main() {
	r := mux.NewRouter()

	db := storage.ConnectPsql()
	storage.Migrate(db)

	userService := services.NewUserService(db)
	messageService := services.NewMessageService(db)
	simService := services.NewSimulatorService()

	sim := controllers.NewSimulator(messageService, userService, simService)
	sim.SetupRoutes(r)
	log.Fatalln(http.ListenAndServe(":8081", r))
}

// func main() {
// 	r := mux.NewRouter()
// 	log.Fatalln(http.ListenAndServe(":8080", r))
// }
