package types

import (
	"fmt"
	"io"
	"strconv"
)

type EngineState string

const (
	SAFE  EngineState = "SAFE"
	ARMED EngineState = "ARMED"
	TEST  EngineState = "TEST"
	ESTOP EngineState = "ESTOP"
	ALL   EngineState = "ALL"
)

var AllEngineState = []EngineState{
	SAFE,
	ARMED,
	TEST,
	ESTOP,
}

func (e EngineState) IsValid() bool {
	switch e {
	case SAFE, ARMED, TEST, ESTOP:
		return true
	}
	return false
}

func (e EngineState) String() string {
	return string(e)
}

func (e *EngineState) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EngineState(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EngineState", str)
	}
	return nil
}

func (e EngineState) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
