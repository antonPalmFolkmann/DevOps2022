package services_test

import (
	"testing"

	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/jinzhu/gorm"
)

var UserService services.IUser = services.NewUserService(&gorm.DB{})

func TestCreateUser(t *testing.T)  {
	// Arrange

	// Act
	UserService.CreateUser("Username", "Email", "Password")

	// Assert
}

