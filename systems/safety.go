package systems

import (
	"log"

	"github.com/Liquid-Propulsion/mainland-server/database/sql"
	"github.com/Liquid-Propulsion/mainland-server/types"
)

type SafetySystem struct {
	safetyChecks []types.SafetyCheck
}

func NewSafetySystem() *SafetySystem {
	safety := new(SafetySystem)
	safety.safetyChecks = make([]types.SafetyCheck, 0)
	return safety
}

func (safety *SafetySystem) LoadChecks() {
	res := sql.Database.Find(&safety.safetyChecks)
	if res.Error != nil {
		log.Printf("Couldn't query for safety checks: %s", res.Error)
	}
}

func (safety *SafetySystem) Tick(state types.EngineState, sensors *SensorsSystem) {
	for _, check := range safety.safetyChecks {
		if state == check.ValidState || check.ValidState == types.ALL {
			safety.RunCheck(check, sensors)
		}
	}
}

func (safety *SafetySystem) RunCheck(check types.SafetyCheck, sensors *SensorsSystem) error {
	return CurrentEngine.EvalSystem.RunCode(check.Code)
}
