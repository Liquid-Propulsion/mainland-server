package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SolenoidState [64]bool

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *SolenoidState) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	err := json.Unmarshal(bytes, &j)
	return err
}

// Value return json value, implement driver.Valuer interface
func (j SolenoidState) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(&j)
}

type Stage struct {
	gorm.Model
	Name          string
	Description   string
	SolenoidState SolenoidState
	Duration      time.Duration
}

func (node *Stage) IsNode() {}
