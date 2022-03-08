package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/antonPalmFolkmann/DevOps2022/models"
	"github.com/gorilla/mux"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users = models.GetUsers()
	var resp Response
	err := dbconn.Find(&users).Error
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
	var user = models.GetUser()
	err := dbconn.Where("user_id = ?",  id).Find(&user).Error
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
	var user = models.GetUser()
	err := dbconn.Where("username = ?", username).Find(&user).Error
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