package services

import (
	"github.com/antonPalmFolkmann/DevOps2022/models"
	"github.com/jinzhu/gorm"
)

var dbconn *gorm.DB

func SetDB(db *gorm.DB) {
	dbconn = db
	var user = models.GetUser()
	var message = models.GetMessage()
	dbconn.AutoMigrate(&user)
	dbconn.AutoMigrate(&message)
}