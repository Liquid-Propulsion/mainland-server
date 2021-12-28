package types

import (
	"gorm.io/gorm"
)

type Node struct {
	gorm.Model
	Name        string
	Description string
	NodeID      uint8
}
