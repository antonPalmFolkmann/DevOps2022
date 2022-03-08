package services

import (
	"github.com/antonPalmFolkmann/DevOps2022/models"
	"github.com/jinzhu/gorm"
)

var dbconn *gorm.DB

type Response struct {
	Data []models.User `json:"data"`
	Message string `json:"message"`
}

func SetDB(db *gorm.DB) {
	dbconn = db
	var user = models.GetUser()
	dbconn.AutoMigrate(&user)
}