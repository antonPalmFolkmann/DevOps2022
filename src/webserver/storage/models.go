package storage

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	UserId   int64  `gorm:"type:bigint;" json:"user_id"`
	Username string `gorm:"type:varchar(255);" json:"username"`
	Email    string `gorm:"type:varchar(255);" json:"email"`
	PwHash   string `gorm:"type:varchar(255);" json:"pw_hash"`
}

type Message struct {
	gorm.Model
	MessageId int64  `gorm:"type:bigint;" json:"message_id"`
	AuthorId  int64  `gorm:"type:bigint;" json:"author_id"`
	Text      string `gorm:"type:text;" json:"text"`
	PubDate   int64  `gorm:"type:bigint;" json:"pub_date"`
	Flagged   int    `gorm:"type:int;" json:"flagged"`
}

type Follows struct {
	gorm.Model
	WhoId  int64 `gorm:"primaryKey"`
	WhomId int64 `gorm:"primaryKey"`
}
