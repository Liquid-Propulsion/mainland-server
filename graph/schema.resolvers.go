package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/Liquid-Propulsion/mainland-server/database/sql"
	"github.com/Liquid-Propulsion/mainland-server/graph/generated"
	"github.com/Liquid-Propulsion/mainland-server/systems"
	"github.com/Liquid-Propulsion/mainland-server/types"
)

func (r *mutationResolver) SetEngineState(ctx context.Context, state types.EngineState) (*types.Engine, error) {
	err := systems.CurrentEngine.SetState(state)
	info := systems.CurrentEngine.EngineInfo()
	return &info, err
}

func (r *mutationResolver) CreateUser(ctx context.Context, user types.CreateUserInput) (*types.User, error) {
	hash, err := EncodePassword(user.Password)
	if err != nil {
		return nil, err
	}
	userType := types.User{
		Name:         user.Name,
		Username:     user.Username,
		PasswordHash: hash,
	}
	tx := sql.Database.Create(&userType)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &userType, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, user types.UpdateUserInput) (*types.User, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	userOut := types.User{
		ID: &idInt,
	}
	tx := sql.Database.Where(&userOut).Take(&userOut)
	if tx.Error != nil {
		return nil, tx.Error
	}
	userOut.Name = user.Name
	userOut.Username = user.Username
	tx = sql.Database.Save(&userOut)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &userOut, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*types.User, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	user := types.User{
		ID: &idInt,
	}
	tx := sql.Database.Delete(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

func (r *mutationResolver) CreateStage(ctx context.Context, stage types.StageInput) (*types.Stage, error) {
	parsedDuration, err := time.ParseDuration(stage.Duration)
	if err != nil {
		return nil, err
	}
	stageType := types.Stage{
		Name:         stage.Name,
		Description:  stage.Description,
		CANID:        uint8(stage.CanID),
		PreStageCode: stage.PreStageCode,
		Duration:     parsedDuration,
	}
	tx := sql.Database.Create(&stageType)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &stageType, nil
}

func (r *mutationResolver) UpdateStage(ctx context.Context, id string, stage types.StageInput) (*types.Stage, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	stageOut := types.Stage{
		ID: &idInt,
	}
	tx := sql.Database.Where(&stageOut).Take(&stageOut)
	if tx.Error != nil {
		return nil, tx.Error
	}
	parsedDuration, err := time.ParseDuration(stage.Duration)
	if err != nil {
		return nil, err
	}
	stageOut.Name = stage.Name
	stageOut.Description = stage.Description
	stageOut.CANID = uint8(stage.CanID)
	stageOut.PreStageCode = stage.PreStageCode
	stageOut.Duration = parsedDuration
	tx = sql.Database.Save(&stageOut)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &stageOut, nil
}

func (r *mutationResolver) DeleteStage(ctx context.Context, id string) (*types.Stage, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	stage := types.Stage{
		ID: &idInt,
	}
	tx := sql.Database.Delete(&stage)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &stage, nil
}

func (r *mutationResolver) CreateSolenoid(ctx context.Context, solenoid types.SolenoidInput) (*types.Solenoid, error) {
	solenoidType := types.Solenoid{
		Name:        solenoid.Name,
		Description: solenoid.Description,
		CANID:       uint8(solenoid.CanID),
	}
	tx := sql.Database.Create(&solenoidType)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &solenoidType, nil
}

func (r *mutationResolver) UpdateSolenoid(ctx context.Context, id string, solenoid types.SolenoidInput) (*types.Solenoid, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	solenoidOut := types.Solenoid{
		ID: &idInt,
	}
	tx := sql.Database.Where(&solenoidOut).Take(&solenoidOut)
	if tx.Error != nil {
		return nil, tx.Error
	}
	solenoidOut.Name = solenoid.Name
	solenoidOut.Description = solenoid.Description
	solenoidOut.CANID = uint8(solenoid.CanID)
	tx = sql.Database.Save(&solenoidOut)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &solenoidOut, nil
}

func (r *mutationResolver) DeleteSolenoid(ctx context.Context, id string) (*types.Solenoid, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	solenoid := types.Solenoid{
		ID: &idInt,
	}
	tx := sql.Database.Delete(&solenoid)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &solenoid, nil
}

func (r *mutationResolver) CreateSensor(ctx context.Context, sensor types.SensorInput) (*types.Sensor, error) {
	sensorType := types.Sensor{
		Name:          sensor.Name,
		Description:   sensor.Description,
		NodeID:        uint8(sensor.NodeID),
		SensorID:      uint8(sensor.SensorID),
		TransformCode: sensor.TransformCode,
	}
	tx := sql.Database.Create(&sensorType)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &sensorType, nil
}

func (r *mutationResolver) UpdateSensor(ctx context.Context, id string, sensor types.SensorInput) (*types.Sensor, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	sensorOut := types.Sensor{
		ID: &idInt,
	}
	tx := sql.Database.Where(&sensorOut).Take(&sensorOut)
	if tx.Error != nil {
		return nil, tx.Error
	}
	sensorOut.Name = sensor.Name
	sensorOut.Description = sensor.Description
	sensorOut.NodeID = uint8(sensor.NodeID)
	sensorOut.SensorID = uint8(sensor.SensorID)
	sensorOut.TransformCode = sensor.TransformCode
	tx = sql.Database.Save(&sensorOut)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &sensorOut, nil
}

func (r *mutationResolver) DeleteSensor(ctx context.Context, id string) (*types.Sensor, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	sensor := types.Sensor{
		ID: &idInt,
	}
	tx := sql.Database.Delete(&sensor)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &sensor, nil
}

func (r *mutationResolver) CreateSafetyCheck(ctx context.Context, check types.SafetyCheckInput) (*types.SafetyCheck, error) {
	safetyType := types.SafetyCheck{
		Name:        check.Name,
		Description: check.Description,
		ValidState:  check.ValidState,
		Code:        check.Code,
	}
	tx := sql.Database.Create(&safetyType)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &safetyType, nil
}

func (r *mutationResolver) UpdateSafetyCheck(ctx context.Context, id string, check types.SafetyCheckInput) (*types.SafetyCheck, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	checkOut := types.SafetyCheck{
		ID: &idInt,
	}
	tx := sql.Database.Where(&checkOut).Take(&checkOut)
	if tx.Error != nil {
		return nil, tx.Error
	}
	checkOut.Name = check.Name
	checkOut.Description = check.Description
	checkOut.ValidState = check.ValidState
	checkOut.Code = check.Code
	tx = sql.Database.Save(&checkOut)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &checkOut, nil
}

func (r *mutationResolver) DeleteSafetyCheck(ctx context.Context, id string) (*types.SafetyCheck, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	safety := types.SafetyCheck{
		ID: &idInt,
	}
	tx := sql.Database.Delete(&safety)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &safety, nil
}

func (r *queryResolver) Engine(ctx context.Context) (*types.Engine, error) {
	info := systems.CurrentEngine.EngineInfo()
	return &info, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*types.User, error) {
	var users []*types.User
	out := sql.Database.Find(&users)
	if out.Error != nil {
		return users, out.Error
	}
	return users, nil
}

func (r *queryResolver) Stages(ctx context.Context) ([]*types.Stage, error) {
	var stages []*types.Stage
	out := sql.Database.Find(&stages)
	if out.Error != nil {
		return stages, out.Error
	}
	return stages, nil
}

func (r *queryResolver) Solenoids(ctx context.Context) ([]*types.Solenoid, error) {
	var solenoids []*types.Solenoid
	out := sql.Database.Find(&solenoids)
	if out.Error != nil {
		return solenoids, out.Error
	}
	return solenoids, nil
}

func (r *queryResolver) Sensors(ctx context.Context) ([]*types.Sensor, error) {
	var sensors []*types.Sensor
	out := sql.Database.Find(&sensors)
	if out.Error != nil {
		return sensors, out.Error
	}
	return sensors, nil
}

func (r *queryResolver) SafetyChecks(ctx context.Context) ([]*types.SafetyCheck, error) {
	var safetyChecks []*types.SafetyCheck
	out := sql.Database.Find(&safetyChecks)
	if out.Error != nil {
		return safetyChecks, out.Error
	}
	return safetyChecks, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
