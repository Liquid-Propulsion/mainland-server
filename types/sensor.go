package types

import (
	"gorm.io/gorm"
)

type Sensor struct {
	gorm.Model
	ID            *uint
	Name          string
	Description   string
	NodeID        uint8
	SensorID      uint8
	TransformCode string
}
