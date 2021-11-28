package sql

import "gorm.io/gorm"

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

type Sensor struct {
	NodeID   uint8
	SensorID uint8
}

type Condition struct {
	Sensor1       Sensor
	ConditionType ConditionType
	Sensor2       Sensor
}

type SafetyCheck struct {
	gorm.Model
	Name        string
	Description string
	Conditions  []Condition
	Action      Action
}
