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
	ReadAllMessagesOfFollowedUsers(username string) ([]storage.MessageDTO ,error)
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
	err := m.db.Create(message).Error
	return err
}

func (m *Message) ReadAllMessages(limit int, offset int) ([]storage.MessageDTO, error) {
	var messages = make([]storage.Message, 0)
	var messageDTOs = make([]storage.MessageDTO, 0)
	err := m.db.Select("user_id", "text", "pub_date", "flagged").
				Where("flagged = 0").
				Find(&messages).Error
				
	for _, v := range messages {
		username := m.getUsernameFromUserID(v.UserID)
		messageDTO := storage.MessageDTO{UserID: v.UserID, Username: username, Text: v.Text, PubDate: time.Unix(v.PubDate, 0), Flagged: v.Flagged}
		messageDTOs = append(messageDTOs, messageDTO)
	}
	return messageDTOs, err
}

func (m *Message) ReadAllMessagesByUsername(username string) ([]storage.MessageDTO, error) {
	var messages = make([]storage.Message, 0)
	var messageDTOs = make([]storage.MessageDTO, 0)
	ID := m.getUserIDFromUsername(username)
	err := m.db.Select("user_id", "text", "pub_date", "flagged").
				Where("user_id = ?", ID).
				Find(&messages).Error
				
	for _, v := range messages {
		messageDTO := storage.MessageDTO{UserID: v.UserID, Username: username, Text: v.Text, PubDate: time.Unix(v.PubDate, 0), Flagged: v.Flagged}
		messageDTOs = append(messageDTOs, messageDTO)
	}
	return messageDTOs, err
}

func (m *Message) ReadAllMessagesOfFollowedUsers(username string) ([]storage.MessageDTO ,error) {
	var messages = make([]storage.Message, 0)
	var messageDTOs = make([]storage.MessageDTO, 0)
	ID := m.getUserIDFromUsername(username)
	err := m.db.Select("user_id", "text", "pub_date", "flagged").
				Where("message_id = ?", ID).
				Find(&messages).Error
	
	for _, v := range messages {
		messageDTO := storage.MessageDTO{UserID: v.UserID, Username: username, Text: v.Text, PubDate: time.Unix(v.PubDate, 0), Flagged: v.Flagged}
		messageDTOs = append(messageDTOs, messageDTO)
	}
	return messageDTOs, err
}

func (m *Message) FlagMessage(ID uint) error {
	var message storage.Message
	return m.db.Model(&message).
				Where("message_id = ?", ID).
				Update("flagged", 1).Error
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
