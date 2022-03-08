package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/gorilla/mux"
)

/*
	GET, POST, DELETE.

	https://gobyexample.com/interfaces
	https://betterprogramming.pub/implementing-interfaces-with-golang-51a3b7f527b4
*/

type IUserController interface {
	GetAllUsers(http.ResponseWriter, http.Request)
	ReadUserByUsername(http.ResponseWriter, http.Request)
	CreateUser(http.ResponseWriter, http.Request)
	GetUserIdByUsername(http.ResponseWriter, http.Request)
	DeleteUser(http.ResponseWriter, http.Request)
}

type UserController struct {
	userService services.IUserService
}

func NewUserController(userService services.IUserService) *UserController {
	return &UserController{userService: userService}
}

func (u *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var resp Response
	users, err := u.userService.ReadAllUsers()
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

func (u *UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["user_id"])
	var resp Response
	user, err := u.userService.ReadUserById(id)
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

func (u *UserController) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	username := params["username"]
	var resp Response
	user, err := u.userService.ReadUserByUsername(username)
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

func (u *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var resp Response
	var user storage.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	log.Println(user)
	err := u.userService.CreateUser(user)
	if err != nil {
		http.Error(w, "Error Creating Record", 400)
		return
	}
	resp.Message = "CREATED"
	json.NewEncoder(w).Encode(resp)
}

func (u *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["user_id"])
	var resp Response
	var user storage.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	err := u.userService.UpdateUser(user, id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	resp.Message = "UPDATED"
	json.NewEncoder(w).Encode(resp)
}

func (u *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["user_id"])
	var resp Response
	err := u.userService.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	resp.Message = "DELETED"
	json.NewEncoder(w).Encode(resp)
}
