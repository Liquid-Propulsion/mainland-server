package types

import (
	"gorm.io/gorm"
)

type SafetyCheck struct {
	gorm.Model
	ID          *uint
	Name        string
	Description string
	ValidState  EngineState
	Code        string
}
