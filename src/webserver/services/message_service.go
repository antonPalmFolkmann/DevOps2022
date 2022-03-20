package services

import (
	"time"

	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/jinzhu/gorm"
)

type IMessage interface {
	CreateMessage(username string, text string) error
	ReadAllMessages(limit int, offset int) ([]storage.MessageDTO, error)
	ReadAllMessagesByUsername(username string) ([]storage.MessageDTO, error)
	ReadAllMessagesOfFollowedUsers(username string) ([]storage.MessageDTO, error)
	FlagMessage(ID uint) error
}

type Message struct {
	db *gorm.DB
}

func NewMessageService(db *gorm.DB) *Message {
	return &Message{db: db}
}

func (m *Message) CreateMessage(username string, text string) error {
	userID := m.getUserIDFromUsername(username)
	message := storage.Message{UserID: userID, Text: text, PubDate: time.Now().Unix(), Flagged: false}
	err := m.db.Create(&message).Error
	return err
}

func (m *Message) ReadAllMessages(limit int, offset int) ([]storage.MessageDTO, error) {
	var messages = make([]storage.Message, 0)
	var messageDTOs = make([]storage.MessageDTO, 0)
	err := m.db.Where("flagged = 0").Find(&messages).Error

	for _, v := range messages {
		username := m.getUsernameFromUserID(v.UserID)
		messageDTO := storage.MessageDTO{UserID: v.UserID, Username: username, Text: v.Text, PubDate: time.Unix(v.PubDate, 0), Flagged: v.Flagged}
		messageDTOs = append(messageDTOs, messageDTO)
	}
	return messageDTOs, err
}

func (m *Message) ReadAllMessagesByUsername(username string) ([]storage.MessageDTO, error) {
	var messages = make([]storage.Message, 0)
	ID := m.getUserIDFromUsername(username)
	err := m.db.Where("user_id = ?", ID).Find(&messages).Error

	var messageDTOs = make([]storage.MessageDTO, 0)
	for _, v := range messages {
		messageDTO := storage.MessageDTO{UserID: v.UserID, Username: username, Text: v.Text, PubDate: time.Unix(v.PubDate, 0), Flagged: v.Flagged}
		messageDTOs = append(messageDTOs, messageDTO)
	}
	return messageDTOs, err
}

func (m *Message) ReadAllMessagesOfFollowedUsers(username string) ([]storage.MessageDTO, error) {
	var user storage.User
	err := m.db.Preload("Follows").Where("username = ?", username).Find(&user).Error
	if err != nil {
		return nil, err
	}

	var messages = make([]storage.Message, 0)
	err = m.db.Where("user_id IN (?)", userToIDs(user.Follows)).Find(&messages).Error
	if err != nil {
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
	err := m.db.Raw("UPDATE messages SET flagged=true WHERE messages.id = ?", ID).Error
	return err
}

func (m *Message) getUserIDFromUsername(username string) uint {
	var user storage.User
	m.db.Unscoped().
		Where("username = ?", username).
		Select("id").
		Find(&user)
	return user.ID
}

func (m *Message) getUsernameFromUserID(userID uint) string {
	var user storage.User
	m.db.First(&user, userID)
	return user.Username
}
