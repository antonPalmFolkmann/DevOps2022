package routes

import (
	"github.com/gorilla/mux"
	"github.com/antonPalmFolkmann/DevOps2022/services"
)


func CreateRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users", services.GetAllUsers).Methods("GET")
	return router
}