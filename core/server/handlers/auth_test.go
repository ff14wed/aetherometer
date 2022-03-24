package handlers_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sync"

	"github.com/99designs/gqlgen/graphql/handler/transport"
	"go.uber.org/zap"

	"github.com/ff14wed/aetherometer/core/config"
	"github.com/ff14wed/aetherometer/core/server/handlers"
	"github.com/ff14wed/aetherometer/core/testhelpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Auth", func() {
	var (
		cp         *config.Provider
		configFile string

		auth *handlers.Auth

		logBuf *testhelpers.LogBuffer
		once   sync.Once
	)

	BeforeEach(func() {
		once.Do(func() {
			logBuf = new(testhelpers.LogBuffer)
			err := zap.RegisterSink("authhandlertest", func(*url.URL) (zap.Sink, error) {
				return logBuf, nil
			})
			Expect(err).ToNot(HaveOccurred())
		})
		logBuf.Reset()

		zapCfg := zap.NewDevelopmentConfig()
		zapCfg.OutputPaths = []string{"authhandlertest://"}
		logger, err := zapCfg.Build()
		Expect(err).ToNot(HaveOccurred())

		f, err := ioutil.TempFile("", "authtest")
		Expect(err).ToNot(HaveOccurred())
		configFile = f.Name()

		cfg := config.Config{}
		cp = config.NewProvider(configFile, cfg, logger)

		auth, err = handlers.NewAuth(cp, logger)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		_ = os.Remove(configFile)
	})

	Describe("adding a plugin", func() {
		It("authorizes the plugin to access the API", func() {
			oldPlugins := auth.GetRegisteredPlugins()

			Expect(cp.AddPlugin("Foo Plugin", "https://example.com/foo/plugin")).To(Succeed())
			Expect(auth.RefreshConfig()).To(Succeed())

			plugins := auth.GetRegisteredPlugins()
			Expect(plugins).ToNot(Equal(oldPlugins))
			pluginInfo := plugins["Foo Plugin"]

			ctx := handlers.ContextWithToken(pluginInfo.APIToken)
			Expect(auth.AuthorizePluginToken(ctx)).To(Succeed())
			Expect(auth.AllowOriginFunc("https://example.com")).To(BeTrue())

			Eventually(logBuf).Should(gbytes.Say("DEBUG"))
			Eventually(logBuf).Should(gbytes.Say("auth-handler"))
			Eventually(logBuf).Should(gbytes.Say("Generated new plugin auth"))
			Eventually(logBuf).Should(gbytes.Say("https://example.com.*id"))
		})

		It("returns an error if the plugin URL is invalid", func() {
			Expect(cp.AddPlugin("Bad Plugin", ":bad-url.com")).To(Succeed())
			Expect(auth.RefreshConfig()).To(BeAssignableToTypeOf(new(url.Error)))

			plugins := auth.GetRegisteredPlugins()
			Expect(plugins).To(BeEmpty())
		})

		It("returns an error if the plugin URL is missing the scheme", func() {
			Expect(cp.AddPlugin("Bad Plugin", "example.com/foo/plugin")).To(Succeed())
			Expect(auth.RefreshConfig()).To(MatchError("could not parse plugin URL"))

			plugins := auth.GetRegisteredPlugins()
			Expect(plugins).To(BeEmpty())
		})

		It("does not change the registered plugins list if a bad plugin URL is added", func() {
			Expect(cp.AddPlugin("Foo Plugin", "https://example.com/foo/plugin")).To(Succeed())
			Expect(auth.RefreshConfig()).To(Succeed())
			plugins := auth.GetRegisteredPlugins()
			Expect(plugins).To(HaveKey("Foo Plugin"))

			Expect(cp.AddPlugin("Bad Plugin", "example.com/foo/plugin")).To(Succeed())
			Expect(auth.RefreshConfig()).To(MatchError("could not parse plugin URL"))

			plugins = auth.GetRegisteredPlugins()
			Expect(plugins).To(HaveKey("Foo Plugin"))
			Expect(plugins).ToNot(HaveKey("Bad Plugin"))

			pluginInfo := plugins["Foo Plugin"]
			ctx := handlers.ContextWithToken(pluginInfo.APIToken)
			Expect(auth.AuthorizePluginToken(ctx)).To(Succeed())
			Expect(auth.AllowOriginFunc("https://example.com")).To(BeTrue())
		})
	})

	Describe("removing a plugin", func() {
		const pluginName = "Foo Plugin"
		var apiToken string

		BeforeEach(func() {
			Expect(cp.AddPlugin(pluginName, "https://example.com/foo/plugin")).To(Succeed())
			Expect(auth.RefreshConfig()).To(Succeed())

			plugins := auth.GetRegisteredPlugins()
			pluginInfo := plugins[pluginName]

			apiToken = pluginInfo.APIToken
			ctx := handlers.ContextWithToken(apiToken)
			Expect(auth.AuthorizePluginToken(ctx)).To(Succeed())
			Expect(auth.AllowOriginFunc("https://example.com")).To(BeTrue())
		})

		It("revokes authorization for the plugin to access the API", func() {
			Expect(cp.RemovePlugin(pluginName)).To(Succeed())
			Expect(auth.RefreshConfig()).To(Succeed())

			ctx := handlers.ContextWithToken(apiToken)
			Expect(auth.AuthorizePluginToken(ctx)).To(MatchError(handlers.ErrAuth))
			Expect(auth.AllowOriginFunc("https://example.com")).To(BeFalse())
		})

		Context("when adding a second plugin with the same origin", func() {
			const altPluginName string = "Bar Plugin"
			var altAPIToken string

			BeforeEach(func() {
				Expect(cp.AddPlugin(altPluginName, "https://example.com/bar/plugin")).To(Succeed())
				Expect(auth.RefreshConfig()).To(Succeed())

				plugins := auth.GetRegisteredPlugins()
				pluginInfo := plugins[altPluginName]

				altAPIToken = pluginInfo.APIToken
			})

			It("continues to allow the origin when removing a single plugin", func() {
				Expect(cp.RemovePlugin(pluginName)).To(Succeed())
				Expect(auth.RefreshConfig()).To(Succeed())

				ctx := handlers.ContextWithToken(apiToken)
				Expect(auth.AuthorizePluginToken(ctx)).To(MatchError(handlers.ErrAuth))

				ctx = handlers.ContextWithToken(altAPIToken)
				Expect(auth.AuthorizePluginToken(ctx)).To(Succeed())
				Expect(auth.AllowOriginFunc("https://example.com")).To(BeTrue())
			})

			It("no longer allows the origin when removing both plugins", func() {
				Expect(cp.RemovePlugin(pluginName)).To(Succeed())
				Expect(cp.RemovePlugin(altPluginName)).To(Succeed())
				Expect(auth.RefreshConfig()).To(Succeed())

				ctx := handlers.ContextWithToken(apiToken)
				Expect(auth.AuthorizePluginToken(ctx)).To(MatchError(handlers.ErrAuth))

				ctx = handlers.ContextWithToken(altAPIToken)
				Expect(auth.AuthorizePluginToken(ctx)).To(MatchError(handlers.ErrAuth))
				Expect(auth.AllowOriginFunc("https://example.com")).To(BeFalse())
			})
		})
	})

	Describe("AuthorizePluginToken", func() {
		Context("when auth is disabled", func() {
			BeforeEach(func() {
				Expect(cp.MutateConfig(func(cfg config.Config) (config.Config, error) {
					cfg.DisableAuth = true
					return cfg, nil
				})).To(Succeed())
				Expect(auth.RefreshConfig()).To(Succeed())
			})

			It("authorizes requests with no token to access the API", func() {
				Expect(auth.AuthorizePluginToken(context.Background())).To(Succeed())
			})
		})

		Context("when the authorization token is sent via websockets", func() {
			It("successfully validates the token", func() {
				Expect(cp.AddPlugin("Foo Plugin", "https://example.com/foo/plugin")).To(Succeed())
				Expect(auth.RefreshConfig()).To(Succeed())

				plugins := auth.GetRegisteredPlugins()
				pluginInfo := plugins["Foo Plugin"]

				initPayload := transport.InitPayload{"Authorization": pluginInfo.APIToken}
				background := context.Background()
				ctx, err := auth.WebsocketInitFunc(background, initPayload)
				Expect(err).ToNot(HaveOccurred())

				Expect(auth.AuthorizePluginToken(ctx)).To(Succeed())
			})

			It("rejects invalid tokens", func() {
				initPayload := transport.InitPayload{"Authorization": "invalid-token"}
				background := context.Background()
				ctx, err := auth.WebsocketInitFunc(background, initPayload)
				Expect(err).ToNot(HaveOccurred())
				Expect(auth.AuthorizePluginToken(ctx)).To(MatchError(handlers.ErrAuth))
			})
		})

		Context("when the provided api token is no longer valid", func() {
			var (
				altAPIToken   string
				altConfigFile string
			)

			BeforeEach(func() {

				f, err := ioutil.TempFile("", "authtest")
				Expect(err).ToNot(HaveOccurred())
				altConfigFile = f.Name()

				cfg := config.Config{}
				altCP := config.NewProvider(altConfigFile, cfg, zap.NewNop())

				altAuth, err := handlers.NewAuth(altCP, zap.NewNop())
				Expect(err).ToNot(HaveOccurred())

				Expect(altCP.AddPlugin("Foo Plugin", "https://example.com/foo/plugin")).To(Succeed())
				Expect(altAuth.RefreshConfig()).To(Succeed())

				plugins := altAuth.GetRegisteredPlugins()
				pluginInfo := plugins["Foo Plugin"]
				altAPIToken = pluginInfo.APIToken
			})

			AfterEach(func() {
				_ = os.Remove(altConfigFile)
			})

			It("rejects the api token", func() {
				ctx := handlers.ContextWithToken(altAPIToken)
				Expect(auth.AuthorizePluginToken(ctx)).To(MatchError(handlers.ErrAuth))
			})
		})
	})

	Describe("Handler", func() {
		It("adds the auth token to the request context", func() {
			req, err := http.NewRequest("POST", "/foo", nil)
			Expect(err).ToNot(HaveOccurred())

			Expect(cp.AddPlugin("Foo Plugin", "https://example.com/foo/plugin")).To(Succeed())
			Expect(auth.RefreshConfig()).To(Succeed())

			plugins := auth.GetRegisteredPlugins()
			pluginInfo := plugins["Foo Plugin"]
			apiToken := pluginInfo.APIToken

			req.Header.Set("Authorization", "Bearer "+apiToken)

			var receivedCtx context.Context
			authHandler := auth.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedCtx = r.Context()
			}))

			rw := httptest.NewRecorder()
			authHandler.ServeHTTP(rw, req)

			Expect(receivedCtx).ToNot(BeNil())

			Expect(auth.AuthorizePluginToken(receivedCtx)).To(Succeed())
		})

		It("calls the provided next handler even if there is no Authorization header", func() {
			req, err := http.NewRequest("POST", "/foo", nil)
			Expect(err).ToNot(HaveOccurred())

			var receivedCtx context.Context
			authHandler := auth.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedCtx = r.Context()
			}))

			rw := httptest.NewRecorder()
			authHandler.ServeHTTP(rw, req)

			Expect(receivedCtx).ToNot(BeNil())
			Expect(auth.AuthorizePluginToken(receivedCtx)).To(MatchError(handlers.ErrAuth))
		})

		It("does not allow unknown origins in the preflight request", func() {
			req, err := http.NewRequest("OPTIONS", "/foo", nil)
			Expect(err).ToNot(HaveOccurred())
			req.Header.Set("Access-Control-Request-Method", "POST")
			req.Header.Set("Access-Control-Request-Headers", "Authorization,Content-Type,X-Apollo-Tracing")
			req.Header.Set("Origin", "https://example.com")

			authHandler := auth.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(502)
			}))

			rw := httptest.NewRecorder()
			authHandler.ServeHTTP(rw, req)
			Expect(rw.Header().Get("Access-Control-Allow-Origin")).To(BeEmpty())
			Expect(rw.Header().Get("Access-Control-Allow-Headers")).To(BeEmpty())
			Expect(rw.Header().Get("Access-Control-Allow-Methods")).To(BeEmpty())
			Expect(rw.Header().Get("Access-Control-Allow-Credentials")).To(BeEmpty())
		})

		Context("when the preflight request is for a plugin that has been authorized", func() {
			BeforeEach(func() {
				Expect(cp.AddPlugin("Foo Plugin", "https://example.com/foo/plugin")).To(Succeed())
				Expect(auth.RefreshConfig()).To(Succeed())
			})

			It("allows the origin in the preflight request", func() {
				req, err := http.NewRequest("OPTIONS", "/foo", nil)
				Expect(err).ToNot(HaveOccurred())
				req.Header.Set("Access-Control-Request-Method", "POST")
				req.Header.Set("Access-Control-Request-Headers", "Authorization,Content-Type,X-Apollo-Tracing")
				req.Header.Set("Origin", "https://example.com")

				authHandler := auth.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(502)
				}))

				rw := httptest.NewRecorder()
				authHandler.ServeHTTP(rw, req)
				Expect(rw.Header().Get("Access-Control-Allow-Origin")).To(Equal("https://example.com"))
				Expect(rw.Header().Get("Access-Control-Allow-Methods")).To(Equal("POST"))
				Expect(rw.Header().Get("Access-Control-Allow-Headers")).To(Equal("Authorization, Content-Type, X-Apollo-Tracing"))
				Expect(rw.Header().Get("Access-Control-Allow-Credentials")).To(Equal("true"))
			})
		})

		Context("when a local plugin provides a correct local token", func() {
			BeforeEach(func() {
				Expect(cp.MutateConfig(func(cfg config.Config) (config.Config, error) {
					cfg.LocalToken = "some-local-token"
					return cfg, nil
				})).To(Succeed())
				Expect(auth.RefreshConfig()).To(Succeed())
			})

			It("allows local plugins to use the local token to authenticate", func() {
				req, err := http.NewRequest("POST", "/foo", nil)
				Expect(err).ToNot(HaveOccurred())

				req.Header.Set("Origin", "app://bar")
				req.Header.Set("Authorization", "Bearer some-local-token")

				var receivedCtx context.Context
				authHandler := auth.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					receivedCtx = r.Context()
				}))

				rw := httptest.NewRecorder()
				authHandler.ServeHTTP(rw, req)

				Expect(receivedCtx).ToNot(BeNil())

				Expect(auth.AuthorizePluginToken(receivedCtx)).To(Succeed())
			})

			It("does not allow non-local plugins to use the local token to authenticate", func() {
				req, err := http.NewRequest("POST", "/foo", nil)
				Expect(err).ToNot(HaveOccurred())

				req.Header.Set("Authorization", "Bearer some-local-token")

				var receivedCtx context.Context
				authHandler := auth.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					receivedCtx = r.Context()
				}))

				rw := httptest.NewRecorder()
				authHandler.ServeHTTP(rw, req)

				Expect(receivedCtx).ToNot(BeNil())

				Expect(auth.AuthorizePluginToken(receivedCtx)).To(MatchError(handlers.ErrAuth))
			})

			It("does not allow local plugins to authenticate without the local token", func() {
				req, err := http.NewRequest("POST", "/foo", nil)
				Expect(err).ToNot(HaveOccurred())

				req.Header.Set("Origin", "app://bar")

				var receivedCtx context.Context
				authHandler := auth.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					receivedCtx = r.Context()
				}))

				rw := httptest.NewRecorder()
				authHandler.ServeHTTP(rw, req)

				Expect(receivedCtx).ToNot(BeNil())

				Expect(auth.AuthorizePluginToken(receivedCtx)).To(MatchError(handlers.ErrAuth))
			})

			It("does not allow local plugins to authenticate with an incorrect token", func() {
				req, err := http.NewRequest("POST", "/foo", nil)
				Expect(err).ToNot(HaveOccurred())

				req.Header.Set("Origin", "app://bar")
				req.Header.Set("Authorization", "Bearer incorrect-local-token")

				var receivedCtx context.Context
				authHandler := auth.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					receivedCtx = r.Context()
				}))

				rw := httptest.NewRecorder()
				authHandler.ServeHTTP(rw, req)

				Expect(receivedCtx).ToNot(BeNil())

				Expect(auth.AuthorizePluginToken(receivedCtx)).To(MatchError(handlers.ErrAuth))

			})

			Context("when a local plugin is registered", func() {
				BeforeEach(func() {
					Expect(cp.AddPlugin("Foo Plugin", "https://localhost:12345/foo/plugin")).To(Succeed())
					Expect(auth.RefreshConfig()).To(Succeed())
				})

				It("allows local plugins to use its assigned API token instead", func() {
					req, err := http.NewRequest("POST", "/foo", nil)
					Expect(err).ToNot(HaveOccurred())

					req.Header.Set("Origin", "https://localhost")

					plugins := auth.GetRegisteredPlugins()
					pluginInfo := plugins["Foo Plugin"]
					apiToken := pluginInfo.APIToken

					req.Header.Set("Authorization", "Bearer "+apiToken)

					var receivedCtx context.Context
					authHandler := auth.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						receivedCtx = r.Context()
					}))

					rw := httptest.NewRecorder()
					authHandler.ServeHTTP(rw, req)

					Expect(receivedCtx).ToNot(BeNil())

					Expect(auth.AuthorizePluginToken(receivedCtx)).To(Succeed())
				})
			})
		})
	})

	Describe("AllowOriginFunc", func() {
		It("always allows cross-origin requests from file://", func() {
			Expect(auth.AllowOriginFunc("file:///D:/dummy/path/index.html")).To(BeTrue())
		})

		It("always allows cross-origin requests from app://", func() {
			Expect(auth.AllowOriginFunc("app://./index.html")).To(BeTrue())
		})

		It("always allows cross-origin requests from localhost", func() {
			Expect(auth.AllowOriginFunc("http://localhost:9001")).To(BeTrue())
		})
	})
})
