package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/gorilla/mux"
)

/*
	GET, POST, DELETE. 


	Er vi enige om at UserService er det der svarer til e.g. UserRepository?
	Tror jeg vil prøve at få interfacet ind i userService og få det til at ligne
	vores BDSA projekt. Yes, lad os prøve
	Enig, men jeg ved ikke helt hvordan vi gør. Vi kan prøve herfra

		https://gobyexample.com/interfaces
		https://betterprogramming.pub/implementing-interfaces-with-golang-51a3b7f527b4
	I think so, not sure
	Tror det er nemmere at ændre interface nu, but dunno
*/

type IUser interface {
	ReadUserById(id int64) *User
	ReadUserByUsername(username string) *User
	CreateUser(username string, email string, pw_hash string) error
	GetUserIdByUsername(username string) int64
	DeleteUser(id int64) error
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var resp Response
	users, err := services.GetAllUsers()
	if err == nil {
		log.Println(users)
		resp.Data = users
		resp.Message = "SUCCESS"
		json.NewEncoder(w).Encode(&resp)
	} else {
		log.Println(err)
		http.Error(w, err.Error(), 400)
	}
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["user_id"])
	var resp Response
	user, err := services.GetUserByID(id)
	if err == nil {
		log.Println(user)
		resp.Data = append(resp.Data, user)
		resp.Message = "SUCCESS"
		json.NewEncoder(w).Encode(&resp)
	} else {
		log.Println(err)
		http.Error(w, err.Error(), 400)
	}
}

func GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	username := params["username"]
	var resp Response
	user, err := services.GetUserByUsername(username)
	if err == nil {
		log.Println(user)
		resp.Data = append(resp.Data, user)
		resp.Message = "SUCCESS"
		json.NewEncoder(w).Encode(&resp)
	} else {
		log.Println(err)
		http.Error(w, err.Error(), 400)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var resp Response
	err := services.CreateUser(r)
	if err != nil {
		http.Error(w, "Error Creating Record", 400)
		return
	}
	resp.Message = "CREATED"
	json.NewEncoder(w).Encode(resp)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["user_id"])
	var resp Response
	err := services.UpdateUser(r, id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	resp.Message = "UPDATED"
	json.NewEncoder(w).Encode(resp)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["user_id"])
	var resp Response
	err := services.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	resp.Message = "DELETED"
	json.NewEncoder(w).Encode(resp)
}