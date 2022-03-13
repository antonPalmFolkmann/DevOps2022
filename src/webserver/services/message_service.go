package services

import (
	"time"

	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/jinzhu/gorm"
)

type IMessage interface {
	CreateMessage(UserID uint, text string, flagged bool) error
	CreateUnflaggedMessage(UserID uint, text string) error
	ReadAllMessages() ([]storage.MessageDTO, error)
	ReadAllMessagesByAuthorId(ID uint) ([]storage.MessageDTO, error)
	ReadAllFlaggedMessages() ([]storage.MessageDTO, error)
	ReadAllFlaggedMessagesByAuthorId(ID uint) ([]storage.MessageDTO, error)
	ReadMessageById(ID uint) (storage.MessageDTO, error)
	UpdateMessage(ID uint, UserID uint, text string, flagged bool) error
	DeleteMessage(ID uint) error
}

type Message struct {
	db *gorm.DB
}

func NewMessageService(db *gorm.DB) *Message {
	return &Message{db: db}
}

func (m *Message) CreateMessage(UserID uint, text string, flagged bool) error {
	message := storage.Message{UserID: UserID, Text: text, PubDate: int(time.Now().Unix()), Flagged: flagged}
	err := m.db.Create(message).Error
	return err
}

func (m *Message) CreateUnflaggedMessage(UserID uint, text string) error {
	err := m.CreateMessage(UserID, text, false)
	return err
}

func (m *Message) ReadAllMessages() ([]storage.MessageDTO, error) {
	var messages = make([]storage.MessageDTO, 0)
	err := m.db.Select("user_id", "text", "pub_date", "flagged").
				Where("flagged = 0").
				Find(&messages).Error
	return messages, err
}

func (m *Message) ReadAllMessagesByAuthorId(ID uint) ([]storage.MessageDTO, error) {
	var messages = make([]storage.MessageDTO, 0)
	err := m.db.Select("user_id", "text", "pub_date", "flagged").
				Where("user_id = ?", ID).
				Find(&messages).Error
	return messages, err
}

func (m *Message) ReadAllFlaggedMessages() ([]storage.MessageDTO, error) {
	var messages = make([]storage.MessageDTO, 0)
	err := m.db.Select("user_id", "text", "pub_date", "flagged").
				Where("flagged = 1").
				Find(&messages).Error
	return messages, err
}

func (m *Message) ReadAllFlaggedMessagesByAuthorId(ID uint) ([]storage.MessageDTO, error) {
	var messages = make([]storage.MessageDTO, 0)
	err := m.db.Select("user_id", "text", "pub_date", "flagged").
				Where("user_id = ? AND flagged = 1", ID).
				Find(&messages).Error
	return messages, err
}

func (m *Message) ReadMessageById(ID uint) (storage.MessageDTO, error) {
	var message storage.MessageDTO
	err := m.db.Select("user_id", "text", "pub_date", "flagged").
				Where("message_id = ?", ID).
				Find(&message).Error
	return message, err
}

func (m *Message) UpdateMessage(ID uint, UserId uint, text string, flagged bool) error {
	var message storage.Message
	err := m.db.Where("user_id = ?", ID).
				Find(&message).Error
	if err != nil {
		return err
	}

	message.UserID = UserId
	message.Text = text
	message.PubDate = int(time.Now().Unix())
	message.Flagged = flagged
	
	err = m.db.Save(&message).Error
	return err
}

func (m *Message) DeleteMessage(ID uint) error {
	var message storage.Message
	err := m.db.Delete(&message, ID).Error
	return err
}
