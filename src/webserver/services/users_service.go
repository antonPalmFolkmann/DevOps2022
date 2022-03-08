package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/antonPalmFolkmann/DevOps2022/models"
	"github.com/gorilla/mux"
)

type IUser interface {
	ReadUserById(id int64) *User
	ReadUserByUsername(username string) *User
	CreateUser(username string, email string, pw_hash string) error
	GetUserIdByUsername(username string) int64
	DeleteUser(id int64) error
}

func ControllerGetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var resp Response
	users, err := ServiceGetAllUsers()
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

func ServiceGetAllUsers() ([]models.User, error) {
	var users = models.GetUsers()
	err := dbconn.Find(&users).Error
	return users, err
}

func ControllerGetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["user_id"])
	var resp Response
	user, err := ServicegetUserByID(id)
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

func ServicegetUserByID(id int) (models.User, error) {
	var user = models.GetUser()
	err := dbconn.Where("user_id = ?", id).Find(&user).Error
	return user, err
}

func ControllerGetUserByUsername(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	username := params["username"]
	var resp Response
	user, err := ServiceGetUserByUsername(username)
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

func ServiceGetUserByUsername(username string) (models.User, error) {
	var user = models.GetUser()
	err := dbconn.Where("username = ?", username).Find(&user).Error
	return user, err
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var resp Response
	var user = models.GetUser()
	_ = json.NewDecoder(r.Body).Decode(&user)
	log.Println(user)

	err := dbconn.Create(&user).Error
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
	var resp Response
	var user = models.GetUser()
	_ = json.NewDecoder(r.Body).Decode(&user)
	id, _ := strconv.Atoi(params["user_id"])
	
	err := dbconn.Model(&user).Where("user_id = ?", id).Update(&user).Error
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
	var resp Response
	var user = models.GetUser()
	err := dbconn.Delete(&user, params["user_id"]).Error
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	resp.Message = "DELETED"
	json.NewEncoder(w).Encode(resp)
}