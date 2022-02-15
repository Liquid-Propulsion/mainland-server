package tcp

import (
	"errors"

	canpackets "github.com/Liquid-Propulsion/canpackets/go"
	"github.com/angelodlfrtr/go-can"
	"github.com/angelodlfrtr/go-can/transports"
)

type TCPCANBackend struct {
	bus               *can.Bus
	sensorDataChannel chan canpackets.SensorDataPacket
	pongChannel       chan canpackets.PongPacket
}

func New(host string, port int) (*TCPCANBackend, error) {
	backend := new(TCPCANBackend)
	tr := &transports.TCPCan{
		Host: host,
		Port: port,
	}

	bus := can.NewBus(tr)

	if err := bus.Open(); err != nil {
		return nil, err
	}

	backend.bus = bus
	backend.sensorDataChannel = make(chan canpackets.SensorDataPacket)
	backend.pongChannel = make(chan canpackets.PongPacket)
	go backend.run()
	return backend, nil
}

func (backend *TCPCANBackend) run() {
	for {
		frame := <-backend.bus.ReadChan()
		if frame == nil {
			break
		}
		switch frame.ArbitrationID {
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
}

func (backend *TCPCANBackend) Utilization() (float32, error) {
	return 1.0, nil
}

func (backend *TCPCANBackend) SensorDataChannel() chan canpackets.SensorDataPacket {
	return backend.sensorDataChannel
}

func (backend *TCPCANBackend) PongChannel() chan canpackets.PongPacket {
	return backend.pongChannel
}

func (backend *TCPCANBackend) SendStage(packet canpackets.StagePacket) error {
	frame, err := createFrame(0x00, packet.Encode())
	if err != nil {
		return err
	}
	return backend.bus.Write(frame)
}

func (backend *TCPCANBackend) SendBlink(packet canpackets.BlinkPacket) error {
	frame, err := createFrame(0x06, packet.Encode())
	if err != nil {
		return err
	}
	return backend.bus.Write(frame)
}

func (backend *TCPCANBackend) SendPing() error {
	frame, err := createFrame(0x04, []byte{})
	if err != nil {
		return err
	}
	return backend.bus.Write(frame)
}
func (backend *TCPCANBackend) SendPower(power canpackets.PowerPacket) error {
	frame, err := createFrame(0x00, power.Encode())
	if err != nil {
		return err
	}
	return backend.bus.Write(frame)
}

func createFrame(id uint32, src []byte) (*can.Frame, error) {
	length := len(src)
	if length > 8 {
		return nil, errors.New("packet too long")
	}
	var dst [8]uint8
	copy(dst[:], src[:length])
	frame := can.Frame{
		ArbitrationID: id,
		DLC:           uint8(length),
		Data:          dst,
	}
	return &frame, nil
}
