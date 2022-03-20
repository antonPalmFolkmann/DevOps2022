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
	ReadAllUsers() ([]storage.User, error)
	ReadUserByUsername(username string) (storage.User, error)
	ReadUserIdByUsername(username string) (uint, error)
	Follow(username string, whomname string) error
	Unfollow(username string, whomname string) error
	IsPasswordCorrect(username string, password string) bool
	IsUsernameTaken(username string) bool
}

type User struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *User {
	return &User{db: db}
}

func (u *User) CreateUser(username string, email string, password string) error {
	pwHash := u.hash(password)
	user := storage.User{Username: username, Email: email, PwHash: pwHash}
	err := u.db.Create(&user).Error
	return err
}

func (u *User) ReadAllUsers() ([]storage.User, error) {
	var users []storage.User
	err := u.db.Select([]string{"id", "username", "email", "pw_hash"}).
		Find(&users).Error
	return users, err
}

func (u *User) ReadUserByUsername(username string) (storage.User, error) {
	var user storage.User
	err := u.db.Where("username = ?", username).
				Select([]string{"id", "username", "email", "pw_hash"}).
				Find(&user).Error
	return user, err
}

func (u *User) ReadUserIdByUsername(username string) (uint, error) {
	var user storage.User
	err := u.db.Where("username = ?", username).
				Select("id").
				Find(&user).Error
	return user.ID, err
}

func (u *User) hash(password string) string {
	hash := md5.New()
	io.WriteString(hash, password)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (u *User) Follow(username string, whomname string) error {
	var user storage.User
	err := u.db.
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return err
	}

	var whom storage.User
	err = u.db.
		Where("username = ?", whomname).
		First(&whom).Error
	if err != nil {
		return err
	}

	user.Follows = append(user.Follows, &whom)
	u.db.Save(&user)
	return nil
}

func (u *User) Unfollow(username string, whomname string) error {
	user, err := u.ReadUserByUsername(username)
	if err != nil {
		return err
	}

	whom, err := u.ReadUserByUsername(whomname)
	if err != nil {
		return err
	}

	u.db.Exec("DELETE FROM follows WHERE user_id = ? AND whom_id = ?", user.ID, whom.ID)
	return nil
}

func (u *User) IsPasswordCorrect(username string, password string) bool {
	passwordHashed := u.hash(password)
	usr, err := u.ReadUserByUsername(username)
	if err != nil {
		return false
	}
	
	return (usr.PwHash == passwordHashed)
}

func (u *User) IsUsernameTaken(username string) bool {
	_, err := u.ReadUserByUsername(username)
	return err == nil 
}
