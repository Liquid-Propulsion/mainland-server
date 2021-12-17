package sql

import (
	"gorm.io/gorm"
)

type SOPCondition struct {
	gorm.Model
	ID             *uint
	SOPID          *uint
	SensorNodeID   uint8
	SensorID       uint8
	ConditionType  ConditionType
	ConditionValue uint32
	Reason         string
}

type SOP struct {
	gorm.Model
	ID          *uint
	Name        string
	Description string
	Conditions  []SOPCondition `gorm:"foreignkey:SOPID"`
}
