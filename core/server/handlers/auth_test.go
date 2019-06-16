package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/99designs/gqlgen/handler"

	"github.com/ff14wed/aetherometer/core/config"
	"github.com/ff14wed/aetherometer/core/server/handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {
	var (
		auth     *handlers.Auth
		adminOTP string
	)

	BeforeEach(func() {
		adminOTP = "one-time-password"
		cfg := config.Config{
			AdminOTP: adminOTP,
		}
		var err error
		auth, err = handlers.NewAuth(cfg, nil)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("CreateAdminToken", func() {
		It("successfully creates a token that passes admin validation", func() {
			ctx := handlers.ContextWithToken(adminOTP)
			token, err := auth.CreateAdminToken(ctx)
			Expect(err).ToNot(HaveOccurred())

			ctx = handlers.ContextWithToken(token)
			Expect(auth.AuthorizeAdminToken(ctx)).To(Succeed())
		})

		It("does not allow a token to be created twice", func() {
			ctx := handlers.ContextWithToken(adminOTP)
			_, err := auth.CreateAdminToken(ctx)
			Expect(err).ToNot(HaveOccurred())

			_, err = auth.CreateAdminToken(ctx)
			Expect(err).To(MatchError("No more admin tokens can be created."))
		})

		It("does not create a token if the OTP is incorrect", func() {
			ctx := handlers.ContextWithToken("incorrect")
			_, err := auth.CreateAdminToken(ctx)
			Expect(err).To(MatchError(handlers.AuthError))
		})

		It("allows creation of a token after a failed attempt", func() {
			ctx := handlers.ContextWithToken("incorrect")
			_, err := auth.CreateAdminToken(ctx)
			Expect(err).To(MatchError(handlers.AuthError))

			ctx = handlers.ContextWithToken(adminOTP)
			token, err := auth.CreateAdminToken(ctx)
			Expect(err).ToNot(HaveOccurred())

			ctx = handlers.ContextWithToken(token)
			Expect(auth.AuthorizeAdminToken(ctx)).To(Succeed())
		})
	})

	Describe("AuthorizeAdminToken", func() {
		It("rejects valid tokens that contain the incorrect admin ID", func() {
			ctx := handlers.ContextWithToken(adminOTP)
			oldToken, err := auth.CreateAdminToken(ctx)
			Expect(err).ToNot(HaveOccurred())

			auth.ResetOTPUsed()
			_, err = auth.CreateAdminToken(ctx)
			Expect(err).ToNot(HaveOccurred())

			ctx = handlers.ContextWithToken(oldToken)
			Expect(auth.AuthorizeAdminToken(ctx)).To(MatchError(handlers.AuthError))
		})

		It("rejects the token if it was signed with the wrong key", func() {
			altAuth, err := handlers.NewAuth(config.Config{AdminOTP: "altotp"}, nil)
			Expect(err).ToNot(HaveOccurred())
			ctx := handlers.ContextWithToken("altotp")
			token, err := altAuth.CreateAdminToken(ctx)
			Expect(err).ToNot(HaveOccurred())

			ctx = handlers.ContextWithToken(token)
			Expect(auth.AuthorizeAdminToken(ctx)).To(MatchError(handlers.AuthError))
		})

		Context("when the authorization token is sent via websockets", func() {
			var adminToken string

			BeforeEach(func() {
				cfg := config.Config{
					AdminOTP: adminOTP,
				}
				var err error
				auth, err = handlers.NewAuth(cfg, func(context.Context) handler.InitPayload {
					return handler.InitPayload{"Authorization": adminToken}
				})
				Expect(err).ToNot(HaveOccurred())
			})

			It("successfully validates the token", func() {
				ctx := handlers.ContextWithToken(adminOTP)
				var err error
				adminToken, err = auth.CreateAdminToken(ctx)
				Expect(err).ToNot(HaveOccurred())

				Expect(auth.AuthorizeAdminToken(context.Background())).To(Succeed())
			})

			It("rejects invalid tokens", func() {
				adminToken = "invalid-token"
				Expect(auth.AuthorizeAdminToken(context.Background())).To(MatchError(handlers.AuthError))
			})
		})
	})

	Describe("AddPlugin", func() {
		var adminToken string

		BeforeEach(func() {
			ctx := handlers.ContextWithToken(adminOTP)
			var err error
			adminToken, err = auth.CreateAdminToken(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("authorizes the plugin to access the API", func() {
			ctx := handlers.ContextWithToken(adminToken)
			apiToken, err := auth.AddPlugin(ctx, "https://example.com/foo/plugin")
			Expect(err).ToNot(HaveOccurred())

			ctx = handlers.ContextWithToken(apiToken)
			Expect(auth.AuthorizePluginToken(ctx)).To(Succeed())
			Expect(auth.AllowOriginFunc("https://example.com")).To(BeTrue())
		})

		It("returns an error if the plugin URL is invalid", func() {
			ctx := handlers.ContextWithToken(adminToken)
			apiToken, err := auth.AddPlugin(ctx, ":bad-url.com")
			Expect(err).To(BeAssignableToTypeOf(new(url.Error)))
			Expect(apiToken).To(BeEmpty())
		})

		It("returns an error if the plugin URL is missing the scheme", func() {
			ctx := handlers.ContextWithToken(adminToken)
			apiToken, err := auth.AddPlugin(ctx, "example.com/foo/plugin")
			Expect(err).To(MatchError("Could not parse plugin URL."))
			Expect(apiToken).To(BeEmpty())
		})

		Context("when the incorrect admin token has been provided", func() {
			BeforeEach(func() {
				altAuth, err := handlers.NewAuth(config.Config{AdminOTP: "altotp"}, nil)
				Expect(err).ToNot(HaveOccurred())
				ctx := handlers.ContextWithToken("altotp")
				adminToken, err = altAuth.CreateAdminToken(ctx)
				Expect(err).ToNot(HaveOccurred())
			})

			It("does not authorize the plugin", func() {
				ctx := handlers.ContextWithToken(adminToken)
				apiToken, err := auth.AddPlugin(ctx, "https://example.com/foo/plugin")
				Expect(err).To(MatchError(handlers.AuthError))
				Expect(apiToken).To(BeEmpty())

				Expect(auth.AllowOriginFunc("https://example.com")).To(BeFalse())
			})
		})
	})

	Describe("RemovePlugin", func() {
		var adminToken, apiToken string

		BeforeEach(func() {
			ctx := handlers.ContextWithToken(adminOTP)
			var err error
			adminToken, err = auth.CreateAdminToken(ctx)
			Expect(err).ToNot(HaveOccurred())

			ctx = handlers.ContextWithToken(adminToken)
			apiToken, err = auth.AddPlugin(ctx, "https://example.com/foo/plugin")
			Expect(err).ToNot(HaveOccurred())
		})

		It("revokes authorization for the plugin to access the API", func() {
			ctx := handlers.ContextWithToken(adminToken)
			removed, err := auth.RemovePlugin(ctx, apiToken)
			Expect(err).ToNot(HaveOccurred())
			Expect(removed).To(BeTrue())

			ctx = handlers.ContextWithToken(apiToken)
			Expect(auth.AuthorizePluginToken(ctx)).To(MatchError(handlers.AuthError))
			Expect(auth.AllowOriginFunc("https://example.com")).To(BeFalse())
		})

		It("removes the token if even if it doesn't exist or no longer exists", func() {
			ctx := handlers.ContextWithToken(adminToken)
			for i := 0; i < 3; i++ {
				removed, err := auth.RemovePlugin(ctx, apiToken)
				Expect(err).ToNot(HaveOccurred())
				Expect(removed).To(BeTrue())
			}

			ctx = handlers.ContextWithToken(apiToken)
			Expect(auth.AuthorizePluginToken(ctx)).To(MatchError(handlers.AuthError))
			Expect(auth.AllowOriginFunc("https://example.com")).To(BeFalse())
		})

		Context("when adding a second plugin with the same origin", func() {
			var altAPIToken string

			BeforeEach(func() {
				ctx := handlers.ContextWithToken(adminToken)
				var err error
				altAPIToken, err = auth.AddPlugin(ctx, "https://example.com/foo/plugin")
				Expect(err).ToNot(HaveOccurred())
			})

			It("continues to allow the origin when removing a single plugin", func() {
				ctx := handlers.ContextWithToken(adminToken)
				for i := 0; i < 3; i++ {
					removed, err := auth.RemovePlugin(ctx, apiToken)
					Expect(err).ToNot(HaveOccurred())
					Expect(removed).To(BeTrue())
				}

				ctx = handlers.ContextWithToken(apiToken)
				Expect(auth.AuthorizePluginToken(ctx)).To(MatchError(handlers.AuthError))

				ctx = handlers.ContextWithToken(altAPIToken)
				Expect(auth.AuthorizePluginToken(ctx)).To(Succeed())
				Expect(auth.AllowOriginFunc("https://example.com")).To(BeTrue())
			})

			It("no longer allows the origin when removing both plugins", func() {
				ctx := handlers.ContextWithToken(adminToken)

				removed, err := auth.RemovePlugin(ctx, apiToken)
				Expect(err).ToNot(HaveOccurred())
				Expect(removed).To(BeTrue())

				removed, err = auth.RemovePlugin(ctx, altAPIToken)
				Expect(err).ToNot(HaveOccurred())
				Expect(removed).To(BeTrue())

				ctx = handlers.ContextWithToken(apiToken)
				Expect(auth.AuthorizePluginToken(ctx)).To(MatchError(handlers.AuthError))

				ctx = handlers.ContextWithToken(altAPIToken)
				Expect(auth.AuthorizePluginToken(ctx)).To(MatchError(handlers.AuthError))
				Expect(auth.AllowOriginFunc("https://example.com")).To(BeFalse())
			})
		})

		Context("when the incorrect admin token has been provided", func() {
			BeforeEach(func() {
				altAuth, err := handlers.NewAuth(config.Config{AdminOTP: "altotp"}, nil)
				Expect(err).ToNot(HaveOccurred())
				ctx := handlers.ContextWithToken("altotp")
				adminToken, err = altAuth.CreateAdminToken(ctx)
				Expect(err).ToNot(HaveOccurred())
			})

			It("does not revoke authorization for the plugin", func() {
				ctx := handlers.ContextWithToken(adminToken)
				removed, err := auth.RemovePlugin(ctx, apiToken)
				Expect(err).To(MatchError(handlers.AuthError))
				Expect(removed).To(BeFalse())

				ctx = handlers.ContextWithToken(apiToken)
				Expect(auth.AuthorizePluginToken(ctx)).To(Succeed())
				Expect(auth.AllowOriginFunc("https://example.com")).To(BeTrue())
			})
		})

		Context("when the provided api token is invalid", func() {
			var altAPIToken string

			BeforeEach(func() {
				altAuth, err := handlers.NewAuth(config.Config{AdminOTP: "altotp"}, nil)
				Expect(err).ToNot(HaveOccurred())
				ctx := handlers.ContextWithToken("altotp")
				altAdminToken, err := altAuth.CreateAdminToken(ctx)
				Expect(err).ToNot(HaveOccurred())

				ctx = handlers.ContextWithToken(altAdminToken)
				altAPIToken, err = altAuth.AddPlugin(ctx, "https://example.com/foo/plugin")
				Expect(err).ToNot(HaveOccurred())
			})

			It("does not revoke authorization for the plugin", func() {
				ctx := handlers.ContextWithToken(adminToken)
				removed, err := auth.RemovePlugin(ctx, altAPIToken)
				Expect(err).To(MatchError(handlers.AuthError))
				Expect(removed).To(BeFalse())

				ctx = handlers.ContextWithToken(apiToken)
				Expect(auth.AuthorizePluginToken(ctx)).To(Succeed())
				Expect(auth.AllowOriginFunc("https://example.com")).To(BeTrue())
			})
		})
	})

	Describe("AuthorizePluginToken", func() {
		var adminToken string

		BeforeEach(func() {
			ctx := handlers.ContextWithToken(adminOTP)
			var err error
			adminToken, err = auth.CreateAdminToken(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("authorizes requests with the admin token to access the API", func() {
			ctx := handlers.ContextWithToken(adminToken)
			Expect(auth.AuthorizePluginToken(ctx)).To(Succeed())
		})

		It("rejects valid tokens that contain the incorrect admin ID", func() {
			auth.ResetOTPUsed()
			ctx := handlers.ContextWithToken(adminOTP)
			_, err := auth.CreateAdminToken(ctx)
			Expect(err).ToNot(HaveOccurred())

			ctx = handlers.ContextWithToken(adminToken)
			Expect(auth.AuthorizePluginToken(ctx)).To(MatchError(handlers.AuthError))
		})

		Context("when auth is disabled", func() {
			BeforeEach(func() {
				cfg := config.Config{
					AdminOTP:    adminOTP,
					DisableAuth: true,
				}
				var err error
				auth, err = handlers.NewAuth(cfg, nil)
				Expect(err).ToNot(HaveOccurred())
			})

			It("authorizes requests with no token to access the API", func() {
				Expect(auth.AuthorizePluginToken(context.Background())).To(Succeed())
			})
		})

		Context("when the authorization token is sent via websockets", func() {
			var apiToken string

			BeforeEach(func() {
				cfg := config.Config{
					AdminOTP: adminOTP,
				}
				var err error
				auth, err = handlers.NewAuth(cfg, func(context.Context) handler.InitPayload {
					return handler.InitPayload{"Authorization": apiToken}
				})
				Expect(err).ToNot(HaveOccurred())

				ctx := handlers.ContextWithToken(adminOTP)
				adminToken, err = auth.CreateAdminToken(ctx)
				Expect(err).ToNot(HaveOccurred())
			})

			It("successfully validates the token", func() {
				ctx := handlers.ContextWithToken(adminToken)
				var err error
				apiToken, err = auth.AddPlugin(ctx, "https://example.com/foo/plugin")
				Expect(err).ToNot(HaveOccurred())

				Expect(auth.AuthorizePluginToken(context.Background())).To(Succeed())
			})

			It("rejects invalid tokens", func() {
				apiToken = "invalid-token"
				Expect(auth.AuthorizePluginToken(context.Background())).To(MatchError(handlers.AuthError))
			})
		})
	})

	Describe("Handler", func() {
		It("adds the auth token to the request context", func() {
			req, err := http.NewRequest("POST", "/foo", nil)
			Expect(err).ToNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+adminOTP)

			var token string
			authHandler := auth.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				token, err = auth.CreateAdminToken(r.Context())
				Expect(err).ToNot(HaveOccurred())
			}))

			rw := httptest.NewRecorder()
			authHandler.ServeHTTP(rw, req)

			Expect(token).ToNot(BeEmpty())

			ctx := handlers.ContextWithToken(token)
			Expect(auth.AuthorizeAdminToken(ctx)).To(Succeed())
		})

		It("calls the provided next handler even if there is no Authorization header", func() {
			req, err := http.NewRequest("POST", "/foo", nil)
			Expect(err).ToNot(HaveOccurred())

			var token string
			authHandler := auth.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				token, err = auth.CreateAdminToken(r.Context())
			}))

			rw := httptest.NewRecorder()
			authHandler.ServeHTTP(rw, req)

			Expect(token).To(BeEmpty())
			Expect(err).To(MatchError(handlers.AuthError))
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
				ctx := handlers.ContextWithToken(adminOTP)
				var err error
				adminToken, err := auth.CreateAdminToken(ctx)
				Expect(err).ToNot(HaveOccurred())

				ctx = handlers.ContextWithToken(adminToken)
				_, err = auth.AddPlugin(ctx, "https://example.com/foo/plugin")
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
})
