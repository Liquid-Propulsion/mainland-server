package types

import (
	"time"

	"gorm.io/gorm"
)

type Stage struct {
	gorm.Model
	Name         string
	Description  string
	CANID        uint8 `gorm:"column:can_"`
	PreStageCode string
	Duration     time.Duration
}
