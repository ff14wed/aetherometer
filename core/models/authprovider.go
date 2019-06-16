package models

import "context"

//go:generate counterfeiter . AuthProvider

// AuthProvider describes the expected interface of an auth provider
// that handles creation of auth tokens and authorization of them.
type AuthProvider interface {
	CreateAdminToken(ctx context.Context) (string, error)
	AddPlugin(ctx context.Context, pluginURL string) (string, error)
	RemovePlugin(ctx context.Context, apiToken string) (bool, error)

	AuthorizePluginToken(ctx context.Context) error
}
