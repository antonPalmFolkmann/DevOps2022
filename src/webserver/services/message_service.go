package services

import (
	"time"

	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type IMessage interface {
	CreateMessage(username string, text string) error
	ReadAllMessages(limit int, offset int) ([]storage.MessageDTO, error)
	ReadAllMessagesByUsername(username string) ([]storage.MessageDTO, error)
	ReadAllMessagesOfFollowedUsers(username string) ([]storage.MessageDTO, error)
	FlagMessage(ID uint) error
}

type Message struct {
	db  *gorm.DB
	log *logrus.Logger
}

func NewMessageService(db *gorm.DB, log *logrus.Logger) *Message {
	return &Message{db: db, log: log}
}

func (m *Message) CreateMessage(username string, text string) error {
	m.log.Trace("Creating a message")

	userID := m.getUserIDFromUsername(username)
	message := storage.Message{UserID: userID, Text: text, PubDate: time.Now().Unix(), Flagged: false}
	err := m.db.Create(&message).Error
	m.log.Debug("Created a message on the database")
	return err
}

func (m *Message) ReadAllMessages(limit int, offset int) ([]storage.MessageDTO, error) {
	m.log.Trace("Reading all messages")

	var messages = make([]storage.Message, 0)
	var messageDTOs = make([]storage.MessageDTO, 0)
	err := m.db.Offset(offset).Limit(limit).Where("flagged = 0").Find(&messages).Error
	m.log.Debug("Read messages from database")

	for _, v := range messages {
		username := m.getUsernameFromUserID(v.UserID)
		messageDTO := storage.MessageDTO{UserID: v.UserID, Username: username, Text: v.Text, PubDate: time.Unix(v.PubDate, 0), Flagged: v.Flagged}
		messageDTOs = append(messageDTOs, messageDTO)
	}
	return messageDTOs, err
}

func (m *Message) ReadAllMessagesByUsername(username string) ([]storage.MessageDTO, error) {
	m.log.Trace("Reading all messages by username")

	var messages = make([]storage.Message, 0)
	ID := m.getUserIDFromUsername(username)
	err := m.db.Where("user_id = ?", ID).Find(&messages).Error
	m.log.Debug("Read all messages by username on the database")

	var messageDTOs = make([]storage.MessageDTO, 0)
	for _, v := range messages {
		messageDTO := storage.MessageDTO{UserID: v.UserID, Username: username, Text: v.Text, PubDate: time.Unix(v.PubDate, 0), Flagged: v.Flagged}
		messageDTOs = append(messageDTOs, messageDTO)
	}
	return messageDTOs, err
}

func (m *Message) ReadAllMessagesOfFollowedUsers(username string) ([]storage.MessageDTO, error) {
	m.log.Trace("Reading all messages of followed users")

	var user storage.User
	err := m.db.Preload("Follows").Where("username = ?", username).Find(&user).Error
	if err != nil {
		return nil, err
	}
	m.log.Debug("Read all followers for the user")

	var messages = make([]storage.Message, 0)
	err = m.db.Where("user_id IN (?)", userToIDs(user.Follows)).Find(&messages).Error
	if err != nil {
		m.log.Warn("Could not read messages of followed users with error: ", err.Error())
		return nil, err
	}

	var messageDTOs = make([]storage.MessageDTO, 0)
	for _, v := range messages {
		author := m.getUsernameFromUserID(v.UserID)
		messageDTO := storage.MessageDTO{UserID: v.UserID, Username: author, Text: v.Text, PubDate: time.Unix(v.PubDate, 0), Flagged: v.Flagged}
		messageDTOs = append(messageDTOs, messageDTO)
	}
	return messageDTOs, err
}

func userToIDs(usrs []*storage.User) []uint {
	IDs := make([]uint, 0)
	for _, usr := range usrs {
		IDs = append(IDs, usr.ID)
	}
	return IDs
}

func (m *Message) FlagMessage(ID uint) error {
	m.log.Trace("Flagging a message")
	err := m.db.Raw("UPDATE messages SET flagged=true WHERE messages.id = ?", ID).Error
	m.log.Debug("Flagged the message on the database")
	return err
}

func (m *Message) getUserIDFromUsername(username string) uint {
	m.log.Trace("Getting user ID from username")

	var user storage.User
	m.db.Unscoped().
		Where("username = ?", username).
		Select("id").
		Find(&user)

	m.log.Debug("Read user from the database")
	return user.ID
}

func (m *Message) getUsernameFromUserID(userID uint) string {
	m.log.Trace("Getting username from user ID")

	var user storage.User
	m.db.First(&user, userID)
	m.log.Debug("Read user from the database")
	return user.Username
}
