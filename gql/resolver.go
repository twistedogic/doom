package gql

import (
	"context"

	"github.com/twistedogic/doom/gql/generated"
	"github.com/twistedogic/doom/pkg/model"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Detail() generated.DetailResolver {
	return &detailResolver{r}
}
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

type detailResolver struct{ *Resolver }

func (r *detailResolver) ID(ctx context.Context, obj *model.Detail) (int, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) GetMatch(ctx context.Context, home *string, away *string) ([]*model.Match, error) {
	panic("not implemented")
}
func (r *queryResolver) GetDetail(ctx context.Context, matchID int) ([]*model.Detail, error) {
	panic("not implemented")
}
