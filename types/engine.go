package types

type EngineState uint8

const (
	SAFE EngineState = iota
	ARMED
	TEST
	ALL
)
