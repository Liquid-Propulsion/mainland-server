package types

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           *uint
	Name         string
	Username     string
	PasswordHash []byte
	TOTPSecret   string
}
