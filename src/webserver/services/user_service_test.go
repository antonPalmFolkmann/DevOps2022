package services_test

import (
	"crypto/md5"
	"fmt"
	"io"
	"testing"

	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/antonPalmFolkmann/DevOps2022/storage"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setUp() (*gorm.DB, services.IUser) {
	db, _ := gorm.Open("sqlite3", ":memory:")
	storage.Migrate(db)

	userService := services.NewUserService(db)
	userService.CreateUser("jalle", "jalle@jalle.jalle", "allej")
	userService.CreateUser("yolo", "yolo@yolo.yolo", "oloy")
	userService.CreateUser("chrisser", "chrisser@chrisser.chrisser", "swak420")

	return db, userService
}

// ------------------- TESTS -------------------------


func Test_CreateUser(t *testing.T) {
	// Arrange
	db, service := setUp()
	var (
		actual storage.User
		username = "user"
		email = "user@itu.dk"
		password = "******"
		passwordHashed string
	)

	hash := md5.New()
	io.WriteString(hash, password)
	passwordHashed = fmt.Sprintf("%x", hash.Sum(nil))

	// Act
	service.CreateUser(username, email, password)

	db.Where("username = ?", username).First(&actual)

	// Assert
	assert.Equal(t, username, actual.Username)
	assert.Equal(t, email, actual.Email)
	assert.Equal(t, passwordHashed, actual.PwHash)
}


func Test_ReadAllUsers(t *testing.T) {
	// Arrange
	_, service := setUp()
	var expected = 3

	// Act
	actual, _ := service.ReadAllUsers()

	// Assert
	assert.Equal(t, expected, len(actual))
}

func Test_ReadUserIdByUsername_Found(t *testing.T) {
	// Arrange
	_, service := setUp()
	var (
		username = "jalle"
		expected = uint(1)
	)

	// Act
	actual, _ := service.ReadUserIdByUsername(username)
	

	// Assert
	assert.Equal(t, expected, actual)
}

func Test_ReadUserIdByUsername_Error(t *testing.T) {
	// Arrange
	_, service := setUp()
	var username = "notfound"

	// Act
	_, actual := service.ReadUserIdByUsername(username)

	// Assert
	assert.Error(t, actual)
}

func Test_follow(t *testing.T) {
	db, service := setUp()

	service.Follow("jalle", "yolo")

	var user storage.User
	db.Preload("Follows").Where("username = ?", "jalle").First(&user)
	assert.Equal(t, 1, len(user.Follows))
}

func Test_unfollow(t *testing.T) {
	db, service := setUp()

	service.Follow("jalle", "yolo")
	service.Unfollow("jalle", "yolo")

	var user storage.User
	db.Preload("Follows").Where("username = ?", "jalle").First(&user)
	assert.Equal(t, 0, len(user.Follows))
}

func Test_followFollowed(t *testing.T) {
	db, service := setUp()

	service.Follow("jalle", "yolo")
	service.Follow("jalle", "yolo")
	var user storage.User

	db.Preload("Follows").Where("username = ?", "jalle").First(&user)
	assert.Len(t, user.Follows, 1)
}

func Test_unfollowNotFollowed(t *testing.T) {
	db, service := setUp()

	service.Follow("jalle", "yolo")
	service.Unfollow("jalle", "chrisser")

	var user storage.User
	db.Preload("Follows").Where("username = ?", "jalle").First(&user)
	assert.Len(t, user.Follows, 1)
}

func Test_followNonExistentReturnsError(t *testing.T) {
	_, service := setUp()
	err := service.Follow("jalle", "RNSK RNSK RNSK RNSK RNSK RNSK RNSK RNSK RNSK RNSK RNSK RSNK RSNK RNSK RNSK RNSK")
	assert.NotNil(t, err)
}

func Test_unfollowNonExistentReturnsError(t *testing.T) {
	_, service := setUp()

	service.CreateUser("jalle", "jalle@jalle.jalle", "allej")
	err := service.Unfollow("jalle", "Benjamin, The Destroyer Of Worlds and Harbringer Of Death")
	assert.NotNil(t, err)
}

func Test_IsPasswordCorrect_True(t *testing.T)  {
	// Arrange
	_, service := setUp()
	var username = "jalle"
	var password = "allej"

	// Act
	actual := service.IsPasswordCorrect(username, password)

	// Assert
	assert.True(t, actual)
}

func Test_IsPasswordCorrect_False(t *testing.T)  {
	// Arrange
	_, service := setUp()
	var username = "jalle"
	var password = "sdf"

	// Act
	actual := service.IsPasswordCorrect(username, password)

	// Assert
	assert.False(t, actual)
}

func Test_IsUsernameTaken_True(t *testing.T)  {
	_, service := setUp()

	actual := service.IsUsernameTaken("jaææææ")

	assert.True(t, actual)
}

func Test_IsUsernameTaken_False(t *testing.T)  {
	_, service := setUp()

	actual := service.IsUsernameTaken("jalle")

	assert.False(t, actual)
}
