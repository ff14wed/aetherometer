package example

import (
	"context"

	"github.com/ff14wed/sibyl/backend/models"
)

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Streams(ctx context.Context) ([]models.Stream, error) {
	panic("not implemented")
}
func (r *queryResolver) Stream(ctx context.Context, streamID int) (models.Stream, error) {
	panic("not implemented")
}
func (r *queryResolver) Entity(ctx context.Context, streamID int, entityID int) (models.Entity, error) {
	panic("not implemented")
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) StreamEvents(ctx context.Context) (<-chan models.StreamEventsPayload, error) {
	panic("not implemented")
}
func (r *subscriptionResolver) EntityEvents(ctx context.Context) (<-chan models.EntityEventsPayload, error) {
	panic("not implemented")
}
