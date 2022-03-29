package main

import (
	"log"
	"net/http"
	"os"
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

	// Testing logging

	//create your file with desired read/write permissions
	f, err := os.OpenFile("minitwit.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()

	//set output of logs to f
	log.SetOutput(f)

	//test case
	log.Println("check to make sure it works")

	// Starting app

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
		http.ListenAndServe(":8081", r)
	}()

	r := mux.NewRouter()
	userController.SetupRoutes(r)
	monitoring.SetupRoutes(r)
	messageController.SetupRoutes(r)
	http.ListenAndServe(":8080", r)
}

// func main() {
// 	r := mux.NewRouter()
// 	log.Fatalln(http.ListenAndServe(":8080", r))
// }
