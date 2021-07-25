package message

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ChannelID int64 `gorm:"index"`
	UserID    string
	Text      string
}
