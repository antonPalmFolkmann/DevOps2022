package models

import (
	"github.com/jinzhu/gorm"
)

type Follows struct {
	gorm.Model
	Who_id	int64 `gorm:"primaryKey"`
	Whom_id	int64 `gorm:"primaryKey"`
}

func GetFollows() Follows{
	var follows Follows
	return follows
}

func GetAllFollows() []Follows{
	var followers []Follows
	return followers
}