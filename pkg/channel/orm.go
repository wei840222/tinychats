package channel

import (
	"gorm.io/gorm"
)

type Channel struct {
	gorm.Model
	Name string
}
