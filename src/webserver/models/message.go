package models

import (
	"github.com/jinzhu/gorm"
)

type Message struct {
	gorm.Model
	Message_id int64  `gorm:"type:bigint;" json:"message_id"`
	Author_id  int64  `gorm:"type:bigint;" json:"author_id"`
	Text       string `gorm:"type:text;" json:"text"`
	Pub_date   int64  `gorm:"type:bigint;" json:"pub_date"`
	Flagged    int    `gorm:"type:int;" json:"flagged"`
}
