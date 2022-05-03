package main

import (
	"net/http"

	"github.com/antonPalmFolkmann/DevOps2022/controllers"
	"github.com/antonPalmFolkmann/DevOps2022/monitoring"
	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetLevel(logrus.TraceLevel)
	db := storage.ConnectPsql()
	storage.Migrate(db)

	userService := services.NewUserService(db, log)
	messageService := services.NewMessageService(db, log)
	simulatorService := services.NewSimulatorService(log)

	store := sessions.NewCookieStore([]byte("supersecret1234"))
	userController := controllers.NewUserController(userService, messageService, store, log)
	messageController := controllers.NewMessage(store, messageService, userService, log)
	serviceController := controllers.NewSimulator(messageService, userService, simulatorService, log)

	log.Println("Pre go func")

	go func() {
		log.Trace("Starting the simulator router")
		r := mux.NewRouter()
		monitoring.SetupRoutes(r)
		serviceController.SetupRoutes(r)
		err := http.ListenAndServe(":8081", r)
		if err != nil {
			log.Fatalf("Failed to listen and serve port: %s", err.Error())
		}
	}()

	log.Trace("Starting the minitwit router")
	r := mux.NewRouter()
	userController.SetupRoutes(r)
	monitoring.SetupRoutes(r)
	messageController.SetupRoutes(r)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Failed to listen and serve port: %s", err.Error())
	}
}
