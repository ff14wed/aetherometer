package handlers

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/99designs/gqlgen/handler"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/ff14wed/aetherometer/core/config"
	"github.com/rs/cors"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

type ctxKey struct {
	name string
}

var authCtxKey = &ctxKey{name: "auth"}

var AuthError = errors.New(
	"Unauthorized: The request lacks valid authentication credentials for the target resource.",
)

// InitPayloadGetter provides an alternate method of retrieving credentials
// for when the token cannot be sent via HTTP headers (basically the case of
// websockets)
type InitPayloadGetter func(context.Context) handler.InitPayload

// Auth provides a middleware handler that handles cross-origin
// requests and pulling of authentication from incoming requests.
// It also handles the registration and deregistration of plugins
// and updates its authorizer methods accordingly.
type Auth struct {
	otpUsed  int32
	adminOTP string
	adminID  atomic.Value

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

		adminOTP:    c.AdminOTP,
		disableAuth: c.DisableAuth,

		allowedPlugins: make(map[string]string),
		allowedOrigins: make(map[string]int),

		initPayloadGetter: initPayloadGetter,

		logger: l.Named("auth-handler"),
	}

	invalidID := xid.New().String()
	a.adminID.Store(invalidID)

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
	if strings.Contains(origin, "http://localhost") {
		return true
	}

	a.allowedLock.RLock()
	count := a.allowedOrigins[origin]
	a.allowedLock.RUnlock()
	return count > 0
}

// CreateAdminToken returns a persistent token that is authorized to not only
// the same requests that plugins can make, but also register or unregister
// plugins.
func (a *Auth) CreateAdminToken(ctx context.Context) (string, error) {
	authString, _ := ctx.Value(authCtxKey).(string)
	if authString != a.adminOTP {
		return "", AuthError
	}

	swapped := atomic.CompareAndSwapInt32(&a.otpUsed, 0, 1)
	if !swapped {
		return "", errors.New("No more admin tokens can be created.")
	}

	adminID := xid.New().String()
	a.adminID.Store(adminID)

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"adminID": adminID,
	})

	adminToken, err := t.SignedString(a.privKey)
	if err != nil {
		return "", err
	}

	return adminToken, nil
}

// AddPlugin registers the provided plugin URL in the system and returns the
// token that the plugin can use to make requests. If the plugin URL cannot be
// parsed, it fails to create a token and returns an error.
func (a *Auth) AddPlugin(ctx context.Context, pluginURL string) (string, error) {
	if err := a.AuthorizeAdminToken(ctx); err != nil {
		return "", err
	}

	parsedURL, err := url.Parse(pluginURL)
	if err != nil {
		return "", err
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", errors.New("Could not parse plugin URL.")
	}

	origin := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	guid := xid.New().String()

	a.allowedLock.Lock()
	a.allowedPlugins[guid] = origin
	a.allowedOrigins[origin] += 1
	a.allowedLock.Unlock()

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"url": pluginURL,
		"id":  guid,
	})

	apiToken, err := t.SignedString(a.privKey)
	if err != nil {
		return "", err
	}

	a.logger.Debug("Added Plugin", zap.String("url", pluginURL), zap.String("id", guid))
	return apiToken, nil
}

// RemovePlugin unregisters the plugin corresponding to this specific apiToken.
// It will return success even if the apiToken is no longer registered,
// provided that this token is still a valid token.
// If all plugins corresponding to a particular origin are removed, this origin
// will no longer be allowed in the CORS mechanism.
func (a *Auth) RemovePlugin(ctx context.Context, apiToken string) (bool, error) {
	if err := a.AuthorizeAdminToken(ctx); err != nil {
		return false, err
	}

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

	return false, AuthError
}

// AuthorizeAdminToken checks to make sure that the provided admin token
// is valid and is authorized to make admin-level requests.
func (a *Auth) AuthorizeAdminToken(ctx context.Context) error {
	claims, err := a.extractClaimsFromCtx(ctx)
	if err != nil {
		return err
	}
	if adminID, found := claims["adminID"]; found {
		correctID := a.adminID.Load().(string)
		if adminID.(string) != correctID {
			return AuthError
		}
		return nil
	}
	return AuthError
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

	if adminID, found := claims["adminID"]; found {
		correctID := a.adminID.Load().(string)
		if adminID.(string) != correctID {
			return AuthError
		}
		return nil
	}

	if pluginID, found := claims["id"]; found {
		a.allowedLock.RLock()
		defer a.allowedLock.RUnlock()
		if _, allowed := a.allowedPlugins[pluginID.(string)]; !allowed {
			return AuthError
		}
		return nil
	}
	return AuthError
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
		return nil, AuthError
	}
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return &a.privKey.PublicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, AuthError
	}
	return claims, nil
}
