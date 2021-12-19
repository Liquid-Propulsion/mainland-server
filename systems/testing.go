package systems

import (
	"log"

	canpackets "github.com/Liquid-Propulsion/canpackets/go"
	"github.com/Liquid-Propulsion/mainland-server/canbackend"
	"github.com/Liquid-Propulsion/mainland-server/database/sql"
	"github.com/Liquid-Propulsion/mainland-server/types"
)

type TestSystem struct {
	solenoids map[uint8]canpackets.SolenoidState
}

func NewTestSystem() *TestSystem {
	test := new(TestSystem)
	test.solenoids = make(map[uint8]canpackets.SolenoidState)
	return test
}

func (test *TestSystem) Reset() {
	var solenoids []types.Solenoid
	res := sql.Database.Find(&solenoids)
	if res.Error != nil {
		log.Printf("Couldn't query for solenoids: %s", res.Error)
	}
	test.solenoids = make(map[uint8]canpackets.SolenoidState)
	for _, solenoid := range solenoids {
		test.solenoids[solenoid.CANID] = canpackets.CLOSED
	}
}

func (test *TestSystem) Tick() {
	for id, state := range test.solenoids {
		err := canbackend.CurrentCANBackend.SendSolenoidCommand(canpackets.SolenoidStatePacket{
			Id:    canpackets.ID(id),
			State: state,
		})
		log.Printf("Couldn't send solenoid command: %s", err)
	}
}

func (test *TestSystem) SetSolenoidState(solenoidID uint8, state canpackets.SolenoidState) {
	test.solenoids[solenoidID] = state
}
