package mock

import (
	"log"
	"math"
	"time"

	canpackets "github.com/Liquid-Propulsion/canpackets/go"
)

var fakeNodeIDs = []uint8{0, 1, 2, 3}
var fakeSensorIDs = []uint8{0, 1, 2}

type MockCANBackend struct {
	sensorDataChannel chan canpackets.SensorDataPacket
	pongChannel       chan canpackets.PongPacket
}

func New() *MockCANBackend {
	backend := new(MockCANBackend)
	backend.sensorDataChannel = make(chan canpackets.SensorDataPacket)
	backend.pongChannel = make(chan canpackets.PongPacket)
	go func() {
		for {
			for _, sensorID := range fakeSensorIDs {
				backend.sensorDataChannel <- canpackets.SensorDataPacket{
					SensorId:   sensorID,
					SensorData: uint32(math.Abs(10000 + (math.Sin(float64(time.Now().UnixMilli())/10000.0) * 1000))),
				}
			}
			time.Sleep(time.Millisecond * 20)
		}
	}()
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

func (backend *MockCANBackend) SendStage(packet canpackets.StagePacket) error {
	log.Printf("MOCKCAN: Set State to %v", packet.SolenoidState)
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
			NodeId:   node,
			NodeType: canpackets.SOLENOID_NODE,
		}
	}
	return nil
}

func (backend *MockCANBackend) SendPower(power canpackets.PowerPacket) error {
	log.Printf("MOCKCAN: Power is ", power.SystemPowered)
	return nil
}
