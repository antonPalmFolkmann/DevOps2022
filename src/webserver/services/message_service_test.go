package services_test

import (
	"testing"

	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func setUpMessageTestDB() (*gorm.DB, services.IMessage) {
	db, _ := gorm.Open("sqlite3", ":memory:")
	storage.Migrate(db)

	userService := services.NewUserService(db)
	userService.CreateUser("jalle", "jalle@jalle.jalle", "allej")
	userService.CreateUser("yolo", "yolo@yolo.yolo", "oloy")
	userService.CreateUser("chrisser", "chrisser@chrisser.chrisser", "swak420")

	userService.Follow("chrisser", "jalle")

	messageService := services.NewMessageService(db)
	messageService.CreateMessage("jalle", "en hel masse ting")
	messageService.CreateMessage("yolo", "skriver også en hel masse ting")
	messageService.CreateMessage("chrisser", "niet")

	return db, messageService
}

func TestCreateMessage(t *testing.T) {
	_, service := setUpMessageTestDB()

	msgs, _ := service.ReadAllMessages(10, 10)
	assert.Len(t, msgs, 3)

	service.CreateMessage("jalle", "new message")
	actual, _ := service.ReadAllMessages(10, 10)
	assert.Len(t, actual, 4)
}

func TestReadAllMessagesReturns3(t *testing.T) {
	_, service := setUpMessageTestDB()

	actual, _ := service.ReadAllMessages(10, 10)
	assert.Len(t, actual, 3)
}

func TestReadAllMessagesByUsername(t *testing.T) {
	_, service := setUpMessageTestDB()

	actual, _ := service.ReadAllMessagesByUsername("chrisser")
	assert.Len(t, actual, 1)
	assert.Equal(t, actual[0].Username, "chrisser")
}

func TestReadAllMessagesOfFollowedUsers(t *testing.T) {
	_, service := setUpMessageTestDB()

	actual, _ := service.ReadAllMessagesOfFollowedUsers("chrisser")
	assert.Len(t, actual, 1)
	assert.Equal(t, "jalle", actual[0].Username)
}
