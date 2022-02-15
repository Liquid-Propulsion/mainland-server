package types

import (
	"gorm.io/gorm"
)

type SOP struct {
	gorm.Model
	Name        string
	Description string
	Code        string
}

func (node *SOP) IsNode() {}
