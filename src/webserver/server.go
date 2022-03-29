package main

import (
	"log"
	"log/syslog"
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

	// Log to syslog
    logWriter, err := syslog.New(syslog.LOG_SYSLOG, "My Awesome App")
    if err != nil {
        log.Fatalln("Unable to set logfile:", err.Error())
    }

    // + set log flag
    log.SetFlags(log.Lshortfile)

    // set the log output
    log.SetOutput(logWriter)

    log.Println("This is a log from GOLANG")

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
