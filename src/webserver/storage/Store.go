package storage

import "time"

type Store interface {
	ReadUserById(id int64) *User
	ReadUserByUsername(username string) *User
	CreateUser(username string, email string, pw_hash string) error
	GetUserIdByUsername(username string) int64

	CreateFollower(whoId int64, whomId int64)
	DeleteFollower(whoId int64, whomId int64)

	CreateMessage(authorId int64, text string, pubDate *time.Time, flagged bool) error
	ReadMessages() []*Message
	ReadMessagesByFlagged(isFlagged bool) []*Message
	ReadMessagesByUserId(id int64) []*Message
}