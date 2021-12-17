package sql

import (
	"time"

	"gorm.io/gorm"
)

type Stage struct {
	gorm.Model
	ID          *uint
	Name        string
	Description string
	CANID       uint8
	Duration    time.Duration
}
