package main

import (
	"log"
	"net/http"
	"time"

	"github.com/antonPalmFolkmann/DevOps2022/controllers"
	"github.com/antonPalmFolkmann/DevOps2022/monitoring"
	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func init() {

}

func main() {
	time.Sleep(2 * time.Second)

	db := storage.ConnectPsql()
	storage.Migrate(db)

	userService := services.NewUserService(db)
	messageService := services.NewMessageService(db)
	simulatorService := services.NewSimulatorService()

	store := sessions.NewCookieStore([]byte("supersecret1234"))
	userController := controllers.NewUserController(userService, messageService, store)
	messageController := controllers.NewMessage(store, messageService, userService)
	serviceController := controllers.NewSimulator(messageService, userService, simulatorService)

	go func() {
		r := mux.NewRouter()
		monitoring.SetupRoutes(r)
		serviceController.SetupRoutes(r)
		err := http.ListenAndServe(":8081", r)
		if err != nil {
			log.Fatalf("Failed to listen and serve port: %s", err.Error())
		}
	}()

	r := mux.NewRouter()
	userController.SetupRoutes(r)
	monitoring.SetupRoutes(r)
	messageController.SetupRoutes(r)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Failed to listen and serve port: %s", err.Error())
	}
}

// func main() {
// 	r := mux.NewRouter()
// 	log.Fatalln(http.ListenAndServe(":8080", r))
// }
