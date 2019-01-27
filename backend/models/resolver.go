package models

import (
	"context"
)

const SibylAPIVersion = "v0.0.0-beta"

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

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query allows graphql to resolve queries made on the system
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// Subscription allows graphql to resolve subscriptions added to the system
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) SendAdapterRequest(ctx context.Context, req AdapterRequest) (string, error) {
	return r.db.SendAdapterRequest(req)
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) APIVersion(ctx context.Context) (string, error) {
	return SibylAPIVersion, nil
}
func (r *queryResolver) Streams(ctx context.Context) ([]Stream, error) {
	return r.db.Streams(), nil
}
func (r *queryResolver) Stream(ctx context.Context, streamID int) (Stream, error) {
	return r.db.Stream(streamID)
}
func (r *queryResolver) Entity(ctx context.Context, streamID int, entityID uint64) (Entity, error) {
	return r.db.Entity(streamID, entityID)
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) StreamEvent(ctx context.Context) (<-chan StreamEvent, error) {
	return r.db.StreamEvent(ctx)
}
func (r *subscriptionResolver) EntityEvent(ctx context.Context) (<-chan EntityEvent, error) {
	return r.db.EntityEvent(ctx)
}
