package models

import "context"

// AuthProvider describes the expected interface of an auth provider
// that handles creation of auth tokens and authorization of them.
type AuthProvider interface {
	AuthorizePluginToken(ctx context.Context) error
}
