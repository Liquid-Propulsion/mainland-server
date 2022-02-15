package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/Liquid-Propulsion/mainland-server/graph/generated"
	"github.com/Liquid-Propulsion/mainland-server/systems"
	"github.com/Liquid-Propulsion/mainland-server/types"
)

func (r *islandNodeResolver) ID(ctx context.Context, obj *types.IslandNode) (string, error) {
	return EncodeID("node", obj.Model.ID), nil
}

func (r *islandNodeResolver) IsOnline(ctx context.Context, obj *types.IslandNode) (bool, error) {
	return systems.CurrentEngine.NodeSystem.NodeOnline(obj.ID), nil
}

// IslandNode returns generated.IslandNodeResolver implementation.
func (r *Resolver) IslandNode() generated.IslandNodeResolver { return &islandNodeResolver{r} }

type islandNodeResolver struct{ *Resolver }
