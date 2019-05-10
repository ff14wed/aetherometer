package handlers_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"sync"

	"github.com/ff14wed/aetherometer/core/config"
	"github.com/ff14wed/aetherometer/core/server/handlers"
	"github.com/ff14wed/aetherometer/core/testhelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/ghttp"
	"go.uber.org/zap"
)

var _ = Describe("MapHandler", func() {
	var (
		apiServer *ghttp.Server
		cachePath string

		mapHandler http.Handler

		logBuf *testhelpers.LogBuffer
		once   sync.Once

		request *http.Request
	)

	BeforeEach(func() {
		once.Do(func() {
			logBuf = new(testhelpers.LogBuffer)
			err := zap.RegisterSink("maphandlertest", func(*url.URL) (zap.Sink, error) {
				return logBuf, nil
			})
			Expect(err).ToNot(HaveOccurred())
		})
		logBuf.Reset()

		zapCfg := zap.NewDevelopmentConfig()
		zapCfg.OutputPaths = []string{"maphandlertest://"}
		logger, err := zapCfg.Build()
		Expect(err).ToNot(HaveOccurred())

		apiServer = ghttp.NewServer()

		cachePath, err = ioutil.TempDir("", "map-handler")
		Expect(err).ToNot(HaveOccurred())

		mapHandler = handlers.NewMapHandler(
			"/prefix/",
			config.Config{
				Maps: config.MapConfig{
					Cache:   cachePath,
					APIPath: "http://" + apiServer.Addr(),
				},
			},
			logger,
		)

		request, err = http.NewRequest("GET", "/prefix/123", nil)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		apiServer.Close()
		os.RemoveAll(cachePath)
	})

	Context("when the map is found in the cache", func() {
		BeforeEach(func() {
			Expect(ioutil.WriteFile(
				path.Join(cachePath, "123.png"),
				[]byte("Image Bytes from Cache"),
				0777,
			)).To(Succeed())
		})

		It("serves the map from the cache", func() {
			rr := httptest.NewRecorder()
			mapHandler.ServeHTTP(rr, request)

			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Header().Get("Content-Type")).To(Equal("image/png"))
			Expect(rr.Header().Get("ETag")).To(Equal(`W/"map-123"`))
			Expect(rr.Header().Get("Cache-Control")).To(Equal("max-age=86400"))
			Expect(rr.Body.String()).To(Equal("Image Bytes from Cache"))
		})

		It("doesn't log anything", func() {
			rr := httptest.NewRecorder()
			mapHandler.ServeHTTP(rr, request)

			Consistently(logBuf).ShouldNot(gbytes.Say("map-handler"))
		})

		It("returns StatusNotModified only when the If-None-Match key has been previously sent as an ETag", func() {
			request.Header.Set("If-None-Match", `"abcd", W/"map-123"`)
			rr := httptest.NewRecorder()
			mapHandler.ServeHTTP(rr, request)

			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Header().Get("Content-Type")).To(Equal("image/png"))

			rr = httptest.NewRecorder()
			mapHandler.ServeHTTP(rr, request)
			Expect(rr.Code).To(Equal(http.StatusNotModified))
			Expect(rr.Body.String()).To(BeEmpty())
		})

		Context("when the If-None-Match key is recognized, but it doesn't match the map ID", func() {
			BeforeEach(func() {
				Expect(ioutil.WriteFile(
					path.Join(cachePath, "124.png"),
					[]byte("Dummy file"),
					0777,
				)).To(Succeed())
			})

			It("does not pass the condition for returning a StatusNotModified", func() {
				primingRequest, err := http.NewRequest("GET", "/prefix/124", nil)
				Expect(err).ToNot(HaveOccurred())

				rr := httptest.NewRecorder()
				mapHandler.ServeHTTP(rr, primingRequest)

				request.Header.Set("If-None-Match", `"abcd", W/"map-124"`)
				rr = httptest.NewRecorder()
				mapHandler.ServeHTTP(rr, request)

				Expect(rr.Code).To(Equal(http.StatusOK))
				Expect(rr.Header().Get("Content-Type")).To(Equal("image/png"))
			})
		})

		Context("when the jpg file also exists", func() {
			BeforeEach(func() {
				Expect(ioutil.WriteFile(
					path.Join(cachePath, "123.jpg"),
					[]byte("JPG bytes from Cache"),
					0777,
				)).To(Succeed())
			})

			It("prioritizes the png map from the cache", func() {
				rr := httptest.NewRecorder()
				mapHandler.ServeHTTP(rr, request)

				Expect(rr.Code).To(Equal(http.StatusOK))
				Expect(rr.Header().Get("Content-Type")).To(Equal("image/png"))
				Expect(rr.Body.String()).To(Equal("Image Bytes from Cache"))
			})
		})

		Context("when only the jpg file exists", func() {
			BeforeEach(func() {
				Expect(ioutil.WriteFile(
					path.Join(cachePath, "124.jpg"),
					[]byte("JPG Bytes from Cache"),
					0777,
				)).To(Succeed())

				var err error
				request, err = http.NewRequest("GET", "/prefix/124", nil)
				Expect(err).ToNot(HaveOccurred())
			})

			It("serves the jpg file from the cache", func() {
				rr := httptest.NewRecorder()
				mapHandler.ServeHTTP(rr, request)

				Expect(rr.Code).To(Equal(http.StatusOK))
				Expect(rr.Header().Get("Content-Type")).To(Equal("image/jpeg"))
				Expect(rr.Body.String()).To(Equal("JPG Bytes from Cache"))
			})
		})
	})

	Context("when the map is not found in the cache", func() {
		BeforeEach(func() {
			apiServer.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/Map/123"),
				ghttp.RespondWith(
					http.StatusOK,
					`{
						"ID":          123,
						"MapFilename": "\/m\/w1i1\/w1i1.01.jpg"
					}`,
				),
			))

			apiServer.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/m/w1i1/w1i1.01.jpg"),
				ghttp.RespondWith(
					http.StatusOK,
					"Image Bytes from API",
					http.Header{"Content-Type": []string{"image/jpeg"}},
				),
			))
		})

		It("downloads the map from the API and serves it", func() {
			rr := httptest.NewRecorder()
			mapHandler.ServeHTTP(rr, request)

			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Header().Get("Content-Type")).To(Equal("image/jpeg"))
			Expect(rr.Header().Get("ETag")).To(Equal(`W/"map-123"`))
			Expect(rr.Header().Get("Cache-Control")).To(Equal("max-age=86400"))
			Expect(rr.Body.String()).To(Equal("Image Bytes from API"))
		})

		It("saves the map in the cache", func() {
			_, err := os.Stat(path.Join(cachePath, "123.jpg"))
			Expect(err).To(HaveOccurred())
			Expect(os.IsNotExist(err)).To(BeTrue())

			rr := httptest.NewRecorder()
			mapHandler.ServeHTTP(rr, request)

			_, err = os.Stat(path.Join(cachePath, "123.jpg"))
			Expect(err).ToNot(HaveOccurred())
		})

		It("serves the map from the cache when requested again", func() {
			rr := httptest.NewRecorder()
			mapHandler.ServeHTTP(rr, request)

			// The fake apiServer will reject additional requests since we onl
			// added 2 one-off handlers to it.
			rr = httptest.NewRecorder()
			mapHandler.ServeHTTP(rr, request)

			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(rr.Header().Get("Content-Type")).To(Equal("image/jpeg"))
			Expect(rr.Body.String()).To(Equal("Image Bytes from API"))
		})

		It("doesn't log anything", func() {
			rr := httptest.NewRecorder()
			mapHandler.ServeHTTP(rr, request)

			Consistently(logBuf).ShouldNot(gbytes.Say("map-handler"))
		})

		Context("when the map is not found in the API", func() {
			BeforeEach(func() {
				apiServer.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/Map/123"),
					ghttp.RespondWith(http.StatusNotFound, "not found"),
				))
			})

			It("returns a 404", func() {
				rr := httptest.NewRecorder()
				mapHandler.ServeHTTP(rr, request)

				Expect(rr.Code).To(Equal(http.StatusNotFound))
			})

			It("logs a debug log with both cache and API errors", func() {
				rr := httptest.NewRecorder()
				mapHandler.ServeHTTP(rr, request)

				Eventually(logBuf).Should(gbytes.Say("DEBUG"))
				Eventually(logBuf).Should(gbytes.Say("Unable to find map in cache or API"))
				Eventually(logBuf).Should(gbytes.Say("cacheErr.*?not found"))
				Eventually(logBuf).Should(gbytes.Say("apiErr.*?not found"))
			})
		})

		Context("when the API returns an invalid response", func() {
			BeforeEach(func() {
				apiServer.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/Map/123"),
					ghttp.RespondWith(http.StatusOK, `{"ID":123}`),
				))
			})

			It("returns a 404", func() {
				rr := httptest.NewRecorder()
				mapHandler.ServeHTTP(rr, request)

				Expect(rr.Code).To(Equal(http.StatusNotFound))
			})

			It("logs a debug log with both cache and API errors", func() {
				rr := httptest.NewRecorder()
				mapHandler.ServeHTTP(rr, request)

				Eventually(logBuf).Should(gbytes.Say("DEBUG"))
				Eventually(logBuf).Should(gbytes.Say("Unable to find map in cache or API"))
				Eventually(logBuf).Should(gbytes.Say("cacheErr.*?not found"))
				Eventually(logBuf).Should(gbytes.Say(`apiErr.*?invalid response from API:.*123`))
			})
		})
	})

	Context("when the path is not prefixed with the appropriate path", func() {
		BeforeEach(func() {
			var err error
			request, err = http.NewRequest("GET", "/unknown/123", nil)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns a 404", func() {
			rr := httptest.NewRecorder()
			mapHandler.ServeHTTP(rr, request)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

		It("logs a debug log about the request path", func() {
			rr := httptest.NewRecorder()
			mapHandler.ServeHTTP(rr, request)

			Eventually(logBuf).Should(gbytes.Say("DEBUG"))
			Eventually(logBuf).Should(gbytes.Say("Error handling request path"))
			Eventually(logBuf).Should(gbytes.Say("path.*/unknown/123"))
		})
	})

	Context("when the mapID is not an integer", func() {
		BeforeEach(func() {
			var err error
			request, err = http.NewRequest("GET", "/prefix/foo", nil)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns a 404", func() {
			rr := httptest.NewRecorder()
			mapHandler.ServeHTTP(rr, request)

			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

		It("logs a debug log about the request path", func() {
			rr := httptest.NewRecorder()
			mapHandler.ServeHTTP(rr, request)

			Eventually(logBuf).Should(gbytes.Say("DEBUG"))
			Eventually(logBuf).Should(gbytes.Say("Error parsing map ID"))
			Eventually(logBuf).Should(gbytes.Say("path.*/prefix/foo"))
		})
	})
})
