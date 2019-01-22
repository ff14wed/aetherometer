package server_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/ff14wed/sibyl/backend/config"
	"github.com/ff14wed/sibyl/backend/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/thejerf/suture"
)

var _ = Describe("Server", func() {
	var (
		s      *server.Server
		buf    *gbytes.Buffer
		logger *log.Logger
	)
	BeforeEach(func() {
		buf = gbytes.NewBuffer()
		logger = log.New(buf, "", 0)
		cfg := config.Config{}
		s = server.New(cfg, logger)
	})

	Describe("Address", func() {
		It("returns nil if the server has not been started", func() {
			Expect(s.Address()).To(BeNil())
		})
	})

	Context("when the server has been started", func() {
		var supervisor *suture.Supervisor

		JustBeforeEach(func() {
			supervisor = suture.New("test-server", suture.Spec{
				Log: func(line string) {
					_, _ = GinkgoWriter.Write([]byte(line))
				},
				FailureThreshold: 1,
			})
			supervisor.ServeBackground()
			_ = supervisor.Add(s)
			s.WaitForStart()
		})

		AfterEach(func() {
			supervisor.Stop()
		})

		It("starts the server on some port", func() {
			Expect(s.Address().Port).ToNot(BeZero())
		})

		It("logs the port the server started on", func() {
			Expect(buf).To(gbytes.Say(`Server: running at `))
			Expect(buf).To(gbytes.Say(s.Address().String()))
		})

		Context("when no handlers have been added to the server", func() {
			It("returns a status 404 page", func() {
				url := fmt.Sprintf("http://%s/", s.Address().String())
				req, err := http.NewRequest("GET", url, nil)
				Expect(err).ToNot(HaveOccurred())

				resp, err := http.DefaultClient.Do(req)
				Expect(err).ToNot(HaveOccurred())

				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				respBytes, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Body.Close()).To(Succeed())
				Expect(string(respBytes)).To(Equal("404 page not found\n"))
			})
		})

		Context("when a handler is added to the server", func() {
			var reqValue uint64
			BeforeEach(func() {
				atomic.StoreUint64(&reqValue, 0)
				handler := http.HandlerFunc(func(r http.ResponseWriter, req *http.Request) {
					defer func() {
						r.WriteHeader(http.StatusOK)
						_, _ = r.Write([]byte("OK"))
					}()

					v := req.Header.Get("foo")
					if v == "" {
						return
					}

					val, err := strconv.ParseUint(v, 10, 64)
					if err != nil {
						return
					}

					atomic.StoreUint64(&reqValue, val)
				})
				s.AddHandler("/foo", handler)
			})

			Context("when making the request to the correct endpoint", func() {
				var url string
				JustBeforeEach(func() {
					url = fmt.Sprintf("http://%s/foo", s.Address().String())
				})

				It("returns the correct response", func() {
					req, err := http.NewRequest("GET", url, nil)
					Expect(err).ToNot(HaveOccurred())
					req.Header.Set("foo", "1234")

					resp, err := http.DefaultClient.Do(req)
					Expect(err).ToNot(HaveOccurred())

					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					respBytes, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())
					Expect(resp.Body.Close()).To(Succeed())
					Expect(string(respBytes)).To(Equal("OK"))
				})

				It("calls the handler", func() {
					req, err := http.NewRequest("GET", url, nil)
					Expect(err).ToNot(HaveOccurred())
					req.Header.Set("foo", "1234")

					_, err = http.DefaultClient.Do(req)
					Expect(err).ToNot(HaveOccurred())

					Expect(atomic.LoadUint64(&reqValue)).To(Equal(uint64(1234)))
				})
			})

			Context("when making the request to the wrong endpoint", func() {
				var url string
				JustBeforeEach(func() {
					url = fmt.Sprintf("http://%s/notfoo", s.Address().String())
				})

				It("returns a status 404 page", func() {
					req, err := http.NewRequest("GET", url, nil)
					Expect(err).ToNot(HaveOccurred())
					req.Header.Set("foo", "1234")

					resp, err := http.DefaultClient.Do(req)
					Expect(err).ToNot(HaveOccurred())

					Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
					respBytes, err := ioutil.ReadAll(resp.Body)
					Expect(err).ToNot(HaveOccurred())
					Expect(resp.Body.Close()).To(Succeed())
					Expect(string(respBytes)).To(Equal("404 page not found\n"))
				})
			})
		})

		Context("when the server is shut down", func() {
			JustBeforeEach(func() {
				supervisor.Stop()
			})

			It("no longer serves requests", func() {
				url := fmt.Sprintf("http://%s/", s.Address().String())
				req, err := http.NewRequest("GET", url, nil)
				Expect(err).ToNot(HaveOccurred())

				_, err = http.DefaultClient.Do(req)
				Expect(err).To(HaveOccurred())
			})

			It("logs that it is shutting down", func() {
				Expect(buf).To(gbytes.Say(`Server: stopping...`))
			})
		})
	})
})
