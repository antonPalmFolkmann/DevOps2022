package services

import (
	"log"
	"net/http"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/antonPalmFolkmann/DevOps2022/models"
)

var dbconn *gorm.DB

type Response struct {
	Data []models.User `json:"data"`
	Message string `json:"message"`
}

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

func SetDB(db *gorm.DB) {
	dbconn = db
	var user = models.GetUser()
	dbconn.AutoMigrate(&user)
}