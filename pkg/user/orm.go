package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID    string `gorm:"uniqueIndex"`
	Name      string
	AvatarURL string
}
