package todo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	UUID uuid.UUID
	Text string
	Done bool
}
