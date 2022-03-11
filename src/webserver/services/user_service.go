package services

import (
	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/jinzhu/gorm"
)

type IUser interface {
	CreateUser(username string, email string, pwHash string) error
	ReadAllUsers() ([]storage.UserDTO, error)
	ReadUserById(ID uint) (storage.UserDTO, error)
	ReadUserByUsername(username string) (storage.UserDTO, error)
	ReadUserIdByUsername(username string) (uint, error)
	UpdateUser(ID uint, username string, email string, pwHash string) error
	DeleteUser(ID uint) error
	Hash(password string) string
	Unhash(hash string) string
}

type User struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *User {
	return &User{db: db}
}

func (u *User) CreateUser(user storage.User) error {
	err := u.db.Create(&user).Error
	return err
}

func (u *User) ReadAllUsers() ([]storage.User, error) {
	var users = make([]storage.User, 0)
	err := u.db.Find(&users).Error
	return users, err
}

func (u *User) ReadUserById(id uint) (storage.User, error) {
	var user storage.User
	err := u.db.Where("user_id = ?", id).Find(&user).Error
	return user, err
}

func (u *User) ReadUserByUsername(username string) (storage.User, error) {
	var user storage.User
	err := u.db.Where("username = ?", username).Find(&user).Error
	return user, err
}

func (u *User) ReadUserIdByUsername(username string) (uint, error) {
	var user storage.User
	err := u.db.Where("username = ?", username).Find(&user).Error
	return user.ID, err
}

func (u *User) UpdateUser(user storage.User, id uint) error {
	err := u.db.Model(&user).Where("user_id = ?", id).Update(&user).Error
	return err
}

func (u *User) DeleteUser(id uint) error {
	var user storage.User
	err := u.db.Delete(&user, id).Error
	return err
}
