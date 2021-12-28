package types

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string
	Username     string
	PasswordHash []byte
	TOTPSecret   string
}
