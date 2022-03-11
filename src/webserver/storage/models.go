package storage

import "github.com/jinzhu/gorm"

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

type UserDTO struct {
	ID   uint
	Username string
	Email    string
	PwHash   string
	Messages []MessageDTO
	Follows  []UserDTO
}

type Message struct {
	// gorm.Model provides ID
	gorm.Model
	// Creates a "message belongs-to one user" relationship
	UserID  uint
	Text    string
	PubDate int
	Flagged bool
}

type MessageDTO struct {
	UserID  uint
	Text    string
	PubDate int
	Flagged bool
}
