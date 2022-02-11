package handlers

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/ff14wed/aetherometer/core/config"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

type ctxKey struct {
	name string
}

var authCtxKey = &ctxKey{name: "auth"}

var ErrAuth = errors.New(
	"unauthorized error: the request lacks valid authentication credentials for the target resource",
)

// InitPayloadGetter provides an alternate method of retrieving credentials
// for when the token cannot be sent via HTTP headers (basically the case of
// websockets)
type InitPayloadGetter func(context.Context) transport.InitPayload

type PluginInfo struct {
	PluginID  string
	PluginURL string
	APIToken  string
}
type authConfig struct {
	disableAuth bool

	// plugins maps -> plugin name -> (pluginID, apiToken)
	plugins map[string]PluginInfo
	// allowedPluginIDs caches the set of allowed pluginIDs
	allowedPluginIDs map[string]struct{}
	// allowedOrigins caches the set of allowed origins
	allowedOrigins map[string]struct{}
}

// Auth provides a middleware handler that handles cross-origin
// requests and pulling of authentication from incoming requests.
// It also handles the registration and deregistration of plugins
// and updates its authorizer methods accordingly.
type Auth struct {
	cp *config.Provider

	authConfig     authConfig
	authConfigLock sync.RWMutex

	cors *cors.Cors

	privKey *rsa.PrivateKey

	initPayloadGetter InitPayloadGetter

	logger *zap.Logger
}

// NewAuth creates a new instance of Auth. InitPayloadGetter provides a
// way to grab credentials from the context when the connection is made
// via Websockets.
func NewAuth(cp *config.Provider, initPayloadGetter InitPayloadGetter, l *zap.Logger) (*Auth, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	a := &Auth{
		privKey: key,

		cp: cp,

		authConfig: authConfig{
			disableAuth: false,

			plugins:          make(map[string]PluginInfo),
			allowedPluginIDs: make(map[string]struct{}),
			allowedOrigins:   make(map[string]struct{}),
		},

		initPayloadGetter: initPayloadGetter,

		logger: l.Named("auth-handler"),
	}

	a.RefreshConfig()

	a.cors = cors.New(cors.Options{
		AllowOriginFunc:  a.AllowOriginFunc,
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Apollo-Tracing"},
	})

	return a, nil
}

// RefreshConfig will read the configuration from the config provider and
// add and remove plugins as necessary. Any preexisting plugins will retain
// the same ID and apiToken
func (a *Auth) RefreshConfig() error {
	a.authConfigLock.Lock()
	defer a.authConfigLock.Unlock()

	cfg := a.cp.Config()

	updatedPlugins := make(map[string]PluginInfo)
	updatedAllowedPluginIDs := make(map[string]struct{})
	updatedAllowedOrigins := make(map[string]struct{})

	for pluginName, pluginURL := range cfg.Plugins {
		origin, err := parseOrigin(pluginURL)
		if err != nil {
			return err
		}

		if info, ok := a.authConfig.plugins[pluginName]; ok {
			updatedPlugins[pluginName] = info
			updatedAllowedPluginIDs[info.PluginID] = struct{}{}
		} else {
			newInfo, err := a.generatePluginInfo(pluginURL)
			if err != nil {
				return err
			}
			updatedPlugins[pluginName] = newInfo
			updatedAllowedPluginIDs[newInfo.PluginID] = struct{}{}
		}

		updatedAllowedOrigins[origin] = struct{}{}
	}

	a.authConfig.disableAuth = cfg.DisableAuth
	a.authConfig.plugins = updatedPlugins
	a.authConfig.allowedPluginIDs = updatedAllowedPluginIDs
	a.authConfig.allowedOrigins = updatedAllowedOrigins

	return nil
}

// Handler returns a handler that serves cross-origin requests and pulls
// authentication data from incoming requests.
func (a *Auth) Handler(next http.Handler) http.Handler {
	return a.cors.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authString, err := request.AuthorizationHeaderExtractor.ExtractToken(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), authCtxKey, authString)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}))
}

// AllowOriginFunc returns true only if the origin is referenced by any URL
// of the registered plugins.
func (a *Auth) AllowOriginFunc(origin string) bool {
	if strings.Contains(origin, "file://") {
		return true
	}
	if strings.Contains(origin, "app://") {
		return true
	}
	if strings.Contains(origin, "http://localhost") {
		return true
	}

	a.authConfigLock.RLock()
	defer a.authConfigLock.RUnlock()
	_, ok := a.authConfig.allowedOrigins[origin]
	return ok
}

func parseOrigin(pluginURL string) (string, error) {
	parsedURL, err := url.Parse(pluginURL)
	if err != nil {
		return "", err
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", errors.New("could not parse plugin URL")
	}

	return fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host), nil
}

// generatePluginInfo registers the provided plugin URL in the system and
// returns the token that the plugin can use to make requests. If the plugin URL
// cannot be parsed, it fails to create a token and returns an error.
func (a *Auth) generatePluginInfo(pluginURL string) (PluginInfo, error) {
	pluginID, err := generateRandomString(32)
	if err != nil {
		return PluginInfo{}, errors.New("system error: could not generate plugin ID")
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"url": pluginURL,
		"id":  pluginID,
	})
	apiToken, err := t.SignedString(a.privKey)
	if err != nil {
		return PluginInfo{}, err
	}

	a.logger.Debug("Generated new plugin auth", zap.String("url", pluginURL), zap.String("id", pluginID))

	return PluginInfo{
		PluginID:  pluginID,
		PluginURL: pluginURL,
		APIToken:  apiToken,
	}, nil
}

// AuthorizePluginToken checks to make sure that the provided plugin token
// is valid and is authorized to make plugin-level requests.
func (a *Auth) AuthorizePluginToken(ctx context.Context) error {
	a.authConfigLock.RLock()
	defer a.authConfigLock.RUnlock()

	if a.authConfig.disableAuth {
		return nil
	}

	claims, err := a.extractClaimsFromCtx(ctx)
	if err != nil {
		return err
	}

	pluginID, found := claims["id"]
	if !found {
		return ErrAuth
	}
	if _, allowed := a.authConfig.allowedPluginIDs[pluginID.(string)]; !allowed {
		return ErrAuth
	}
	return nil
}

// GetRegisteredPlugins returns a copy of the map containing all the plugins
// and their API tokens.
func (a *Auth) GetRegisteredPlugins() map[string]PluginInfo {
	a.authConfigLock.RLock()
	defer a.authConfigLock.RUnlock()

	plugins := make(map[string]PluginInfo)

	for name, info := range a.authConfig.plugins {
		plugins[name] = info
	}

	return plugins
}

func (a *Auth) extractClaimsFromCtx(ctx context.Context) (jwt.MapClaims, error) {
	authString, found := ctx.Value(authCtxKey).(string)
	if !found && a.initPayloadGetter != nil {
		authString = a.initPayloadGetter(ctx).Authorization()
	}
	return a.extractClaimsFromToken(authString)
}

func (a *Auth) extractClaimsFromToken(tokenString string) (jwt.MapClaims, error) {
	if len(tokenString) == 0 {
		return nil, ErrAuth
	}
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return &a.privKey.PublicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrAuth
	}
	return claims, nil
}

func generateRandomString(length int) (string, error) {
	const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	s := make([]byte, length)
	for i := 0; i < length; i++ {
		c, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		s[i] = chars[c.Int64()]
	}
	return string(s), nil
}
