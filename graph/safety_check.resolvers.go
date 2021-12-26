package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/Liquid-Propulsion/mainland-server/graph/generated"
	"github.com/Liquid-Propulsion/mainland-server/types"
)

func (r *safetyCheckResolver) ID(ctx context.Context, obj *types.SafetyCheck) (string, error) {
	return EncodeID("safety_check", *obj.ID), nil
}

// SafetyCheck returns generated.SafetyCheckResolver implementation.
func (r *Resolver) SafetyCheck() generated.SafetyCheckResolver { return &safetyCheckResolver{r} }

type safetyCheckResolver struct{ *Resolver }
