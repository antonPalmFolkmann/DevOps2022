package services

import (
	"errors"

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
	// message := storage.Message{UserID: userID, Text: text, PubDate: time.Now(), Flagged: false}
	// err := m.db.Create(message).Error
	return errors.New("not implemented yet")
}

func (m *Message) ReadAllMessages(limit int, offset int) ([]storage.MessageDTO, error) {
	// var messages = make([]storage.Message, 0)
	// err := m.db.Select("user_id", "text", "pub_date", "flagged").
	// 	Where("flagged = 0").
	// 	Find(&messages).Error
	return make([]storage.MessageDTO, 0), errors.New("not implemented yet")
}

func (m *Message) ReadAllMessagesByUsername(username string) ([]storage.MessageDTO, error) {
	// var messages = make([]storage.Message, 0)
	// err := m.db.Select("user_id", "text", "pub_date", "flagged").
	// 	Where("user_id = ?", id).
	// 	Find(&messages).Error
	return make([]storage.MessageDTO, 0), errors.New("not implemented yet")
}

func (m *Message) ReadAllMessagesOfFollowedUsers(username string) ([]storage.MessageDTO, error) {
	// var message storage.Message
	// err := m.db.Select("user_id", "text", "pub_date", "flagged").
	// 	Where("message_id = ?", id).
	// 	Find(&message).Error
	return make([]storage.MessageDTO, 0), errors.New("not implemented yet")
}

func (m *Message) FlagMessage(ID uint) error {
	return errors.New("not implemented yet")
}
