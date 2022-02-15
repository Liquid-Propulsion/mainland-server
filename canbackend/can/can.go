package can

import (
	"errors"

	canpackets "github.com/Liquid-Propulsion/canpackets/go"
	"github.com/brutella/can"
)

type OpenCANBackend struct {
	bus               *can.Bus
	sensorDataChannel chan canpackets.SensorDataPacket
	pongChannel       chan canpackets.PongPacket
}

func New(interfaceName string) (*OpenCANBackend, error) {
	backend := new(OpenCANBackend)
	bus, err := can.NewBusForInterfaceWithName(interfaceName)
	if err != nil {
		return nil, err
	}
	backend.bus = bus
	backend.sensorDataChannel = make(chan canpackets.SensorDataPacket)
	backend.pongChannel = make(chan canpackets.PongPacket)
	bus.Subscribe(backend)
	return backend, nil
}

func (backend *OpenCANBackend) Handle(frame can.Frame) {
	switch frame.ID {
	case 0x03:
		data := canpackets.SensorDataPacket{}
		data.Decode(frame.Data[:])
		backend.sensorDataChannel <- data
	case 0x05:
		data := canpackets.PongPacket{}
		data.Decode(frame.Data[:])
		backend.pongChannel <- data
	}
}

func (backend *OpenCANBackend) Utilization() (float32, error) {
	return 1.0, nil
}

func (backend *OpenCANBackend) SensorDataChannel() chan canpackets.SensorDataPacket {
	return backend.sensorDataChannel
}

func (backend *OpenCANBackend) PongChannel() chan canpackets.PongPacket {
	return backend.pongChannel
}

func (backend *OpenCANBackend) SendStage(packet canpackets.StagePacket) error {
	frame, err := createFrame(0x01, packet.Encode())
	if err != nil {
		return err
	}
	return backend.bus.Publish(frame)
}

func (backend *OpenCANBackend) SendBlink(packet canpackets.BlinkPacket) error {
	frame, err := createFrame(0x06, packet.Encode())
	if err != nil {
		return err
	}
	return backend.bus.Publish(frame)
}

func (backend *OpenCANBackend) SendPing() error {
	frame, err := createFrame(0x04, []byte{})
	if err != nil {
		return err
	}
	return backend.bus.Publish(frame)
}
func (backend *OpenCANBackend) SendPower(power canpackets.PowerPacket) error {
	frame, err := createFrame(0x00, power.Encode())
	if err != nil {
		return err
	}
	return backend.bus.Publish(frame)
}

func createFrame(id uint32, src []byte) (can.Frame, error) {
	length := len(src)
	if length > 8 {
		return can.Frame{}, errors.New("packet too long")
	}
	var dst [8]uint8
	copy(dst[:], src[:length])
	frame := can.Frame{
		ID:     id,
		Length: uint8(length),
		Flags:  0,
		Res0:   0,
		Res1:   0,
		Data:   dst,
	}
	return frame, nil
}
