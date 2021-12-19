package types

import (
	"gorm.io/gorm"
)

type SOP struct {
	gorm.Model
	ID          *uint
	Name        string
	Description string
	Code        string
}
