package models

import (
	"context"
	"errors"
)

const SibylAPIVersion = "v0.0.0-beta"

type StreamRequestHandler func(pid int, data []byte) (resp string, err error)

// Resolver is a resolver for the queried data
type Resolver struct {
	sp      StoreProvider
	handler StreamRequestHandler
}

// NewResolver creates a new query resolver
// It takes the sp as an argument to use as a backing store for the queried
// data
func NewResolver(sp StoreProvider, streamRequestHandler StreamRequestHandler) *Resolver {
	return &Resolver{sp: sp, handler: streamRequestHandler}
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

// SendStreamRequest allows a client to send data upstream to configure the
// stream (the data source) at runtime.
func (r *mutationResolver) SendStreamRequest(ctx context.Context, req StreamRequest) (string, error) {
	if r.handler == nil {
		return "", errors.New("Request handler is missing")
	}
	return r.handler(req.StreamID, []byte(req.Data))
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) APIVersion(ctx context.Context) (string, error) {
	return SibylAPIVersion, nil
}
func (r *queryResolver) Streams(ctx context.Context) ([]Stream, error) {
	return r.sp.Streams()
}
func (r *queryResolver) Stream(ctx context.Context, streamID int) (Stream, error) {
	return r.sp.Stream(streamID)
}
func (r *queryResolver) Entity(ctx context.Context, streamID int, entityID uint64) (Entity, error) {
	return r.sp.Entity(streamID, entityID)
}

type subscriptionResolver struct{ *Resolver }

// StreamEvent returns an event channel that can be used for subscriptions to
// Stream events
func (r *subscriptionResolver) StreamEvent(ctx context.Context) (<-chan StreamEvent, error) {
	ses := r.sp.StreamEventSource()
	ch, id := ses.Subscribe()
	go func() {
		<-ctx.Done()
		ses.Unsubscribe(id)
	}()
	return ch, nil
}

// EntityEvent returns an event channel that can be used for subscriptions to
// Entity events
func (r *subscriptionResolver) EntityEvent(ctx context.Context) (<-chan EntityEvent, error) {
	ees := r.sp.EntityEventSource()
	ch, id := ees.Subscribe()
	go func() {
		<-ctx.Done()
		ees.Unsubscribe(id)
	}()
	return ch, nil
}
