package services

import (
	"github.com/antonPalmFolkmann/DevOps2022/models"
)

type IUserService interface {
	CreateUser(user models.User) error
	ReadAllUsers() ([]models.User, error)
	ReadUserById(id int) (models.User, error)
	ReadUserByUsername(username string) (models.User, error)
	ReadUserIdByUsername(username string) (int64, error)
	UpdateUser(user models.User, id int) error
	DeleteUser(id int) error
}

type UserService struct{}

func (u MessageService) CreateUser(user *models.User) error {
	err := dbconn.Create(&user).Error
	return err
}

func (u MessageService) ReadAllUsers() ([]models.User, error) {
	var users = make([]models.User, 0)
	err := dbconn.Find(&users).Error
	return users, err
}

func (u MessageService) ReadUserById(id int) (models.User, error) {
	var user models.User
	err := dbconn.Where("user_id = ?", id).Find(&user).Error
	return user, err
}

func (u MessageService) ReadUserByUsername(username string) (models.User, error) {
	var user models.User
	err := dbconn.Where("username = ?", username).Find(&user).Error
	return user, err
}

func (u MessageService) ReadUserIdByUsername(username string) (int64, error) {
	var user models.User
	err := dbconn.Where("username = ?", username).Find(&user).Error
	return user.User_id, err
}

func (u MessageService) UpdateUser(user *models.User, id int) error {
	err := dbconn.Model(&user).Where("user_id = ?", id).Update(&user).Error
	return err
}

func (u MessageService) DeleteUser(id int) error {
	var user models.User
	err := dbconn.Delete(&user, id).Error
	return err
}
