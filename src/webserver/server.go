package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/antonPalmFolkmann/DevOps2022/controllers"
	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func init() {

}

func main() {
	time.Sleep(5 * time.Second)
	db := storage.ConnectPsql()
	defer db.Close()
	storage.Migrate(db)

	store := sessions.NewCookieStore([]byte(os.Getenv("SECURE_COOKIE_KEY")))

	userService := services.NewUserService(db)
	messageService := services.NewMessageService(db)

	userController := controllers.NewUserController(userService, messageService, store)

	r := mux.NewRouter()
	userController.SetupRoutes(r)
	log.Fatalln(http.ListenAndServe(":8080", r))
}
