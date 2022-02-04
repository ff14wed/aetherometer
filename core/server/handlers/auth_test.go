package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
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

		cfg := config.Config{}
		auth, err = handlers.NewAuth(cfg, nil, logger)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("AddPlugin", func() {
		It("authorizes the plugin to access the API", func() {
			apiToken, err := auth.AddPlugin("https://example.com/foo/plugin")
			Expect(err).ToNot(HaveOccurred())

			ctx := handlers.ContextWithToken(apiToken)
			Expect(auth.AuthorizePluginToken(ctx)).To(Succeed())
			Expect(auth.AllowOriginFunc("https://example.com")).To(BeTrue())

			Eventually(logBuf).Should(gbytes.Say("DEBUG"))
			Eventually(logBuf).Should(gbytes.Say("auth-handler"))
			Eventually(logBuf).Should(gbytes.Say("Added Plugin"))
			Eventually(logBuf).Should(gbytes.Say("https://example.com.*id"))
		})

		It("returns an error if the plugin URL is invalid", func() {
			apiToken, err := auth.AddPlugin(":bad-url.com")
			Expect(err).To(BeAssignableToTypeOf(new(url.Error)))
			Expect(apiToken).To(BeEmpty())
		})

		It("returns an error if the plugin URL is missing the scheme", func() {
			apiToken, err := auth.AddPlugin("example.com/foo/plugin")
			Expect(err).To(MatchError("could not parse plugin URL"))
			Expect(apiToken).To(BeEmpty())
		})
	})

	Describe("RemovePlugin", func() {
		var apiToken string

		BeforeEach(func() {
			var err error
			apiToken, err = auth.AddPlugin("https://example.com/foo/plugin")
			Expect(err).ToNot(HaveOccurred())
		})

		It("revokes authorization for the plugin to access the API", func() {
			removed, err := auth.RemovePlugin(apiToken)
			Expect(err).ToNot(HaveOccurred())
			Expect(removed).To(BeTrue())

			ctx := handlers.ContextWithToken(apiToken)
			Expect(auth.AuthorizePluginToken(ctx)).To(MatchError(handlers.ErrAuth))
			Expect(auth.AllowOriginFunc("https://example.com")).To(BeFalse())

			Eventually(logBuf).Should(gbytes.Say("DEBUG"))
			Eventually(logBuf).Should(gbytes.Say("auth-handler"))
			Eventually(logBuf).Should(gbytes.Say("Deleted Plugin.*id"))
		})

		It("removes the token if even if it doesn't exist or no longer exists", func() {
			for i := 0; i < 3; i++ {
				removed, err := auth.RemovePlugin(apiToken)
				Expect(err).ToNot(HaveOccurred())
				Expect(removed).To(BeTrue())
			}

			ctx := handlers.ContextWithToken(apiToken)
			Expect(auth.AuthorizePluginToken(ctx)).To(MatchError(handlers.ErrAuth))
			Expect(auth.AllowOriginFunc("https://example.com")).To(BeFalse())
		})

		Context("when adding a second plugin with the same origin", func() {
			var altAPIToken string

			BeforeEach(func() {
				var err error
				altAPIToken, err = auth.AddPlugin("https://example.com/foo/plugin")
				Expect(err).ToNot(HaveOccurred())
			})

			It("continues to allow the origin when removing a single plugin", func() {
				for i := 0; i < 3; i++ {
					removed, err := auth.RemovePlugin(apiToken)
					Expect(err).ToNot(HaveOccurred())
					Expect(removed).To(BeTrue())
				}

				ctx := handlers.ContextWithToken(apiToken)
				Expect(auth.AuthorizePluginToken(ctx)).To(MatchError(handlers.ErrAuth))

				ctx = handlers.ContextWithToken(altAPIToken)
				Expect(auth.AuthorizePluginToken(ctx)).To(Succeed())
				Expect(auth.AllowOriginFunc("https://example.com")).To(BeTrue())
			})

			It("no longer allows the origin when removing both plugins", func() {
				removed, err := auth.RemovePlugin(apiToken)
				Expect(err).ToNot(HaveOccurred())
				Expect(removed).To(BeTrue())

				removed, err = auth.RemovePlugin(altAPIToken)
				Expect(err).ToNot(HaveOccurred())
				Expect(removed).To(BeTrue())

				ctx := handlers.ContextWithToken(apiToken)
				Expect(auth.AuthorizePluginToken(ctx)).To(MatchError(handlers.ErrAuth))

				ctx = handlers.ContextWithToken(altAPIToken)
				Expect(auth.AuthorizePluginToken(ctx)).To(MatchError(handlers.ErrAuth))
				Expect(auth.AllowOriginFunc("https://example.com")).To(BeFalse())
			})
		})

		Context("when the provided api token is invalid", func() {
			var altAPIToken string

			BeforeEach(func() {
				altAuth, err := handlers.NewAuth(config.Config{}, nil, zap.NewNop())
				Expect(err).ToNot(HaveOccurred())
				altAPIToken, err = altAuth.AddPlugin("https://example.com/foo/plugin")
				Expect(err).ToNot(HaveOccurred())
			})

			It("does not revoke authorization for the plugin", func() {
				removed, err := auth.RemovePlugin(altAPIToken)
				Expect(err).To(MatchError(handlers.ErrAuth))
				Expect(removed).To(BeFalse())

				ctx := handlers.ContextWithToken(apiToken)
				Expect(auth.AuthorizePluginToken(ctx)).To(Succeed())
				Expect(auth.AllowOriginFunc("https://example.com")).To(BeTrue())
			})
		})
	})

	Describe("AuthorizePluginToken", func() {
		Context("when auth is disabled", func() {
			BeforeEach(func() {
				cfg := config.Config{
					DisableAuth: true,
				}
				var err error
				auth, err = handlers.NewAuth(cfg, nil, zap.NewNop())
				Expect(err).ToNot(HaveOccurred())
			})

			It("authorizes requests with no token to access the API", func() {
				Expect(auth.AuthorizePluginToken(context.Background())).To(Succeed())
			})
		})

		Context("when the authorization token is sent via websockets", func() {
			var apiToken string

			BeforeEach(func() {
				cfg := config.Config{}
				apiToken = ""

				var err error
				auth, err = handlers.NewAuth(cfg, func(context.Context) transport.InitPayload {
					return transport.InitPayload{"Authorization": apiToken}
				}, zap.NewNop())
				Expect(err).ToNot(HaveOccurred())
			})

			It("successfully validates the token", func() {
				var err error
				apiToken, err = auth.AddPlugin("https://example.com/foo/plugin")
				Expect(err).ToNot(HaveOccurred())

				Expect(auth.AuthorizePluginToken(context.Background())).To(Succeed())
			})

			It("rejects invalid tokens", func() {
				apiToken = "invalid-token"
				Expect(auth.AuthorizePluginToken(context.Background())).To(MatchError(handlers.ErrAuth))
			})
		})
	})

	Describe("Handler", func() {
		It("adds the auth token to the request context", func() {
			req, err := http.NewRequest("POST", "/foo", nil)
			Expect(err).ToNot(HaveOccurred())

			apiToken, err := auth.AddPlugin("https://example.com/foo/plugin")
			Expect(err).ToNot(HaveOccurred())

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
				_, err := auth.AddPlugin("https://example.com/foo/plugin")
				Expect(err).ToNot(HaveOccurred())
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
