package models

import (
	"context"
)

// Resolver is a resolver for the queried data
type Resolver struct {
	db *DB
}

// NewResolver creates a new query resolver
// It takes the db as an argument to use as a backing store for the queried
// data
func NewResolver(db *DB) *Resolver {
	return &Resolver{db: db}
}

// Query allows graphql to resolve queries made on the system
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// Subscription allows graphql to resolve subscriptions added to the system
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Streams(ctx context.Context) ([]Stream, error) {
	return r.db.Streams(), nil
}
func (r *queryResolver) Stream(ctx context.Context, streamID int) (Stream, error) {
	return r.db.Stream(streamID)
}
func (r *queryResolver) Entity(ctx context.Context, streamID int, entityID int) (Entity, error) {
	return r.db.Entity(streamID, entityID)
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) StreamEvents(ctx context.Context) (<-chan StreamEventsPayload, error) {
	return r.db.StreamEvents(ctx)
}
func (r *subscriptionResolver) EntityEvents(ctx context.Context) (<-chan EntityEventsPayload, error) {
	return r.db.EntityEvents(ctx)
}
