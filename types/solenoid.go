package types

import (
	"gorm.io/gorm"
)

type Solenoid struct {
	gorm.Model
	Name        string
	Description string
	CANID       uint8
}
