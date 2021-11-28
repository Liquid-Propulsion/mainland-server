package sql

import "gorm.io/gorm"

type SOP struct {
	gorm.Model
	Name        string
	Description string
	Conditions  []Condition
}
