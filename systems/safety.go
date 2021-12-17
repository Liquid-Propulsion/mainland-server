package systems

import (
	"fmt"
	"log"

	"github.com/Liquid-Propulsion/mainland-server/database/sql"
)

type SafetySystem struct {
	safetyChecks  []sql.SafetyCheck
	warnings      []string
	info          []string
	reasonForStop string
}

func NewSafetySystem() *SafetySystem {
	safety := new(SafetySystem)
	safety.safetyChecks = make([]sql.SafetyCheck, 0)
	safety.warnings = make([]string, 0)
	safety.info = make([]string, 0)
	safety.reasonForStop = "none"
	return safety
}

func (safety *SafetySystem) LoadChecks() {
	res := sql.Database.Find(&safety.safetyChecks)
	if res.Error != nil {
		log.Printf("Couldn't query for safety checks: %s", res.Error)
	}
}

func (safety *SafetySystem) Reset() {
	safety.warnings = make([]string, 0)
	safety.info = make([]string, 0)
	safety.reasonForStop = ""
}

func (safety *SafetySystem) Tick(state EngineState, sensors *SensorsSystem) bool {
	for _, check := range safety.safetyChecks {
		if check.ValidState == sql.ALL {
			return safety.RunCheck(check, sensors)
		}
		if EngineStateToDBState(state) == check.ValidState {
			return safety.RunCheck(check, sensors)
		}
	}
	return false
}

func (safety *SafetySystem) RunCheck(check sql.SafetyCheck, sensors *SensorsSystem) bool {
	for _, condition := range check.Conditions {
		conditionMet := safety.RunCondition(condition, sensors)
		if conditionMet {
			switch check.Action {
			case sql.STOPSYSTEM:
				safety.reasonForStop = fmt.Sprintf("%s: %s", check.Name, condition.Reason)
				return true
			case sql.WARN:
				safety.warnings = append(safety.warnings, fmt.Sprintf("%s: %s", check.Name, condition.Reason))
			case sql.INFO:
				safety.info = append(safety.info, fmt.Sprintf("%s: %s", check.Name, condition.Reason))
			}
		}
	}
	return false
}

func (safety *SafetySystem) RunCondition(condition sql.SafetyCondition, sensors *SensorsSystem) bool {
	data, err := sensors.GetLatestSensorData(condition.SensorNodeID, condition.SensorID)
	if err != nil {
		return false
	}
	switch condition.ConditionType {
	case sql.EQUAL:
		return data.SensorData == condition.ConditionValue
	case sql.GREATERTHAN:
		return data.SensorData > condition.ConditionValue
	case sql.LESSTHAN:
		return data.SensorData < condition.ConditionValue
	case sql.GREATERTHANEQUAL:
		return data.SensorData >= condition.ConditionValue
	case sql.LESSTHANEQUAL:
		return data.SensorData <= condition.ConditionValue
	case sql.NOTEQUAL:
		return data.SensorData != condition.ConditionValue
	}
	return false
}

func EngineStateToDBState(state EngineState) sql.State {
	switch state {
	case SAFE:
		return sql.SAFE
	case ARMED:
		return sql.ARMED
	case TEST:
		return sql.TEST
	}
	return sql.ALL
}
