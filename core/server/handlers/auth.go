package handlers

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"math"
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

// Auth provides a middleware handler that handles cross-origin
// requests and pulling of authentication from incoming requests.
// It also handles the registration and deregistration of plugins
// and updates its authorizer methods accordingly.
type Auth struct {
	disableAuth bool

	allowedPlugins map[string]string
	allowedOrigins map[string]int
	allowedLock    sync.RWMutex

	cors *cors.Cors

	privKey *rsa.PrivateKey

	initPayloadGetter InitPayloadGetter

	logger *zap.Logger
}

// NewAuth creates a new instance of Auth. InitPayloadGetter provides a
// way to grab credentials from the context when the connection is made
// via Websockets.
func NewAuth(c config.Config, initPayloadGetter InitPayloadGetter, l *zap.Logger) (*Auth, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	a := &Auth{
		privKey: key,

		disableAuth: c.DisableAuth,

		allowedPlugins: make(map[string]string),
		allowedOrigins: make(map[string]int),

		initPayloadGetter: initPayloadGetter,

		logger: l.Named("auth-handler"),
	}

	if len(c.Plugins) > 0 {
		for _, plugin := range c.Plugins {
			a.allowedOrigins[plugin] = math.MaxInt32
		}
	}

	a.cors = cors.New(cors.Options{
		AllowOriginFunc:  a.AllowOriginFunc,
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Apollo-Tracing"},
	})

	return a, nil
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

	a.allowedLock.RLock()
	count := a.allowedOrigins[origin]
	a.allowedLock.RUnlock()
	return count > 0
}

// AddPlugin registers the provided plugin URL in the system and returns the
// token that the plugin can use to make requests. If the plugin URL cannot be
// parsed, it fails to create a token and returns an error.
func (a *Auth) AddPlugin(pluginURL string) (string, error) {
	parsedURL, err := url.Parse(pluginURL)
	if err != nil {
		return "", err
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", errors.New("could not parse plugin URL")
	}

	origin := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	pluginID, err := generateRandomString(32)
	if err != nil {
		return "", errors.New("system error: could not generate plugin ID")
	}

	a.allowedLock.Lock()
	a.allowedPlugins[pluginID] = origin
	a.allowedOrigins[origin] += 1
	a.allowedLock.Unlock()

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"url": pluginURL,
		"id":  pluginID,
	})

	apiToken, err := t.SignedString(a.privKey)
	if err != nil {
		return "", err
	}

	a.logger.Debug("Added Plugin", zap.String("url", pluginURL), zap.String("id", pluginID))
	return apiToken, nil
}

// RemovePlugin unregisters the plugin corresponding to this specific apiToken.
// It will return success even if the apiToken is no longer registered,
// provided that this token is still a valid token.
// If all plugins corresponding to a particular origin are removed, this origin
// will no longer be allowed in the CORS mechanism.
func (a *Auth) RemovePlugin(apiToken string) (bool, error) {
	claims, err := a.extractClaimsFromToken(apiToken)
	if err != nil {
		return false, err
	}

	if pluginID, found := claims["id"]; found {
		id := pluginID.(string)
		a.allowedLock.Lock()
		if origin, found := a.allowedPlugins[id]; found {
			a.allowedOrigins[origin] -= 1
		}
		delete(a.allowedPlugins, id)
		a.allowedLock.Unlock()
		a.logger.Debug("Deleted Plugin", zap.String("id", id))
		return true, nil
	}

	return false, ErrAuth
}

// AuthorizePluginToken checks to make sure that the provided plugin token
// is valid and is authorized to make plugin-level requests.
func (a *Auth) AuthorizePluginToken(ctx context.Context) error {
	if a.disableAuth {
		return nil
	}

	claims, err := a.extractClaimsFromCtx(ctx)
	if err != nil {
		return err
	}

	if pluginID, found := claims["id"]; found {
		a.allowedLock.RLock()
		defer a.allowedLock.RUnlock()
		if _, allowed := a.allowedPlugins[pluginID.(string)]; !allowed {
			return ErrAuth
		}
		return nil
	}
	return ErrAuth
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
