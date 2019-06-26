package models

import (
	"context"
	"errors"
)

// AetherometerAPIVersion returns the current semantic version of the API. Generally,
// incremental additions to the API will be introduced with new patch versions.
// Minor breaking changes are introduced with new minor versions of the API.
// Major API changes and rewrites will be introduced with new major versions
// of the API
const AetherometerAPIVersion = "v0.1.0-beta"

// StreamRequestHandler defines the type of a client request handler that can
// be attached to the resolver.
type StreamRequestHandler func(streamID int, data []byte) (resp string, err error)

// Resolver is a resolver for the queried data
type Resolver struct {
	sp      StoreProvider
	auth    AuthProvider
	handler StreamRequestHandler
}

// NewResolver creates a new query resolver
// It takes the sp as an argument to use as a backing store for the queried
// data
func NewResolver(
	sp StoreProvider,
	auth AuthProvider,
	streamRequestHandler StreamRequestHandler,
) *Resolver {
	return &Resolver{sp: sp, auth: auth, handler: streamRequestHandler}
}

// Mutation allows graphql to handle mutation requests for the system
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
	if err := r.auth.AuthorizePluginToken(ctx); err != nil {
		return "", err
	}
	if r.handler == nil {
		return "", errors.New("Request handler is missing")
	}
	return r.handler(req.StreamID, []byte(req.Data))
}

// CreateAdminToken creates a token that is used to authorize the AddPlugin
// and RemovePlugin mutations. It can only be called once and will fail
// the subsequent attempts.
func (r *mutationResolver) CreateAdminToken(ctx context.Context) (string, error) {
	return r.auth.CreateAdminToken(ctx)
}

// AddPlugin registers the pluginURL with the API. It returns the apiToken
// that the plugin can use to authenticate with the API .
func (r *mutationResolver) AddPlugin(ctx context.Context, pluginURL string) (string, error) {
	return r.auth.AddPlugin(ctx, pluginURL)
}

// RemovePlugin revokes the access rights for the plugin associated with the API
// token.
func (r *mutationResolver) RemovePlugin(ctx context.Context, apiToken string) (bool, error) {
	return r.auth.RemovePlugin(ctx, apiToken)
}

type queryResolver struct{ *Resolver }

// APIVersion returns the version of the API (see AetherometerAPIVersion).
func (r *queryResolver) APIVersion(ctx context.Context) (string, error) {
	return AetherometerAPIVersion, nil
}

// Streams returns all of the streams known to the API
func (r *queryResolver) Streams(ctx context.Context) ([]Stream, error) {
	if err := r.auth.AuthorizePluginToken(ctx); err != nil {
		return nil, err
	}
	return r.sp.Streams()
}

// Stream returns the stream identified by streamID.
func (r *queryResolver) Stream(ctx context.Context, streamID int) (Stream, error) {
	if err := r.auth.AuthorizePluginToken(ctx); err != nil {
		return Stream{}, err
	}
	return r.sp.Stream(streamID)
}

// Entity returns the entity identified by entityID in the requested stream
// identified by streamID.
func (r *queryResolver) Entity(ctx context.Context, streamID int, entityID uint64) (Entity, error) {
	if err := r.auth.AuthorizePluginToken(ctx); err != nil {
		return Entity{}, err
	}
	return r.sp.Entity(streamID, entityID)
}

type subscriptionResolver struct{ *Resolver }

// StreamEvent returns an event channel that can be used for subscriptions to
// Stream events
func (r *subscriptionResolver) StreamEvent(ctx context.Context) (<-chan StreamEvent, error) {
	if err := r.auth.AuthorizePluginToken(ctx); err != nil {
		return nil, err
	}
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
	if err := r.auth.AuthorizePluginToken(ctx); err != nil {
		return nil, err
	}
	ees := r.sp.EntityEventSource()
	ch, id := ees.Subscribe()
	go func() {
		<-ctx.Done()
		ees.Unsubscribe(id)
	}()
	return ch, nil
}
