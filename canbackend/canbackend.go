package canbackend

import (
	"log"

	"github.com/Liquid-Propulsion/mainland-server/canbackend/can"
	"github.com/Liquid-Propulsion/mainland-server/canbackend/mock"
	"github.com/Liquid-Propulsion/mainland-server/config"
)

var CurrentCANBackend CANBackend

func Init(canType config.CANType) {
	switch canType {
	case config.OPENCAN:
		backend, err := can.New(config.CurrentConfig.CAN.InterfaceName)
		if err != nil {
			log.Printf("Cannot initialize OpenCAN bus: %s", err)
			break
		}
		log.Println("Successfully initialized OpenCAN bus.")
		CurrentCANBackend = backend
		return
	}
	log.Println("Successfully initialized Mock CAN Bus.")
	CurrentCANBackend = mock.New()
}
