package services

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"

	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
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
	db  *gorm.DB
	log *logrus.Logger
}

func NewUserService(db *gorm.DB, log *logrus.Logger) *User {
	return &User{db: db, log: log}
}

func (u *User) CreateUser(username string, email string, password string) error {
	u.log.Trace("Creating a user")

	pwHash := u.hash(password)
	user := storage.User{Username: username, Email: email, PwHash: pwHash}
	err := u.db.Create(&user).Error
	u.log.Debug("Created a user on the database")
	return err
}

func (u *User) ReadAllUsers() ([]storage.User, error) {
	u.log.Trace("Reading all users")

	var users []storage.User
	err := u.db.Select([]string{"id", "username", "email", "pw_hash"}).
		Find(&users).Error

	return users, err
}

func (u *User) ReadUserByUsername(username string) (storage.User, error) {
	u.log.Trace("Reading user by username")

	var user storage.User
	err := u.db.Where("username = ?", username).
		Select([]string{"id", "username", "email", "pw_hash"}).
		Find(&user).Error
	u.log.Debug("Read user database")
	return user, err
}

func (u *User) ReadUserIdByUsername(username string) (uint, error) {
	u.log.Trace("Reading user ID by username")

	var user storage.User
	err := u.db.Where("username = ?", username).
		Select("id").
		Find(&user).Error
	u.log.Debug("Reader user from database, returning their ID")
	return user.ID, err
}

func (u *User) hash(password string) string {
	u.log.Trace("Hashing a password")

	hash := md5.New()
	_, err := io.WriteString(hash, password)
	if err != nil {
		log.Fatalf("Failed to hash password: %s", err)
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (u *User) Follow(username string, whomname string) error {
	u.log.Trace("Following a user")

	var user storage.User
	err := u.db.
		Where("username = ?", username).
		First(&user).Error
	u.log.Debug("Read the user wanting to follow from the database")
	if err != nil {
		u.log.Warn("Unable to find the user wanting to follow another user with error: ", err.Error())
		return err
	}

	var whom storage.User
	err = u.db.
		Where("username = ?", whomname).
		First(&whom).Error
	u.log.Debug("Read the user to be followed from the database")
	if err != nil {
		u.log.Warn("Unable to find the user that should be followed with error: ", err.Error())
		return err
	}

	user.Follows = append(user.Follows, &whom)
	u.db.Save(&user)
	return nil
}

func (u *User) Unfollow(username string, whomname string) error {
	u.log.Trace("A user is unfollowing another user")

	user, err := u.ReadUserByUsername(username)
	if err != nil {
		u.log.Warn("Unable to find the user wanting to unfollow another user with error: ", err.Error())
		return err
	}

	whom, err := u.ReadUserByUsername(whomname)
	if err != nil {
		u.log.Warn("Unable to find the user that should be unfollowed with error: ", err.Error())
		return err
	}

	u.db.Exec("DELETE FROM follows WHERE user_id = ? AND whom_id = ?", user.ID, whom.ID)
	return nil
}

func (u *User) IsPasswordCorrect(username string, password string) bool {
	u.log.Trace("Checking if a password is correct")

	passwordHashed := u.hash(password)
	usr, err := u.ReadUserByUsername(username)
	if err != nil {
		return false
	}

	return (usr.PwHash == passwordHashed)
}

func (u *User) IsUsernameTaken(username string) bool {
	u.log.Trace("Checking if a username is taken")

	_, err := u.ReadUserByUsername(username)
	return err == nil
}
