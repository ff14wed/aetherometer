package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/ff14wed/aetherometer/core/config"

	. "github.com/onsi/ginkgo/v2"
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
				APIPort: 9000,
				Sources: config.Sources{
					DataPath: dummyPath,
					Maps: config.MapConfig{
						Cache: dummyPath,
					},
				},
			}
			Expect(c.Validate()).To(Succeed())
		})

		Describe("api_port", func() {
			BeforeEach(func() {
				c = &config.Config{
					Sources: config.Sources{
						DataPath: dummyPath,
						Maps: config.MapConfig{
							Cache: dummyPath,
						},
					},
				}
			})

			It("does not error when zero", func() {
				Expect(c.Validate()).To(Succeed())
			})
		})

		Describe("data_path", func() {
			var dataDir string

			JustBeforeEach(func() {
				c = &config.Config{
					APIPort: 9000,
					Sources: config.Sources{
						DataPath: dataDir,
					},
				}
			})

			Context("when data_path is empty", func() {
				BeforeEach(func() {
					dataDir = ""
				})

				It("errors", func() {
					Expect(c.Validate()).To(MatchError("config error in [sources]: data_path must be provided"))
				})
			})

			Context("when data_path does not exist", func() {
				BeforeEach(func() {
					dataDir = `Z:\foo\does\not\exist`
				})
				It("errors", func() {
					Expect(c.Validate()).To(MatchError(`config error in [sources]: data_path directory ("Z:\foo\does\not\exist") does not exist`))
				})
			})

			Context("when data_path is not a directory", func() {
				BeforeEach(func() {
					dataDir = dummyFile
				})
				It("errors", func() {
					Expect(c.Validate()).To(MatchError(fmt.Sprintf(`config error in [sources]: data_path ("%s") must be a directory`, dummyFile)))
				})
			})
		})

		Describe("Maps", func() {
			var mapsDir string

			JustBeforeEach(func() {
				c = &config.Config{
					APIPort: 9000,
					Sources: config.Sources{
						DataPath: dummyPath,
						Maps: config.MapConfig{
							Cache: mapsDir,
						},
					},
				}
			})

			Context("when cache is empty", func() {
				BeforeEach(func() {
					mapsDir = ""
				})

				It("errors", func() {
					Expect(c.Validate()).To(MatchError("config error in [sources.maps]: cache must be provided"))
				})
			})

			Context("when cache does not exist", func() {
				BeforeEach(func() {
					mapsDir = `Z:\foo\does\not\exist`
				})

				It("errors", func() {
					Expect(c.Validate()).To(MatchError(`config error in [sources.maps]: cache directory ("Z:\foo\does\not\exist") does not exist`))
				})
			})

			Context("when cache is not a directory", func() {
				BeforeEach(func() {
					mapsDir = dummyFile
				})

				It("errors", func() {
					Expect(c.Validate()).To(MatchError(fmt.Sprintf(`config error in [sources.maps]: cache ("%s") must be a directory`, dummyFile)))
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
					Expect(a.IsEnabled("test")).To(BeTrue())
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
				`disable_auth = false`,
				`local_token = "some-token"`,
				`auto_update = true`,
				`[sources]`,
				`data_path = "dummy-path"`,
				`maps.cache = "some-map-dir"`,
				`maps.api_path = "www.maps.com"`,
				`[adapters.hook]`,
				`enabled = true`,
				`[plugins]`,
				`"My Plugin" = "https://foo.com/my/plugin"`,
				`"Other Plugin" = "https://bar.com/other/plugin"`,
			}
			input = strings.Join(lines, "\n")

			c = &config.Config{
				APIPort:     9000,
				DisableAuth: false,
				LocalToken:  "some-token",
				AutoUpdate:  true,
				Sources: config.Sources{
					DataPath: "dummy-path",
					Maps: config.MapConfig{
						Cache:   "some-map-dir",
						APIPath: "www.maps.com",
					},
				},
				Adapters: config.Adapters{
					Hook: config.HookConfig{
						Enabled: true,
					},
				},
				Plugins: map[string]string{
					"My Plugin":    "https://foo.com/my/plugin",
					"Other Plugin": "https://bar.com/other/plugin",
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
