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
	"github.com/Liquid-Propulsion/mainland-server/types"
)

type Engine struct {
	state            types.EngineState
	LockoutSystem    *lockout.Lockout
	SensorsSystem    *SensorsSystem
	StagingSystem    *StagingSystem
	SafetySystem     *SafetySystem
	TestSystem       *TestSystem
	EvalSystem       *EvalSystem
	TestButton       *TestButton
	previousTickTime time.Time
}

var CurrentEngine *Engine

func Init() {
	engine := new(Engine)
	engine.state = types.SAFE
	engine.LockoutSystem = lockout.New()
	engine.TestButton = NewTestButton()
	engine.SensorsSystem = NewSensorsSystem()
	engine.StagingSystem = NewStagingSystem()
	engine.SafetySystem = NewSafetySystem()
	engine.TestSystem = NewTestSystem()
	engine.EvalSystem = NewEvalSystem()
	engine.previousTickTime = time.Now()
	engine.start()
	CurrentEngine = engine
}

func (engine *Engine) start() {
	engine.StagingSystem.Reset()
	engine.SafetySystem.Reset()
	engine.TestSystem.Reset()
	go engine.LockoutSystem.Run()
	go engine.tickLoop()
}

func (engine *Engine) tickLoop() {
	for {
		// Run all Safety Checks in the Safety System, disabling the system if necessary
		engine.SafetySystem.Tick(engine.state, engine.SensorsSystem)
		if engine.state != types.ESTOP {
			err := canbackend.CurrentCANBackend.SendPower(canpackets.PowerPacket{
				SystemPowered: true,
			})
			if err != nil {
				log.Printf("Couldn't send power packet: %s", err)
			}
		}
		switch engine.state {
		case types.ARMED:
			if engine.LockoutSystem.LockedOut() {
				// If a lockout key is removed while armed, the system is automatically safed.
				engine.SetState(types.SAFE)
				break
			}
			engine.StagingSystem.DecrementTime(time.Since(engine.previousTickTime))
			if !engine.StagingSystem.HasTimeLeft() {
				if engine.StagingSystem.NextStage() == nil {
					engine.SetState(types.SAFE)
				}
			}
			err := canbackend.CurrentCANBackend.SendPower(canpackets.PowerPacket{
				SystemPowered: true,
			})
			if err != nil {
				log.Printf("Couldn't send power packet: %s", err)
			}
			stage := engine.StagingSystem.GetCurrentStage()
			if stage != nil {
				err = canbackend.CurrentCANBackend.SendStage(canpackets.StagePacket{
					SystemReady: true,
					Stage:       canpackets.Stage(engine.StagingSystem.GetCurrentStage().CANID),
				})
				if err != nil {
					log.Printf("Couldn't send stage packet: %s", err)
				}
			}
		case types.TEST:
			if !engine.TestButtonHeld() {
				engine.SetState(types.SAFE)
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
	return engine.TestButton.ButtonHeld()
}

func (engine *Engine) HasRPIO() bool {
	return engine.TestButton.HasRPIO()
}

func (engine *Engine) SetState(state types.EngineState) error {
	switch state {
	case types.SAFE:
		engine.StagingSystem.Reset()
	case types.ARMED:
		if engine.LockoutSystem.LockedOut() {
			return errors.New("the engine is currently locked out: check that all keys are returned")
		}
	case types.TEST:
		if engine.LockoutSystem.LockedOut() && !engine.TestButtonHeld() {
			return errors.New("when the engine is locked out, the test button must be held to go into test state")
		}
		engine.TestSystem.Reset()
	}
	engine.state = state
	return nil
}

func (engine *Engine) EngineInfo() types.Engine {
	return types.Engine{
		EngineState:    engine.state,
		IsRpio:         engine.HasRPIO(),
		TestButtonHeld: engine.TestButtonHeld(),
		LockoutEnabled: engine.LockoutSystem.LockedOut(),
	}
}

func (engine *Engine) State() types.EngineState {
	return engine.state
}
