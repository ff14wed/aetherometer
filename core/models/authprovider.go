package models

import "context"

// AuthProvider describes the expected interface of an auth provider
// that handles creation of auth tokens and authorization of them.
type AuthProvider interface {
	AddPlugin(pluginURL string) (string, error)
	RemovePlugin(apiToken string) (bool, error)

	AuthorizePluginToken(ctx context.Context) error
}
