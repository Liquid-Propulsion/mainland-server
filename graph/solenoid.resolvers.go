package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/Liquid-Propulsion/mainland-server/graph/generated"
	"github.com/Liquid-Propulsion/mainland-server/types"
)

func (r *solenoidResolver) ID(ctx context.Context, obj *types.Solenoid) (string, error) {
	return EncodeID("solenoid", obj.Model.ID), nil
}

// Solenoid returns generated.SolenoidResolver implementation.
func (r *Resolver) Solenoid() generated.SolenoidResolver { return &solenoidResolver{r} }

type solenoidResolver struct{ *Resolver }
