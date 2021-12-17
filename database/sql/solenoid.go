package sql

import (
	"gorm.io/gorm"
)

type Solenoid struct {
	gorm.Model
	ID          *uint
	Name        string
	Description string
	CANID       uint8
}
