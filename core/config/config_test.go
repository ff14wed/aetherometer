package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/ff14wed/sibyl/backend/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var (
		c         *config.Config
		dummyFile string
		dummyPath string
	)

	Describe("Validate", func() {
		BeforeEach(func() {
			var err error
			dummyFile, err = os.Executable()
			Expect(err).ToNot(HaveOccurred())
			dummyPath = filepath.Dir(dummyFile)
		})

		It("is successful on a correct config", func() {
			c = &config.Config{
				APIPort:  9000,
				DataPath: dummyPath,
				Maps: config.MapConfig{
					Cache: dummyPath,
				},
			}
			Expect(c.Validate()).To(Succeed())
		})

		Describe("api_port", func() {
			BeforeEach(func() {
				c = &config.Config{}
			})

			It("errors when zero", func() {
				Expect(c.Validate()).To(MatchError("config error: api_port must be provided"))
			})
		})

		Describe("Maps", func() {
			var mapsDir string

			JustBeforeEach(func() {
				c = &config.Config{
					APIPort:  9000,
					DataPath: dummyPath,
					Maps: config.MapConfig{
						Cache: mapsDir,
					},
				}
			})

			Context("when cache is empty", func() {
				BeforeEach(func() {
					mapsDir = ""
				})

				It("errors", func() {
					Expect(c.Validate()).To(MatchError("config error in [maps]: cache must be provided"))
				})
			})

			Context("when cache does not exist", func() {
				BeforeEach(func() {
					mapsDir = `Z:\foo\does\not\exist`
				})

				It("errors", func() {
					Expect(c.Validate()).To(MatchError(`config error in [maps]: cache directory ("Z:\foo\does\not\exist") does not exist`))
				})
			})

			Context("when cache is not a directory", func() {
				BeforeEach(func() {
					mapsDir = dummyFile
				})

				It("errors", func() {
					Expect(c.Validate()).To(MatchError(fmt.Sprintf(`config error in [maps]: cache ("%s") must be a directory`, dummyFile)))
				})
			})
		})

		Describe("data_path", func() {
			var dataDir string

			JustBeforeEach(func() {
				c = &config.Config{
					APIPort:  9000,
					DataPath: dataDir,
					Maps: config.MapConfig{
						Cache: dummyPath,
					},
				}
			})

			Context("when data_path is empty", func() {
				BeforeEach(func() {
					dataDir = ""
				})

				It("errors", func() {
					Expect(c.Validate()).To(MatchError("config error: data_path must be provided"))
				})
			})

			Context("when data_path does not exist", func() {
				BeforeEach(func() {
					dataDir = `Z:\foo\does\not\exist`
				})
				It("errors", func() {
					Expect(c.Validate()).To(MatchError(`config error: data_path directory ("Z:\foo\does\not\exist") does not exist`))
				})
			})

			Context("when data_path is not a directory", func() {
				BeforeEach(func() {
					dataDir = dummyFile
				})
				It("errors", func() {
					Expect(c.Validate()).To(MatchError(fmt.Sprintf(`config error: data_path ("%s") must be a directory`, dummyFile)))
				})
			})
		})
	})

	Describe("Adapters", func() {
		var a config.Adapters

		Describe("IsEnabled", func() {
			Context("when the adapter is not enabled", func() {
				BeforeEach(func() {
					a.Hook.Enabled = false
				})

				It("returns false if the adapter is not enabled", func() {
					Expect(a.IsEnabled("Hook")).To(BeFalse())
				})
			})

			Context("when the adapter is enabled", func() {
				BeforeEach(func() {
					a.Hook.Enabled = true
				})

				It("returns true", func() {
					Expect(a.IsEnabled("Hook")).To(BeTrue())
				})
			})

			Context("when the adapter has no Enabled option", func() {
				It("returns true", func() {
					Expect(a.IsEnabled("Test")).To(BeTrue())
				})
			})

			It("panics if the adapter config does not exist", func() {
				var panicMsg interface{}
				Expect(func() {
					defer func() {
						if err := recover(); err != nil {
							panicMsg = err
							panic(panicMsg)
						}
					}()
					_ = a.IsEnabled("Unknown")
				}).To(Panic())
				Expect(panicMsg).To(Equal("ERROR: Adapter config for Unknown does not exist"))
			})
		})
	})

	Describe("toml.Decode", func() {
		var (
			input string
		)

		BeforeEach(func() {
			lines := []string{
				`api_port = 9000`,
				`data_path = "dummy-path"`,
				`[maps]`,
				`cache = "some-map-dir"`,
				`api_path = "www.maps.com"`,
				`[adapters.hook]`,
				`enabled = true`,
			}
			input = strings.Join(lines, "\n")

			c = &config.Config{
				APIPort:  9000,
				DataPath: "dummy-path",
				Maps: config.MapConfig{
					Cache:   "some-map-dir",
					APIPath: "www.maps.com",
				},
				Adapters: config.Adapters{
					Hook: config.HookConfig{
						Enabled: true,
					},
				},
			}
		})

		It("decodes successfully from TOML", func() {
			var cfg config.Config
			_, err := toml.Decode(input, &cfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(cfg).To(Equal(*c))
		})
	})
})
