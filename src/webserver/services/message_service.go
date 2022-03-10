package services

import (
	"github.com/antonPalmFolkmann/DevOps2022/storage"
	"github.com/jinzhu/gorm"
)

type IMessageService interface {
	CreateMessage(message storage.Message) error
	ReadAllMessages() ([]storage.Message, error)
	ReadAllMessagesByAuthorId(id uint) ([]storage.Message, error)
	ReadAllFlaggedMessages() ([]storage.Message, error)
	ReadAllFlaggedMessagesByAuthorId(id uint) ([]storage.Message, error)
	ReadMessageById(id int) (storage.Message, error)
	UpdateMessage(message storage.Message, id uint) error
	DeleteMessage(id uint) error
}

type MessageService struct {
	db *gorm.DB
}

func NewMessageService(db *gorm.DB) *MessageService {
	return &MessageService{db: db}
}

func (m *MessageService) CreateMessage(message storage.Message) error {
	err := m.db.Create(message).Error
	return err
}

func (m *MessageService) ReadAllMessages() ([]storage.Message, error) {
	var messages = make([]storage.Message, 0)
	err := m.db.Where("flagged = 0").Find(&messages).Error
	return messages, err
}

func (m *MessageService) ReadAllMessagesByAuthorId(id uint) ([]storage.Message, error) {
	var messages = make([]storage.Message, 0)
	err := m.db.Where("author_id = ?", id).Find(&messages).Error
	return messages, err
}

func (m *MessageService) ReadAllFlaggedMessages() ([]storage.Message, error) {
	var messages = make([]storage.Message, 0)
	err := m.db.Where("flagged = 1").Find(&messages).Error
	return messages, err
}

func (m *MessageService) ReadAllFlaggedMessagesByAuthorId(id uint) ([]storage.Message, error) {
	var messages = make([]storage.Message, 0)
	err := m.db.Where("author_id = ? AND flagged = 1", id).Find(&messages).Error
	return messages, err
}

func (m *MessageService) ReadMessageById(id uint) (storage.Message, error) {
	var message storage.Message
	err := m.db.Where("message_id = ?", id).Find(&message).Error
	return message, err
}

func (m *MessageService) UpdateMessage(message storage.Message, id uint) error {
	err := m.db.Model(&message).Where("message_id = ?", id).Update(&message).Error
	return err
}

func (m *MessageService) DeleteMessage(id uint) error {
	var message storage.Message
	err := m.db.Delete(&message, id).Error
	return err
}
