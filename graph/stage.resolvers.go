package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/Liquid-Propulsion/mainland-server/graph/generated"
	"github.com/Liquid-Propulsion/mainland-server/types"
)

func (r *stageResolver) ID(ctx context.Context, obj *types.Stage) (string, error) {
	return EncodeID("stage", obj.Model.ID), nil
}

func (r *stageResolver) CanID(ctx context.Context, obj *types.Stage) (int, error) {
	return int(obj.CANID), nil
}

func (r *stageResolver) Duration(ctx context.Context, obj *types.Stage) (string, error) {
	return obj.Duration.String(), nil
}

// Stage returns generated.StageResolver implementation.
func (r *Resolver) Stage() generated.StageResolver { return &stageResolver{r} }

type stageResolver struct{ *Resolver }
