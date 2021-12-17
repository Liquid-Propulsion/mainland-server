// This package implements the core systems of the Mainland server. This includes the SOP system, Safety system, Staging system, Test system,
// Sensor system, and the engine.
package systems

import (
	"errors"
	"log"
	"time"

	canpackets "github.com/Liquid-Propulsion/canpackets/go"
	"github.com/Liquid-Propulsion/mainland-server/canbackend"
	"github.com/Liquid-Propulsion/mainland-server/lockout"
	"github.com/stianeikeland/go-rpio/v4"
)

type EngineState uint8

const (
	SAFE EngineState = iota
	ARMED
	TEST
)

type Engine struct {
	state            EngineState
	LockoutSystem    *lockout.Lockout
	SensorsSystem    *SensorsSystem
	StagingSystem    *StagingSystem
	SafetySystem     *SafetySystem
	TestSystem       *TestSystem
	previousTickTime time.Time
	hasRPIO          bool
	testButtonPin    rpio.Pin
}

var CurrentEngine *Engine

func Init() {
	engine := new(Engine)
	engine.state = SAFE
	engine.LockoutSystem = lockout.New()
	engine.SensorsSystem = NewSensorsSystem()
	engine.StagingSystem = NewStagingSystem()
	engine.SafetySystem = NewSafetySystem()
	engine.TestSystem = NewTestSystem()
	engine.previousTickTime = time.Now()
	engine.hasRPIO = false
	engine.testButtonPin = 0
	engine.start()
	CurrentEngine = engine
}

func (engine *Engine) start() {
	err := rpio.Open()
	if err != nil {
		log.Printf("couldn't open test button, assuming this is in development mode: %s", err)
	} else {
		engine.hasRPIO = true
		engine.testButtonPin = rpio.Pin(10)
		engine.testButtonPin.Input()
	}
	engine.StagingSystem.LoadStages()
	engine.SafetySystem.LoadChecks()
	engine.TestSystem.Reset()
	go engine.LockoutSystem.Run()
	go engine.tickLoop()
}

func (engine *Engine) tickLoop() {
	// Run all Safety Checks in the Safety System, disabling the system if necessary
	if engine.SafetySystem.Tick(engine.state, engine.SensorsSystem) {
		engine.SetState(SAFE)
		return
	}
	for {
		switch engine.state {
		case ARMED:
			if engine.LockoutSystem.LockedOut() {
				// If a lockout key is removed while armed, the system is automatically safed.
				engine.SetState(SAFE)
				break
			}
			engine.StagingSystem.DecrementTime(time.Since(engine.previousTickTime))
			if !engine.StagingSystem.HasTimeLeft() {
				engine.StagingSystem.NextStage()
			}
			err := canbackend.CurrentCANBackend.SendPower(canpackets.PowerPacket{
				SystemPowered: true,
			})
			if err != nil {
				log.Printf("Couldn't send power packet: %s", err)
			}
			err = canbackend.CurrentCANBackend.SendStage(canpackets.StagePacket{
				SystemReady: true,
				Stage:       canpackets.Stage(engine.StagingSystem.GetCurrentStage().CANID),
			})
			if err != nil {
				log.Printf("Couldn't send stage packet: %s", err)
			}
		case TEST:
			if !engine.TestButtonHeld() {
				engine.SetState(SAFE)
				break
			}
			err := canbackend.CurrentCANBackend.SendPower(canpackets.PowerPacket{
				SystemPowered: true,
			})
			if err != nil {
				log.Printf("Couldn't send power packet: %s", err)
			}
			engine.TestSystem.Tick()
		}
		engine.previousTickTime = time.Now()
		// There are 50 ticks in every second.
		time.Sleep(time.Millisecond * 20)
	}
}

func (engine *Engine) TestButtonHeld() bool {
	if engine.hasRPIO {
		return engine.testButtonPin.Read() == rpio.High
	}
	return true
}

func (engine *Engine) HasRPIO() bool {
	return engine.hasRPIO
}

func (engine *Engine) SetState(state EngineState) error {
	switch state {
	case ARMED:
		if engine.LockoutSystem.LockedOut() {
			return errors.New("the engine is currently locked out: check that all keys are returned")
		}
	case TEST:
		if engine.LockoutSystem.LockedOut() && !engine.TestButtonHeld() {
			return errors.New("when the engine is locked out, the test button must be held to go into test state")
		}
		engine.TestSystem.Reset()
	}
	engine.state = state
	return nil
}

func (engine *Engine) State() EngineState {
	return engine.state
}
