package message

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	UserID string
	Text   string
}
