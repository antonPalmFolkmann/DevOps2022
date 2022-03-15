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
	UpdateUser(ID uint, username string, email string, password string) error
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
	var users []storage.User
	err := u.db.Select([]string{"id", "username", "email", "pw_hash"}).
				Find(&users).Error
	userDTOs := make([]storage.UserDTO, 0)
	for _, v := range users {
		userDTO := storage.UserDTO{ID: v.ID, Username: v.Username, Email: v.Email, PwHash: v.PwHash}
		userDTOs = append(userDTOs, userDTO)
	}
	return userDTOs, err
}

func (u *User) ReadUserById(id uint) (storage.UserDTO, error) {
	var user storage.User
	err := u.db.Unscoped().
				Where("id = ?", id).
				Select([]string{"id", "username", "email", "pw_hash"}).
				Find(&user).Error

	userDTO := storage.UserDTO{ID: user.ID, Username: user.Username, Email: user.Email, PwHash: user.PwHash}
	return userDTO, err
}

func (u *User) ReadUserByUsername(username string) (storage.UserDTO, error) {
	var user storage.User
	err := u.db.Unscoped().
				Where("username = ?", username).
				Select([]string{"id", "username", "email", "pw_hash"}).
				Find(&user).Error

	userDTO := storage.UserDTO{ID: user.ID, Username: user.Username, Email: user.Email, PwHash: user.PwHash}
	return userDTO, err
}

func (u *User) ReadUserIdByUsername(username string) (uint, error) {
	var user storage.User
	err := u.db.Unscoped().
				Where("username = ?", username).
				Select("id").
				Find(&user).Error
	return user.ID, err
}

func (u *User) UpdateUser(ID uint, username string, email string, password string) error {
	var user storage.User
	PwHash := u.Hash(password)
	err := u.db.Model(&user).
				Unscoped().
				Where("id = ?", ID).
				Update(&storage.User{Username: username, Email: email, PwHash: PwHash}).Error
	return err
} 

func (u *User) DeleteUser(ID uint) error {
	var user storage.User
	err := u.db.Unscoped().
				Where("id = ?", ID).
				Delete(&user).Error
	return err
}

func (u *User) Hash(password string) string {
	hash := md5.New()
	io.WriteString(hash, password)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
