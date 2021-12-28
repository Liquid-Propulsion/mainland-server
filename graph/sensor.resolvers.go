package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/Liquid-Propulsion/mainland-server/graph/generated"
	"github.com/Liquid-Propulsion/mainland-server/types"
)

func (r *sensorResolver) ID(ctx context.Context, obj *types.Sensor) (string, error) {
	return EncodeID("sensor", obj.Model.ID), nil
}

func (r *sensorResolver) NodeID(ctx context.Context, obj *types.Sensor) (int, error) {
	return int(obj.NodeID), nil
}

func (r *sensorResolver) SensorID(ctx context.Context, obj *types.Sensor) (int, error) {
	return int(obj.SensorID), nil
}

// Sensor returns generated.SensorResolver implementation.
func (r *Resolver) Sensor() generated.SensorResolver { return &sensorResolver{r} }

type sensorResolver struct{ *Resolver }
