package systems

import (
	"testing"

	"github.com/Liquid-Propulsion/mainland-server/canbackend"
	"github.com/Liquid-Propulsion/mainland-server/config"
	"github.com/Liquid-Propulsion/mainland-server/database/sql"
	"github.com/Liquid-Propulsion/mainland-server/database/timeseries"
	"github.com/Liquid-Propulsion/mainland-server/types"
	"github.com/stretchr/testify/assert"
)

func TestEngine(t *testing.T) {
	config.Init()
	canbackend.Init(config.MOCK)
	sql.Init()
	timeseries.Init("./")
	Init()

	assert.Equal(t, CurrentEngine.State(), types.SAFE, "Engine is not safe.")

	CurrentEngine.SetState(types.ARMED)

	assert.Equal(t, CurrentEngine.State(), types.SAFE, "Engine is not safe (Is locked out).")

	CurrentEngine.SetState(types.TEST)

	assert.Equal(t, CurrentEngine.State(), types.TEST, "Engine is not in test mode.")
}
