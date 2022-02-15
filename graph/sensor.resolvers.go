package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/Liquid-Propulsion/mainland-server/graph/generated"
	"github.com/Liquid-Propulsion/mainland-server/systems"
	"github.com/Liquid-Propulsion/mainland-server/types"
)

func (r *sensorResolver) ID(ctx context.Context, obj *types.Sensor) (string, error) {
	return EncodeID("sensor", obj.Model.ID), nil
}

func (r *sensorResolver) RawValue(ctx context.Context, obj *types.Sensor) (float64, error) {
	data, err := systems.CurrentEngine.SensorsSystem.GetLatestSensorData(obj.ID)
	if err != nil {
		return 0.0, err
	}
	return float64(data.SensorDataRaw), nil
}

func (r *sensorResolver) Value(ctx context.Context, obj *types.Sensor) (float64, error) {
	data, err := systems.CurrentEngine.SensorsSystem.GetLatestSensorData(obj.ID)
	if err != nil {
		return 0.0, err
	}
	return float64(data.SensorData), nil
}

// Sensor returns generated.SensorResolver implementation.
func (r *Resolver) Sensor() generated.SensorResolver { return &sensorResolver{r} }

type sensorResolver struct{ *Resolver }
