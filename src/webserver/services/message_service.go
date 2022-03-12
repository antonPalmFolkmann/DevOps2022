package services

import (
	"time"

	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/jinzhu/gorm"
)

type IMessage interface {
	CreateMessage(authorID uint, text string, flagged bool) error
	CreateUnflaggedMessage(authorID uint, text string) error
	ReadAllMessages() ([]storage.MessageDTO, error)
	ReadAllMessagesByAuthorId(ID uint) ([]storage.MessageDTO, error)
	ReadAllFlaggedMessages() ([]storage.MessageDTO, error)
	ReadAllFlaggedMessagesByAuthorId(ID uint) ([]storage.MessageDTO, error)
	ReadMessageById(ID uint) (storage.MessageDTO, error)
	UpdateMessage(ID uint, authorID uint, text string, flagged bool) error
	DeleteMessage(ID uint) error
}

type Message struct {
	db *gorm.DB
}

func NewMessageService(db *gorm.DB) *Message {
	return &Message{db: db}
}

func (m *Message) CreateMessage(authorID uint, text string, flagged bool) error {
	message := storage.Message{AuthorID: authorID, Text: text, PubDate: time.Now(), Flagged: flagged}
	err := m.db.Create(message).Error
	return err
}

func (m *Message) CreateUnflaggedMessage(authorID uint, text string) error {
	err := m.CreateMessage(authorID, text, false)
	return err
}

func (m *Message) ReadAllMessages() ([]storage.MessageDTO, error) {
	var messages = make([]storage.MessageDTO, 0)
	err := m.db.Select("author_id", "text", "pub_date", "flagged").
				Where("flagged = 0").
				Find(&messages).Error
	return messages, err
}

func (m *Message) ReadAllMessagesByAuthorId(id uint) ([]storage.MessageDTO, error) {
	var messages = make([]storage.MessageDTO, 0)
	err := m.db.Select("author_id", "text", "pub_date", "flagged").
				Where("author_id = ?", id).
				Find(&messages).Error
	return messages, err
}

func (m *Message) ReadAllFlaggedMessages() ([]storage.MessageDTO, error) {
	var messages = make([]storage.MessageDTO, 0)
	err := m.db.Select("author_id", "text", "pub_date", "flagged").
				Where("flagged = 1").
				Find(&messages).Error
	return messages, err
}

func (m *Message) ReadAllFlaggedMessagesByAuthorId(id uint) ([]storage.MessageDTO, error) {
	var messages = make([]storage.MessageDTO, 0)
	err := m.db.Select("author_id", "text", "pub_date", "flagged").
				Where("author_id = ? AND flagged = 1", id).
				Find(&messages).Error
	return messages, err
}

func (m *Message) ReadMessageById(id uint) (storage.MessageDTO, error) {
	var message storage.MessageDTO
	err := m.db.Select("author_id", "text", "pub_date", "flagged").
				Where("message_id = ?", id).
				Find(&message).Error
	return message, err
}

func (m *Message) UpdateMessage(ID uint, authorID uint, text string, flagged bool) error {
	var message storage.Message
	err := m.db.Where("author_id = ?", ID).
				Find(&message).Error
	if err != nil {
		return err
	}

	message.AuthorID = authorID
	message.Text = text
	message.PubDate = time.Now()
	message.Flagged = flagged
	
	err = m.db.Save(&message).Error
	return err
}

func (m *Message) DeleteMessage(id uint) error {
	var message storage.Message
	err := m.db.Delete(&message, id).Error
	return err
}
