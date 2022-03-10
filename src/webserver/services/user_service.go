package services

import (
	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/jinzhu/gorm"
)

type IUserService interface {
	CreateUser(user storage.User) error
	ReadAllUsers() ([]storage.User, error)
	ReadUserById(id uint) (storage.User, error)
	ReadUserByUsername(username string) (storage.User, error)
	ReadUserIdByUsername(username string) (uint, error)
	UpdateUser(user storage.User, id uint) error
	DeleteUser(id uint) error
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (u *UserService) CreateUser(user storage.User) error {
	err := u.db.Create(&user).Error
	return err
}

func (u *UserService) ReadAllUsers() ([]storage.User, error) {
	var users = make([]storage.User, 0)
	err := u.db.Find(&users).Error
	return users, err
}

func (u *UserService) ReadUserById(id uint) (storage.User, error) {
	var user storage.User
	err := u.db.Where("user_id = ?", id).Find(&user).Error
	return user, err
}

func (u *UserService) ReadUserByUsername(username string) (storage.User, error) {
	var user storage.User
	err := u.db.Where("username = ?", username).Find(&user).Error
	return user, err
}

func (u *UserService) ReadUserIdByUsername(username string) (uint, error) {
	var user storage.User
	err := u.db.Where("username = ?", username).Find(&user).Error
	return user.ID, err
}

func (u *UserService) UpdateUser(user storage.User, id uint) error {
	err := u.db.Model(&user).Where("user_id = ?", id).Update(&user).Error
	return err
}

func (u *UserService) DeleteUser(id uint) error {
	var user storage.User
	err := u.db.Delete(&user, id).Error
	return err
}
