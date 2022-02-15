package types

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string
	Username     string `gorm:"unique;not null"`
	PasswordHash []byte
	TOTPEnabled  bool
	TOTPSecret   string
}

func (node *User) IsNode() {}
