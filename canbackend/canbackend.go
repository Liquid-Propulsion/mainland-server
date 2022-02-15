package canbackend

import (
	"log"

	"github.com/Liquid-Propulsion/mainland-server/canbackend/can"
	"github.com/Liquid-Propulsion/mainland-server/canbackend/mock"
	"github.com/Liquid-Propulsion/mainland-server/canbackend/tcp"
	"github.com/Liquid-Propulsion/mainland-server/canbackend/usb"
	"github.com/Liquid-Propulsion/mainland-server/config"
)

var CurrentCANBackend CANBackend

func Init(canType config.CANType) {
	switch canType {
	case config.SOCKETCAN:
		backend, err := can.New(config.CurrentConfig.CAN.InterfaceName)
		if err != nil {
			log.Printf("Cannot initialize OpenCAN bus: %s", err)
			break
		}
		log.Println("Successfully initialized OpenCAN bus.")
		CurrentCANBackend = backend
		return
	case config.TCP:
		backend, err := tcp.New(config.CurrentConfig.CAN.Host, config.CurrentConfig.CAN.Port)
		if err != nil {
			log.Printf("Cannot initialize TCP CAN bus: %s", err)
			break
		}
		log.Printf("Successfully initialized TCP CAN bus at %s:%d.", config.CurrentConfig.CAN.Host, config.CurrentConfig.CAN.Port)
		CurrentCANBackend = backend
		return
	case config.USB:
		backend, err := usb.New(config.CurrentConfig.CAN.SerialPort, config.CurrentConfig.CAN.BaudRate)
		if err != nil {
			log.Printf("Cannot initialize TCP CAN bus: %s", err)
			break
		}
		log.Printf("Successfully initialized USB CAN bus at %s with a baud rate of %d.", config.CurrentConfig.CAN.SerialPort, config.CurrentConfig.CAN.BaudRate)
		CurrentCANBackend = backend
		return
	}
	log.Println("Successfully initialized Mock CAN Bus.")
	CurrentCANBackend = mock.New()
}
