package services

import (
	"github.com/antonPalmFolkmann/DevOps2022/models"
)

type IMessageService interface {
	CreateMessage(message models.Message) error
	ReadAllMessages() ([]models.Message, error)
	ReadAllMessagesByAuthorId(id int) ([]models.Message, error)
	ReadAllFlaggedMessages() ([]models.Message, error)
	ReadAllFlaggedMessagesByAuthorId(id int) ([]models.Message, error)
	ReadMessageById(id int) (models.Message, error)
	UpdateMessage(message models.Message, id int) error
	DeleteMessage(id int) error
}

type MessageService struct{}

func (m MessageService) CreateMessage(message *models.Message) error {
	err := dbconn.Create(&message).Error
	return err
}

func (m MessageService) ReadAllMessages() ([]models.Message, error) {
	var messages = make([]models.Message, 0)
	err := dbconn.Where("flagged = 0").Find(&messages).Error
	return messages, err
}

func (m MessageService) ReadAllMessagesByAuthorId(id int) ([]models.Message, error) {
	var messages = make([]models.Message, 0)
	err := dbconn.Where("author_id = ?", id).Find(&messages).Error
	return messages, err
}

func (m MessageService) ReadAllFlaggedMessages() ([]models.Message, error) {
	var messages = make([]models.Message, 0)
	err := dbconn.Where("flagged = 1").Find(&messages).Error
	return messages, err
}

func (m MessageService) ReadAllFlaggedMessagesByAuthorId(id int) ([]models.Message, error) {
	var messages = make([]models.Message, 0)
	err := dbconn.Where("author_id = ? AND flagged = 1", id).Find(&messages).Error
	return messages, err
}

func (m MessageService) ReadMessageById(id int) (models.Message, error) {
	var message models.Message
	err := dbconn.Where("message_id = ?", id).Find(&message).Error
	return message, err
}

func (m MessageService) UpdateMessage(message *models.Message, id int) error {
	err := dbconn.Model(&message).Where("message_id = ?", id).Update(&message).Error
	return err
}

func (m MessageService) DeleteMessage(id int) error {
	var message models.Message
	err := dbconn.Delete(&message, id).Error
	return err
}
