package canbackend

import (
	canpackets "github.com/Liquid-Propulsion/canpackets/go"
)

type CANBackend interface {
	// This is calculated by the following equation  ((inbound packets per second + outbound packets per second) * 76 bits) / 500000 bps
	// That equation is ran on every function call, and is based upon the current saved values for packets per second.
	Utilization() (float32, error)
	// Returns a channel that returns sensor data when it's recieved.
	SensorDataChannel() chan canpackets.SensorDataPacket
	// Returns a channel that returns a pongs of Island Nodes
	PongChannel() chan canpackets.PongPacket
	// Sends a Stage Command
	SendStage(canpackets.StagePacket) error
	// Sends a Power Command
	SendPower(canpackets.PowerPacket) error
	// Sends a Blink Command
	SendBlink(canpackets.BlinkPacket) error
	// Sends a Ping Command
	SendPing() error
}
