package sql

import (
	"gorm.io/gorm"
)

type ConditionType string

const (
	EQUAL            ConditionType = "EQUAL"
	NOTEQUAL         ConditionType = "NOTEQUAL"
	GREATERTHAN      ConditionType = "GREATERTHAN"
	GREATERTHANEQUAL ConditionType = "GREATERTHANEQUAL"
	LESSTHAN         ConditionType = "LESSTHAN"
	LESSTHANEQUAL    ConditionType = "LESSTHANEQUAL"
)

type Action string

const (
	STOPSYSTEM Action = "STOPSYSTEM"
	WARN       Action = "WARN"
	INFO       Action = "INFO"
)

type State string

const (
	ALL   State = "ALL"
	SAFE  State = "SAFE"
	ARMED State = "ARMED"
	TEST  State = "TEST"
)

type SafetyCondition struct {
	gorm.Model
	ID             *uint
	SafetyCheckID  *uint
	SensorNodeID   uint8
	SensorID       uint8
	ConditionType  ConditionType
	ConditionValue uint32
	Reason         string
}

type SafetyCheck struct {
	gorm.Model
	ID          *uint
	Name        string
	Description string
	ValidState  State
	Conditions  []SafetyCondition `gorm:"foreignkey:SafetyCheckID"`
	Action      Action
}
