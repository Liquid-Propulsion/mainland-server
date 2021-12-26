package tests

import (
	"log"
	"testing"
	"time"

	"github.com/Liquid-Propulsion/mainland-server/canbackend"
	"github.com/Liquid-Propulsion/mainland-server/config"
	"github.com/Liquid-Propulsion/mainland-server/database/sql"
	"github.com/Liquid-Propulsion/mainland-server/database/timeseries"
	"github.com/Liquid-Propulsion/mainland-server/systems"
	"github.com/Liquid-Propulsion/mainland-server/types"
	"github.com/stretchr/testify/assert"
)

func initTestState() {
	config.Init()
	canbackend.Init(config.MOCK)
	sql.Init("file::memory:?cache=shared")
	timeseries.Init(true, "")
	systems.Init()
}

// TestDefaultEngineState ensures that by default the engine state is correct.
func TestDefaultEngineState(t *testing.T) {
	initTestState()

	assert.Equal(t, systems.CurrentEngine.State(), types.SAFE, "Engine is not safe.")

	systems.CurrentEngine.SetState(types.ARMED) // Lockout system is enabled, so it cannot be enabled yet.

	assert.Equal(t, systems.CurrentEngine.State(), types.SAFE, "Engine is not safe (While locked out).")

	systems.CurrentEngine.SetState(types.TEST)

	assert.Equal(t, systems.CurrentEngine.State(), types.TEST, "Engine is not in test mode.")

	// Disable Lockout Safety
	systems.CurrentEngine.LockoutSystem.SetLockoutEnabled(false)

	systems.CurrentEngine.SetState(types.ARMED)

	assert.Equal(t, systems.CurrentEngine.State(), types.ARMED, "Engine is not armed.")

	systems.CurrentEngine.SetState(types.TEST)

	assert.Equal(t, systems.CurrentEngine.State(), types.SAFE, "Engine is not in test mode.")

	// Reenable Lockout Safety and test

	systems.CurrentEngine.SetState(types.ARMED)

	systems.CurrentEngine.LockoutSystem.SetLockoutEnabled(true)

	time.Sleep(time.Millisecond * 40)

	log.Println(systems.CurrentEngine.State())

	assert.Equal(t, systems.CurrentEngine.State(), types.SAFE, "Engine is not safe.")
}

func TestLockout(t *testing.T) {
	initTestState()

	assert.Equal(t, systems.CurrentEngine.LockoutSystem.LockedOut(), true, "Engine not locked out.")

	systems.CurrentEngine.LockoutSystem.FakePacketRecieved()

	assert.Equal(t, systems.CurrentEngine.LockoutSystem.LockedOut(), false, "Engine locked out.")

	time.Sleep(time.Millisecond * 41) // Check timeout

	assert.Equal(t, systems.CurrentEngine.LockoutSystem.LockedOut(), true, "Engine not locked out.")
}
