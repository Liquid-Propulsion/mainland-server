package mock

import (
	"log"

	canpackets "github.com/Liquid-Propulsion/canpackets/go"
)

var fakeNodeIDs = []uint8{0, 1, 2, 3}

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
	log.Printf("MOCKCAN: Set Solenoid %d to %d", packet.Id, packet.State)
	return nil
}

func (backend *MockCANBackend) SendStage(packet canpackets.StagePacket) error {
	log.Printf("MOCKCAN: Set State to %d", packet.Stage)
	return nil
}

func (backend *MockCANBackend) SendBlink(packet canpackets.BlinkPacket) error {
	log.Printf("MOCKCAN: Blink Node %d", packet.NodeId)
	return nil
}

func (backend *MockCANBackend) SendPing() error {
	log.Printf("MOCKCAN: Ping all Nodes")
	for _, node := range fakeNodeIDs {
		backend.pongChannel <- canpackets.PongPacket{
			NodeId:   canpackets.ID(node),
			NodeType: canpackets.SOLENOID_NODE,
		}
	}
	return nil
}

func (backend *MockCANBackend) SendPower(power canpackets.PowerPacket) error {
	log.Printf("MOCKCAN: Power is online")
	return nil
}
