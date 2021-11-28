package mock

import (
	canpackets "github.com/Liquid-Propulsion/canpackets/go"
)

type MockCANBackend struct {
	sensorDataChannel chan canpackets.SensorDataPacket
	pongChannel       chan canpackets.PongPacket
}

func New() *MockCANBackend {
	backend := new(MockCANBackend)
	backend.sensorDataChannel = make(chan canpackets.SensorDataPacket)
	backend.pongChannel = make(chan canpackets.PongPacket)
	return backend
}

func (backend *MockCANBackend) Utilization() (float32, error) {
	return 1.0, nil
}

func (backend *MockCANBackend) SensorDataChannel() chan canpackets.SensorDataPacket {
	return backend.sensorDataChannel
}

func (backend *MockCANBackend) PongChannel() chan canpackets.PongPacket {
	return backend.pongChannel
}

func (backend *MockCANBackend) SendSolenoidCommand(packet canpackets.SolenoidStatePacket) error {
	return nil
}

func (backend *MockCANBackend) SendStage(packet canpackets.StagePacket) error {
	return nil
}

func (backend *MockCANBackend) SendBlink(packet canpackets.BlinkPacket) error {
	return nil
}

func (backend *MockCANBackend) SendPing() error {
	return nil
}

func (backend *MockCANBackend) SendPower(power canpackets.PowerPacket) error {
	return nil
}
