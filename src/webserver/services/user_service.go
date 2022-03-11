package services

import (
	"crypto/md5"
	"fmt"
	"io"

	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/jinzhu/gorm"
)

type IUser interface {
	CreateUser(username string, email string, password string) error
	ReadAllUsers() ([]storage.UserDTO, error)
	ReadUserById(ID uint) (storage.UserDTO, error)
	ReadUserByUsername(username string) (storage.UserDTO, error)
	ReadUserIdByUsername(username string) (uint, error)
	UpdateUser(ID uint, username string, email string, pwHash string) error
	DeleteUser(ID uint) error
	Hash(password string) string
}

type User struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *User {
	return &User{db: db}
}

func (u *User) CreateUser(username string, email string, password string) error {
	pwHash := u.Hash(password)
	user := storage.User{Username: username, Email: email, PwHash: pwHash}
	err := u.db.Create(&user).Error
	return err
}

func (u *User) ReadAllUsers() ([]storage.UserDTO, error) {
	var users = make([]storage.UserDTO, 0)
	err := u.db.Select("user_id", "username", "email", "pw_hash").Find(&users).Error
	return users, err
}

func (u *User) ReadUserById(id uint) (storage.UserDTO, error) {
	var user storage.UserDTO
	err := u.db.Select("user_id", "username", "email", "pw_hash", "messages", "follows").
				Where("user_id = ?", id).
				Find(&user).Error
	return user, err
}

func (u *User) ReadUserByUsername(username string) (storage.UserDTO, error) {
	var user storage.UserDTO
	err := u.db.Select("user_id", "username", "email", "pw_hash", "messages", "follows").
				Where("username = ?", username).
				Find(&user).Error
	return user, err
}

func (u *User) ReadUserIdByUsername(username string) (uint, error) {
	var user storage.UserDTO
	err := u.db.Select("user_id", "username", "email", "pw_hash", "messages", "follows").
				Where("username = ?", username).
				Find(&user).Error
	return user.ID, err
}

func (u *User) UpdateUser(ID uint, username string, email string, pwHash string) error {
	var user storage.User
	err := u.db.Where("user_id = ?", ID).
				Find(&user).Error
	if err != nil {
		return err
	}
	user.Username = username
	user.Email = email
	user.PwHash = pwHash
	err = u.db.Save(&user).Error
	return err
}

func (u *User) DeleteUser(ID uint) error {
	var user storage.User
	err := u.db.Delete(&user, ID).Error
	return err
}

func (u *User) Hash(password string) string {
	hash := md5.New()
	io.WriteString(hash, password)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
