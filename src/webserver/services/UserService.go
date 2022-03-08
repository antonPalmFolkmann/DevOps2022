package services

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/antonPalmFolkmann/DevOps2022/models"
)

type IUserService interface {
	ReadUserById(id int64) (*User, error) // <----- her
	ReadUserByUsername(username string) *User
	CreateUser(username string, email string, pw_hash string) error
	GetUserIdByUsername(username string) int64
	DeleteUser(id int64) error
}

type UserService struct {}

// Vent det her føles forkert. En user burde jo ikke kunne læse andre usere?
// Vi har nu skrevet at user-struct'en implementerer IUser, hvor vi har defineret
// at man kan læse, oprette og slette brugere. Det burde være 'user_servicen' der kan det
// Yeah så burde vi omdøbe interfacet og structen.
// Det her ligner mere BDSA
// Hmm yeah, skal lige tænke to sek 
// Record til DTO
// Hvad gør vi så med vores struct nu? Læser vi ind fra User til UserService? Det gjorde vi i BDSA
func (u UserService) ReadUserById(id int64) (models.User, error) {
	var user = models.GetUser()
	err := dbconn.Where("user_id = ?", id).Find(&user).Error
	return user, err
}

func GetAllUsers() ([]models.User, error) {
	var users = models.GetUsers()
	err := dbconn.Find(&users).Error
	return users, err
}

func GetUserByID(id int) (models.User, error) {
	var user = models.GetUser()
	err := dbconn.Where("user_id = ?", id).Find(&user).Error
	return user, err
}

func GetUserByUsername(username string) (models.User, error) {
	var user = models.GetUser()
	err := dbconn.Where("username = ?", username).Find(&user).Error
	return user, err
}

func CreateUser(r *http.Request) error {
	var user = models.GetUser()
	_ = json.NewDecoder(r.Body).Decode(&user)
	log.Println(user)

	err := dbconn.Create(&user).Error
	return err
}

func UpdateUser(r *http.Request, id int) error {
	var user = models.GetUser()
	_ = json.NewDecoder(r.Body).Decode(&user)

	err := dbconn.Model(&user).Where("user_id = ?", id).Update(&user).Error
	return err
}

func DeleteUser(id int) error {
	var user = models.GetUser()
	err := dbconn.Delete(&user, id).Error
	return err
}