package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Liquid-Propulsion/mainland-server/database/sql"
	"github.com/Liquid-Propulsion/mainland-server/graph/generated"
	"github.com/Liquid-Propulsion/mainland-server/session"
	"github.com/Liquid-Propulsion/mainland-server/systems"
	"github.com/Liquid-Propulsion/mainland-server/types"
	"github.com/dustin/go-broadcast"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

var upsertBroadcast = broadcast.NewBroadcaster(10)
var deleteBroadcast = broadcast.NewBroadcaster(10)

func (r *mutationResolver) SetEngineState(ctx context.Context, state types.EngineState) (*types.Engine, error) {
	err := systems.CurrentEngine.SetState(state)
	info := systems.CurrentEngine.EngineInfo()
	return &info, err
}

func (r *mutationResolver) ResetEngine(ctx context.Context) (*types.Engine, error) {
	systems.CurrentEngine.Reset()
	info := systems.CurrentEngine.EngineInfo()
	return &info, nil
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
		TOTPEnabled:  false,
	}
	tx := sql.Database.Create(&userType)
	if tx.Error != nil {
		return nil, tx.Error
	}
	upsertBroadcast.Submit(EncodeID("user", userType.ID))
	return &userType, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, user types.UpdateUserInput) (*types.User, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	userOut := types.User{
		Model: gorm.Model{
			ID: idInt,
		},
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
	upsertBroadcast.Submit(id)
	return &userOut, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*types.User, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	user := types.User{
		Model: gorm.Model{
			ID: idInt,
		},
	}
	tx := sql.Database.Delete(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	deleteBroadcast.Submit(id)
	return &user, nil
}

func (r *mutationResolver) CreateStage(ctx context.Context, stage types.StageInput) (*types.Stage, error) {
	if len(stage.SolenoidState) != 64 {
		return nil, errors.New("solenoid_state must be 64 boolean values")
	}
	var solenoidState [64]bool
	copy(solenoidState[:], stage.SolenoidState[:64])
	parsedDuration, err := time.ParseDuration(stage.Duration)
	if err != nil {
		return nil, err
	}
	stageType := types.Stage{
		Name:          stage.Name,
		Description:   stage.Description,
		SolenoidState: solenoidState,
		Duration:      parsedDuration,
	}
	tx := sql.Database.Create(&stageType)
	if tx.Error != nil {
		return nil, tx.Error
	}
	upsertBroadcast.Submit(EncodeID("stage", stageType.ID))
	systems.CurrentEngine.StagingSystem.Reset()
	return &stageType, nil
}

func (r *mutationResolver) UpdateStage(ctx context.Context, id string, stage types.StageInput) (*types.Stage, error) {
	if len(stage.SolenoidState) != 64 {
		return nil, errors.New("solenoid_state must be 64 boolean values")
	}
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	stageOut := types.Stage{
		Model: gorm.Model{
			ID: idInt,
		},
	}
	tx := sql.Database.Where(&stageOut).Take(&stageOut)
	if tx.Error != nil {
		return nil, tx.Error
	}
	parsedDuration, err := time.ParseDuration(stage.Duration)
	if err != nil {
		return nil, err
	}
	var solenoidState [64]bool
	copy(solenoidState[:], stage.SolenoidState[:64])
	stageOut.Name = stage.Name
	stageOut.Description = stage.Description
	stageOut.SolenoidState = solenoidState
	stageOut.Duration = parsedDuration
	tx = sql.Database.Save(&stageOut)
	if tx.Error != nil {
		return nil, tx.Error
	}
	upsertBroadcast.Submit(id)
	systems.CurrentEngine.StagingSystem.Reset()
	return &stageOut, nil
}

func (r *mutationResolver) DeleteStage(ctx context.Context, id string) (*types.Stage, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	stage := types.Stage{
		Model: gorm.Model{
			ID: idInt,
		},
	}
	tx := sql.Database.Delete(&stage)
	if tx.Error != nil {
		return nil, tx.Error
	}
	deleteBroadcast.Submit(id)
	systems.CurrentEngine.StagingSystem.Reset()
	return &stage, nil
}

func (r *mutationResolver) CreateSolenoid(ctx context.Context, solenoid types.SolenoidInput) (*types.Solenoid, error) {
	id, err := cast.ToUintE(solenoid.ID)
	if err != nil {
		return nil, err
	}
	solenoidType := types.Solenoid{
		Model: types.Model{
			ID: id,
		},
		Name:        solenoid.Name,
		Description: solenoid.Description,
	}
	tx := sql.Database.Create(&solenoidType)
	if tx.Error != nil {
		return nil, tx.Error
	}
	upsertBroadcast.Submit(EncodeID("solenoid", id))
	return &solenoidType, nil
}

func (r *mutationResolver) UpdateSolenoid(ctx context.Context, id string, solenoid types.SolenoidInput) (*types.Solenoid, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	solenoidOut := types.Solenoid{
		Model: types.Model{
			ID: idInt,
		},
	}
	tx := sql.Database.Where(&solenoidOut).Take(&solenoidOut)
	if tx.Error != nil {
		return nil, tx.Error
	}
	solenoidOut.Name = solenoid.Name
	solenoidOut.Description = solenoid.Description
	tx = sql.Database.Save(&solenoidOut)
	if tx.Error != nil {
		return nil, tx.Error
	}
	upsertBroadcast.Submit(id)
	return &solenoidOut, nil
}

func (r *mutationResolver) DeleteSolenoid(ctx context.Context, id string) (*types.Solenoid, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	solenoid := types.Solenoid{
		Model: types.Model{
			ID: idInt,
		},
	}
	tx := sql.Database.Delete(&solenoid)
	if tx.Error != nil {
		return nil, tx.Error
	}
	deleteBroadcast.Submit(id)
	return &solenoid, nil
}

func (r *mutationResolver) CreateSensor(ctx context.Context, sensor types.SensorInput) (*types.Sensor, error) {
	id, err := cast.ToUintE(sensor.ID)
	if err != nil {
		return nil, err
	}
	sensorType := types.Sensor{
		Model: types.Model{
			ID: id,
		},
		Name:          sensor.Name,
		Description:   sensor.Description,
		TransformCode: sensor.TransformCode,
	}
	tx := sql.Database.Create(&sensorType)
	if tx.Error != nil {
		return nil, tx.Error
	}
	upsertBroadcast.Submit(EncodeID("sensor", id))
	systems.CurrentEngine.SensorsSystem.Reset()
	return &sensorType, nil
}

func (r *mutationResolver) UpdateSensor(ctx context.Context, id string, sensor types.SensorInput) (*types.Sensor, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	sensorOut := types.Sensor{
		Model: types.Model{
			ID: idInt,
		},
	}
	tx := sql.Database.Where(&sensorOut).Take(&sensorOut)
	if tx.Error != nil {
		return nil, tx.Error
	}
	sensorOut.Name = sensor.Name
	sensorOut.Description = sensor.Description
	sensorOut.TransformCode = sensor.TransformCode
	tx = sql.Database.Save(&sensorOut)
	if tx.Error != nil {
		return nil, tx.Error
	}
	upsertBroadcast.Submit(id)
	systems.CurrentEngine.SensorsSystem.Reset()
	return &sensorOut, nil
}

func (r *mutationResolver) DeleteSensor(ctx context.Context, id string) (*types.Sensor, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	sensor := types.Sensor{
		Model: types.Model{
			ID: idInt,
		},
	}
	tx := sql.Database.Delete(&sensor)
	if tx.Error != nil {
		return nil, tx.Error
	}
	deleteBroadcast.Submit(id)
	systems.CurrentEngine.SensorsSystem.Reset()
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
	upsertBroadcast.Submit(EncodeID("safety_check", safetyType.ID))
	systems.CurrentEngine.SafetySystem.Reset()
	return &safetyType, nil
}

func (r *mutationResolver) UpdateSafetyCheck(ctx context.Context, id string, check types.SafetyCheckInput) (*types.SafetyCheck, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	checkOut := types.SafetyCheck{
		Model: gorm.Model{
			ID: idInt,
		},
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
	upsertBroadcast.Submit(id)
	systems.CurrentEngine.SafetySystem.Reset()
	return &checkOut, nil
}

func (r *mutationResolver) DeleteSafetyCheck(ctx context.Context, id string) (*types.SafetyCheck, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	safety := types.SafetyCheck{
		Model: gorm.Model{
			ID: idInt,
		},
	}
	tx := sql.Database.Delete(&safety)
	if tx.Error != nil {
		return nil, tx.Error
	}
	deleteBroadcast.Submit(id)
	systems.CurrentEngine.SafetySystem.Reset()
	return &safety, nil
}

func (r *mutationResolver) CreateIslandNode(ctx context.Context, island types.IslandNodeInput) (*types.IslandNode, error) {
	id, err := cast.ToUintE(island.ID)
	if err != nil {
		return nil, err
	}
	islandType := types.IslandNode{
		Model: types.Model{
			ID: id,
		},
		Name:        island.Name,
		Description: island.Description,
	}
	tx := sql.Database.Create(&islandType)
	if tx.Error != nil {
		return nil, tx.Error
	}
	upsertBroadcast.Submit(EncodeID("node", id))
	systems.CurrentEngine.NodeSystem.Reset()
	return &islandType, nil
}

func (r *mutationResolver) UpdateIslandNode(ctx context.Context, id string, island types.IslandNodeInput) (*types.IslandNode, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	islandOut := types.IslandNode{
		Model: types.Model{
			ID: idInt,
		},
	}
	tx := sql.Database.Where(&islandOut).Take(&islandOut)
	if tx.Error != nil {
		return nil, tx.Error
	}
	islandOut.Name = island.Name
	islandOut.Description = island.Description
	tx = sql.Database.Save(&islandOut)
	if tx.Error != nil {
		return nil, tx.Error
	}
	upsertBroadcast.Submit(id)
	systems.CurrentEngine.NodeSystem.Reset()
	return &islandOut, nil
}

func (r *mutationResolver) DeleteIslandNode(ctx context.Context, id string) (*types.IslandNode, error) {
	_, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	island := types.IslandNode{
		Model: types.Model{
			ID: idInt,
		},
	}
	tx := sql.Database.Delete(&island)
	if tx.Error != nil {
		return nil, tx.Error
	}
	deleteBroadcast.Submit(id)
	systems.CurrentEngine.NodeSystem.Reset()
	return &island, nil
}

func (r *mutationResolver) SignIn(ctx context.Context, username string, password string) (*types.SignInResponse, error) {
	_, ok := session.ForContext(ctx)
	if ok {
		return nil, errors.New("your already logged in")
	}
	userOut := types.User{
		Username: username,
	}
	tx := sql.Database.Where(&userOut).Take(&userOut)
	if tx.Error != nil {
		time.Sleep(time.Second * 1)
		return nil, errors.New("username not found")
	}
	if !VerifyPassword(userOut.PasswordHash, password) {
		time.Sleep(time.Second * 1)
		return nil, errors.New("password invalid")
	}

	if userOut.TOTPEnabled {
		initial_session := session.CurrentSessionManager.Add(userOut.ID)
		new_session, _ := session.CurrentSessionManager.SetLock(initial_session.UUID.String(), true)
		return &types.SignInResponse{
			Type:        types.SignInResponseTypeTotp,
			Session:     new_session.UUID.String(),
			LockOutTime: new_session.LockTime.String(),
			ExpiryTime:  new_session.ExpiryTime.String(),
		}, nil
	} else {
		session := session.CurrentSessionManager.Add(userOut.ID)
		return &types.SignInResponse{
			Type:        types.SignInResponseTypeValid,
			Session:     session.UUID.String(),
			LockOutTime: session.LockTime.String(),
			ExpiryTime:  session.ExpiryTime.String(),
		}, nil
	}
}

func (r *mutationResolver) TotpVerify(ctx context.Context, code string) (*types.SignInResponse, error) {
	current_session, ok := session.ForContext(ctx)
	if ok {
		userOut := types.User{
			Model: gorm.Model{
				ID: current_session.UserID,
			},
		}
		tx := sql.Database.Where(&userOut).Take(&userOut)
		if tx.Error != nil {
			time.Sleep(time.Second * 1)
			return nil, errors.New("unknown user")
		}
		if current_session.Locked && userOut.TOTPEnabled {
			if totp.Validate(code, userOut.TOTPSecret) {
				new_session, err := session.CurrentSessionManager.SetLock(current_session.UUID.String(), false)
				if err != nil {
					return nil, err
				}
				return &types.SignInResponse{
					Type:        types.SignInResponseTypeValid,
					Session:     new_session.UUID.String(),
					LockOutTime: new_session.LockTime.String(),
					ExpiryTime:  new_session.ExpiryTime.String(),
				}, nil
			}
			return nil, errors.New("invalid totp code")
		}
		return nil, errors.New("totp not enabled for user")
	}
	return nil, errors.New("no session found")
}

func (r *mutationResolver) SignOut(ctx context.Context) (bool, error) {
	current_session, _ := session.ForContext(ctx)
	return true, session.CurrentSessionManager.Remove(current_session.UUID.String())
}

func (r *mutationResolver) LockOut(ctx context.Context) (*types.SignInResponse, error) {
	current_session, _ := session.ForContext(ctx)
	new_session, err := session.CurrentSessionManager.SetLock(current_session.UUID.String(), true)
	if err != nil {
		return nil, err
	}
	return &types.SignInResponse{
		Type:        types.SignInResponseTypeValid,
		Session:     new_session.UUID.String(),
		LockOutTime: new_session.LockTime.String(),
		ExpiryTime:  new_session.ExpiryTime.String(),
	}, nil
}

func (r *mutationResolver) PreventLockOut(ctx context.Context) (*types.SignInResponse, error) {
	current_session, _ := session.ForContext(ctx)
	new_session, err := session.CurrentSessionManager.PreventLock(current_session.UUID.String())
	if err != nil {
		return nil, err
	}
	return &types.SignInResponse{
		Type:        types.SignInResponseTypeValid,
		Session:     new_session.UUID.String(),
		LockOutTime: new_session.LockTime.String(),
		ExpiryTime:  new_session.ExpiryTime.String(),
	}, nil
}

func (r *queryResolver) Node(ctx context.Context, id string) (types.Node, error) {
	typeName, idInt, err := DecodeID(id)
	if err != nil {
		return nil, err
	}
	switch typeName {
	case "node":
		islandOut := types.IslandNode{
			Model: types.Model{
				ID: idInt,
			},
		}
		tx := sql.Database.Where(&islandOut).Take(&islandOut)
		if tx.Error != nil {
			return nil, tx.Error
		}
		return &islandOut, nil
	case "safety_check":
		safetyOut := types.SafetyCheck{
			Model: gorm.Model{
				ID: idInt,
			},
		}
		tx := sql.Database.Where(&safetyOut).Take(&safetyOut)
		if tx.Error != nil {
			return nil, tx.Error
		}
		return &safetyOut, nil
	case "sensor":
		sensorOut := types.Sensor{
			Model: types.Model{
				ID: idInt,
			},
		}
		tx := sql.Database.Where(&sensorOut).Take(&sensorOut)
		if tx.Error != nil {
			return nil, tx.Error
		}
		return &sensorOut, nil
	case "solenoid":
		solenoidOut := types.Solenoid{
			Model: types.Model{
				ID: idInt,
			},
		}
		tx := sql.Database.Where(&solenoidOut).Take(&solenoidOut)
		if tx.Error != nil {
			return nil, tx.Error
		}
		return &solenoidOut, nil
	case "stage":
		stageOut := types.Stage{
			Model: gorm.Model{
				ID: idInt,
			},
		}
		tx := sql.Database.Where(&stageOut).Take(&stageOut)
		if tx.Error != nil {
			return nil, tx.Error
		}
		return &stageOut, nil
	case "user":
		userOut := types.User{
			Model: gorm.Model{
				ID: idInt,
			},
		}
		tx := sql.Database.Where(&userOut).Take(&userOut)
		if tx.Error != nil {
			return nil, tx.Error
		}
		return &userOut, nil
	}
	return nil, errors.New("unknown node type")
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

func (r *queryResolver) IslandNodes(ctx context.Context) ([]*types.IslandNode, error) {
	var islandNodes []*types.IslandNode
	out := sql.Database.Find(&islandNodes)
	if out.Error != nil {
		return islandNodes, out.Error
	}
	return islandNodes, nil
}

func (r *queryResolver) LatestSensorData(ctx context.Context, queries []*types.SensorQuery) ([]float64, error) {
	sensorData := make([]float64, len(queries))
	for i, query := range queries {
		_, idInt, err := DecodeID(query.ID)
		if err != nil {
			return nil, err
		}
		data, err := systems.CurrentEngine.SensorsSystem.GetLatestSensorData(idInt)
		if err != nil {
			return sensorData, err
		}
		sensorData[i] = float64(data.SensorData)
		if query.Raw {
			sensorData[i] = float64(data.SensorDataRaw)
		}
	}
	return sensorData, nil
}

func (r *subscriptionResolver) NodeUpserted(ctx context.Context) (<-chan string, error) {
	in_channel := make(chan interface{})
	upsertBroadcast.Register(in_channel)
	out_channel := make(chan string)
	go func() {
		for {
			msg := <-in_channel
			out_channel <- fmt.Sprintf("%v", msg)
		}
	}()
	return out_channel, nil
}

func (r *subscriptionResolver) NodeDeleted(ctx context.Context) (<-chan string, error) {
	in_channel := make(chan interface{})
	deleteBroadcast.Register(in_channel)
	out_channel := make(chan string)
	go func() {
		for {
			msg := <-in_channel
			out_channel <- fmt.Sprintf("%v", msg)
		}
	}()
	return out_channel, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
