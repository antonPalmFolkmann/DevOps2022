package routes

import (
	//controllers "github.com/antonPalmFolkmann/DevOps2022/Controllers"
	"github.com/gorilla/mux"
)


func CreateRouter() *mux.Router {
	router := mux.NewRouter()
	/* router.HandleFunc("/users", UserController.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", controllers.GetUserByID).Methods("GET")
	router.HandleFunc("/users/{username}", controllers.GetUserByUsername).Methods("GET")
	router.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", controllers.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", controllers.DeleteUser).Methods("DELETE") */
	return router
}