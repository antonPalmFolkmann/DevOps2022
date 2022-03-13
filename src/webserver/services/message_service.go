package services

import (
	"time"

	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/jinzhu/gorm"
)

type IMessage interface {
	CreateMessage(userID uint, text string, pubDate time.Time, flagged bool) error
	CreateUnflaggedMessage(userID uint, text string, pubDate time.Time)
	ReadAllMessages() ([]storage.MessageDTO, error)
	ReadAllMessagesByAuthorId(ID uint) ([]storage.MessageDTO, error)
	ReadAllFlaggedMessages() ([]storage.MessageDTO, error)
	ReadAllFlaggedMessagesByAuthorId(ID uint) ([]storage.MessageDTO, error)
	ReadAllMessagesForUsername(username string) ([]storage.MessageDTO, error)
	ReadMessageById(ID uint) (storage.MessageDTO, error)
	ReadAllMessagesForUsername(username string) (storage.MessageDTO, error)
	UpdateMessage(ID uint, userID uint, text string, pubDate time.Time, flagged bool) error
	DeleteMessage(ID uint) error
}

type Message struct {
	db *gorm.DB
}

func NewMessageService(db *gorm.DB) *Message {
	return &Message{db: db}
}

func (m *Message) CreateMessage(message storage.Message) error {
	err := m.db.Create(message).Error
	return err
}

func (m *Message) ReadAllMessages() ([]storage.Message, error) {
	var messages = make([]storage.Message, 0)
	err := m.db.Where("flagged = 0").Find(&messages).Error
	return messages, err
}

func (m *Message) ReadAllMessagesByAuthorId(id uint) ([]storage.Message, error) {
	var messages = make([]storage.Message, 0)
	err := m.db.Where("author_id = ?", id).Find(&messages).Error
	return messages, err
}

func (m *Message) ReadAllFlaggedMessages() ([]storage.Message, error) {
	var messages = make([]storage.Message, 0)
	err := m.db.Where("flagged = 1").Find(&messages).Error
	return messages, err
}

func (m *Message) ReadAllFlaggedMessagesByAuthorId(id uint) ([]storage.Message, error) {
	var messages = make([]storage.Message, 0)
	err := m.db.Where("author_id = ? AND flagged = 1", id).Find(&messages).Error
	return messages, err
}

func (m *Message) ReadMessageById(id uint) (storage.Message, error) {
	var message storage.Message
	err := m.db.Where("message_id = ?", id).Find(&message).Error
	return message, err
}

func (m *Message) UpdateMessage(message storage.Message, id uint) error {
	err := m.db.Model(&message).Where("message_id = ?", id).Update(&message).Error
	return err
}

func (m *Message) DeleteMessage(id uint) error {
	var message storage.Message
	err := m.db.Delete(&message, id).Error
	return err
}
