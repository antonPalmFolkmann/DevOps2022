package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	User_id   	int64	`gorm:"type:bigint;" json:"user_id"`
	Username  	string	`gorm:"type:varchar(255);" json:"username"`
	Email     	string	`gorm:"type:varchar(255);" json:"email"`	
	Pw_hash   	string	`gorm:"type:varchar(255);" json:"pw_hash"`
}

func GetUser() User {
	var user User
	return user
}

func GetUsers() []User {
	var users []User
	return users
}