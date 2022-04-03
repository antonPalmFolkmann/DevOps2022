package services_test

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
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
	err := userService.CreateUser("jalle", "jalle@jalle.jalle", "allej")
	check_if_test_fail(err)
	err = userService.CreateUser("yolo", "yolo@yolo.yolo", "oloy")
	check_if_test_fail(err)
	err = userService.CreateUser("chrisser", "chrisser@chrisser.chrisser", "swak420")
	check_if_test_fail(err)
	return db, userService
}

// ------------------- TESTS -------------------------

func Test_CreateUser(t *testing.T) {
	// Arrange
	db, service := setUp()
	var (
		actual         storage.User
		username       = "user"
		email          = "user@itu.dk"
		password       = "******"
		passwordHashed string
	)

	hash := md5.New()
	_, err := io.WriteString(hash, password)
	check_if_test_fail(err)

	passwordHashed = fmt.Sprintf("%x", hash.Sum(nil))

	// Act
	err = service.CreateUser(username, email, password)
	check_if_test_fail(err)

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

	err := service.Follow("jalle", "yolo")
	check_if_test_fail(err)

	var user storage.User
	db.Preload("Follows").Where("username = ?", "jalle").First(&user)
	assert.Equal(t, 1, len(user.Follows))
}

func Test_unfollow(t *testing.T) {
	db, service := setUp()

	err := service.Follow("jalle", "yolo")
	check_if_test_fail(err)
	err = service.Unfollow("jalle", "yolo")
	check_if_test_fail(err)

	var user storage.User
	db.Preload("Follows").Where("username = ?", "jalle").First(&user)
	assert.Equal(t, 0, len(user.Follows))
}

func Test_followFollowed(t *testing.T) {
	db, service := setUp()

	err := service.Follow("jalle", "yolo")
	check_if_test_fail(err)
	err = service.Follow("jalle", "yolo")
	check_if_test_fail(err)
	var user storage.User

	db.Preload("Follows").Where("username = ?", "jalle").First(&user)
	assert.Len(t, user.Follows, 1)
}

func Test_unfollowNotFollowed(t *testing.T) {
	db, service := setUp()

	err := service.Follow("jalle", "yolo")
	check_if_test_fail(err)
	err = service.Unfollow("jalle", "chrisser")
	check_if_test_fail(err)

	var user storage.User
	db.Preload("Follows").Where("username = ?", "jalle").First(&user)
	assert.Len(t, user.Follows, 1)
}

func Test_followNonExistentReturnsError(t *testing.T) {
	_, service := setUp()
	err := service.Follow("jalle", "RNSK RNSK RNSK RNSK RNSK RNSK RNSK RNSK RNSK RNSK RNSK RSNK RSNK RNSK RNSK RNSK")
	check_if_test_fail(err)
	assert.NotNil(t, err)
}

func Test_unfollowNonExistentReturnsError(t *testing.T) {
	_, service := setUp()

	err := service.CreateUser("jalle", "jalle@jalle.jalle", "allej")
	check_if_test_fail(err)
	err = service.Unfollow("jalle", "Benjamin, The Destroyer Of Worlds and Harbringer Of Death")
	assert.NotNil(t, err)
}

func Test_IsPasswordCorrect_True(t *testing.T) {
	// Arrange
	_, service := setUp()
	var username = "jalle"
	var password = "allej"

	// Act
	actual := service.IsPasswordCorrect(username, password)

	// Assert
	assert.True(t, actual)
}

func Test_IsPasswordCorrect_False(t *testing.T) {
	// Arrange
	_, service := setUp()
	var username = "jalle"
	var password = "sdf"

	// Act
	actual := service.IsPasswordCorrect(username, password)

	// Assert
	assert.False(t, actual)
}

func Test_IsUsernameTaken_False(t *testing.T) {
	_, service := setUp()

	actual := service.IsUsernameTaken("jaææææ")

	assert.False(t, actual)
}

func Test_IsUsernameTaken_True(t *testing.T) {
	_, service := setUp()

	actual := service.IsUsernameTaken("jalle")

	assert.True(t, actual)
}

func check_if_test_fail(err error) {
	if err != nil {
		log.Fatalf("An error occured during test: %s", err.Error())
	}
}
