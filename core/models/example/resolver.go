package example

import (
	"context"

	"github.com/ff14wed/aetherometer/core/models"
)

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) SendStreamRequest(ctx context.Context, request models.StreamRequest) (string, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) APIVersion(ctx context.Context) (string, error) {
	panic("not implemented")
}
func (r *queryResolver) Streams(ctx context.Context) ([]models.Stream, error) {
	panic("not implemented")
}
func (r *queryResolver) Stream(ctx context.Context, streamID int) (models.Stream, error) {
	panic("not implemented")
}
func (r *queryResolver) Entity(ctx context.Context, streamID int, entityID uint64) (models.Entity, error) {
	panic("not implemented")
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) StreamEvent(ctx context.Context) (<-chan models.StreamEvent, error) {
	panic("not implemented")
}
func (r *subscriptionResolver) EntityEvent(ctx context.Context) (<-chan models.EntityEvent, error) {
	panic("not implemented")
}
