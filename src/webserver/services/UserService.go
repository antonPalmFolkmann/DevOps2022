package services

import (
	"github.com/antonPalmFolkmann/DevOps2022/models"
)

type IUserService interface {
	ReadAllUsers() 							([]models.User, error)
	ReadUserById(id int) 					(models.User, error)
	ReadUserByUsername(username string) 	(models.User, error)
	ReadUserIdByUsername(username string) 	(int64, error)
	CreateUser(user models.User) 			(error)
	UpdateUser(user models.User, id int)	(error)
	DeleteUser(id int) 						(error)
}

type UserService struct {}

func (u UserService) ReadAllUsers() ([]models.User, error) {
	var users = models.GetUsers()
	err := dbconn.Find(&users).Error
	return users, err
}

func (u UserService) ReadUserById(id int) (models.User, error) {
	var user = models.GetUser()
	err := dbconn.Where("user_id = ?", id).Find(&user).Error
	return user, err
}


func (u UserService) ReadUserByUsername(username string) (models.User, error) {
	var user = models.GetUser()
	err := dbconn.Where("username = ?", username).Find(&user).Error
	return user, err
}

func (u UserService) ReadUserIdByUsername(username string) (int64, error) {
	var user = models.GetUser()
	err := dbconn.Where("username = ?", username).Find(&user).Error
	return user.User_id, err
}

func (u UserService) CreateUser(user *models.User) error {
	err := dbconn.Create(&user).Error
	return err
}

func (u UserService) UpdateUser(user *models.User, id int) error {
	err := dbconn.Model(&user).Where("user_id = ?", id).Update(&user).Error
	return err
}

func (u UserService) DeleteUser(id int) error {
	var user = models.GetUser()
	err := dbconn.Delete(&user, id).Error
	return err
}