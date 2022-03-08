package routes

import (
	"github.com/gorilla/mux"
	"github.com/antonPalmFolkmann/DevOps2022/services"
)


func CreateRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users", services.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", services.GetUserByID).Methods("GET")
	router.HandleFunc("/users/{username}", services.GetUserByUsername).Methods("GET")
	router.HandleFunc("/users", services.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", services.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", services.DeleteUser).Methods("DELETE")
	return router
}