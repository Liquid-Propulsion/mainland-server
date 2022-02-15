package types

import (
	"gorm.io/gorm"
)

type SafetyCheck struct {
	gorm.Model
	Name        string
	Description string
	ValidState  EngineState
	Code        string
}

func (node *SafetyCheck) IsNode() {}
