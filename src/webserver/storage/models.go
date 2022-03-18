package storage

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	// gorm.Model provides ID
	gorm.Model
	Username string
	Email    string
	PwHash   string
	// Creates a "user has many messages" relationship
	Messages []Message
	// Creates a "user follows many users" relationship
	Follows []*User `gorm:"many2many:follows;association_jointable_foreignkey:whom_id"`
}

type Message struct {
	// gorm.Model provides ID
	gorm.Model
	// Creates a "message belongs-to one user" relationship
	UserID  uint
	Text    string
	PubDate time.Time
	Flagged bool
}
