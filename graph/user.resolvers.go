package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/Liquid-Propulsion/mainland-server/graph/generated"
	"github.com/Liquid-Propulsion/mainland-server/types"
)

func (r *userResolver) ID(ctx context.Context, obj *types.User) (string, error) {
	return EncodeID("user", obj.Model.ID), nil
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
