package systems

import (
	"log"
	"time"

	"github.com/Liquid-Propulsion/mainland-server/database/sql"
)

type StagingSystem struct {
	stages          []sql.Stage
	currentStage    int
	timeLeftInStage time.Duration
}

func NewStagingSystem() *StagingSystem {
	staging := new(StagingSystem)
	staging.stages = make([]sql.Stage, 0)
	staging.currentStage = 0
	staging.timeLeftInStage = time.Microsecond
	return staging
}

func (staging *StagingSystem) LoadStages() {
	res := sql.Database.Order("can_id").Find(&staging.stages)
	if res.Error != nil {
		log.Printf("Couldn't query for stages: %s", res.Error)
	}
	if len(staging.stages) > 0 {
		staging.currentStage = 0
		staging.timeLeftInStage = staging.stages[0].Duration
	}
}

func (staging *StagingSystem) DecrementTime(by time.Duration) {
	staging.timeLeftInStage = staging.timeLeftInStage - by
}

func (staging *StagingSystem) HasTimeLeft() bool {
	return staging.timeLeftInStage > time.Millisecond
}

func (staging *StagingSystem) GetCurrentStage() *sql.Stage {
	if len(staging.stages) > staging.currentStage {
		return &staging.stages[staging.currentStage]
	}
	return nil
}

func (staging *StagingSystem) NextStage() *sql.Stage {
	if len(staging.stages) > staging.currentStage+1 {
		staging.currentStage += 1
		return &staging.stages[staging.currentStage]
	}
	return nil
}
