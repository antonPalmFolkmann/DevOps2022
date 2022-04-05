package services_test

import (
	"testing"

	"github.com/antonPalmFolkmann/DevOps2022/services"
	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func setUpMessageTestDB() (*gorm.DB, services.IMessage) {
	log := logrus.New()
	db, _ := gorm.Open("sqlite3", ":memory:")
	storage.Migrate(db)

	userService := services.NewUserService(db, log)
	err := userService.CreateUser("jalle", "jalle@jalle.jalle", "allej")
	if err != nil {
		log.Fatalf("Failed to create user during DB setup for testing: %s", err)
	}
	err = userService.CreateUser("yolo", "yolo@yolo.yolo", "oloy")
	if err != nil {
		log.Fatalf("Failed to create user during DB setup for testing: %s", err)
	}
	err = userService.CreateUser("chrisser", "chrisser@chrisser.chrisser", "swak420")
	if err != nil {
		log.Fatalf("Failed to create user during DB setup for testing: %s", err)
	}

	err = userService.Follow("chrisser", "jalle")
	if err != nil {
		log.Fatalf("Failed to follow user during DB setup for testing: %s", err)
	}

	messageService := services.NewMessageService(db, log)
	err = messageService.CreateMessage("jalle", "en hel masse ting")
	if err != nil {
		log.Fatalf("Failed to create message during DB setup for testing: %s", err)
	}
	err = messageService.CreateMessage("yolo", "skriver også en hel masse ting")
	if err != nil {
		log.Fatalf("Failed to create message during DB setup for testing: %s", err)
	}
	err = messageService.CreateMessage("chrisser", "niet")
	if err != nil {
		log.Fatalf("Failed to create message during DB setup for testing: %s", err)
	}

	return db, messageService
}

func TestCreateMessage(t *testing.T) {
	_, service := setUpMessageTestDB()

	msgs, _ := service.ReadAllMessages(10, 10)
	assert.Len(t, msgs, 3)

	err := service.CreateMessage("jalle", "new message")
	check_if_test_fail(err)
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
