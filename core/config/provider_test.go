package config_test

import (
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/thejerf/suture"
	"go.uber.org/zap"

	"github.com/ff14wed/aetherometer/core/config"
	"github.com/ff14wed/aetherometer/core/testhelpers"
)

var _ = Describe("Provider", func() {
	var (
		logBuf *testhelpers.LogBuffer
		once   sync.Once

		cp         *config.Provider
		configFile string

		supervisor *suture.Supervisor
	)

	BeforeEach(func() {
		once.Do(func() {
			logBuf = new(testhelpers.LogBuffer)
			err := zap.RegisterSink("configprovidertest", func(*url.URL) (zap.Sink, error) {
				return logBuf, nil
			})
			Expect(err).ToNot(HaveOccurred())
		})
		logBuf.Reset()
		zapCfg := zap.NewDevelopmentConfig()
		zapCfg.OutputPaths = []string{"configprovidertest://"}
		logger, err := zapCfg.Build()
		Expect(err).ToNot(HaveOccurred())

		f, err := ioutil.TempFile("", "configprovidertest")
		Expect(err).ToNot(HaveOccurred())
		configFile = f.Name()
		Expect(f.Close()).To(Succeed())
		Expect(os.RemoveAll(configFile)).To(Succeed())

		defaultCfg := config.Config{
			APIPort: 9000,
			Sources: config.Sources{
				DataPath: os.TempDir(),
				Maps: config.MapConfig{
					Cache:   os.TempDir(),
					APIPath: "www.maps.com",
				},
			},
			Adapters: config.Adapters{
				Hook: config.HookConfig{
					Enabled: false,
				},
			},
		}
		cp = config.NewProvider(configFile, defaultCfg, logger)

		supervisor = suture.New("test-configprovider", suture.Spec{
			Log: func(line string) {
				_, _ = GinkgoWriter.Write([]byte(line))
			},
			FailureThreshold: 1,
		})
		supervisor.ServeBackground()
	})

	JustBeforeEach(func() {
		_ = supervisor.Add(cp)
	})

	AfterEach(func() {
		supervisor.Stop()
		_ = os.RemoveAll(configFile)
	})

	It(`logs "Running" on startup`, func() {
		Eventually(logBuf).Should(gbytes.Say("config-provider.*Running"))
	})

	Describe("WaitUntilReady", func() {
		It("blocks until the provider is running", func() {
			cp.WaitUntilReady()
			Expect(logBuf).Should(gbytes.Say("config-provider.*Running"))
		})
	})

	It(`logs "Stopping..." on shutdown`, func() {
		supervisor.Stop()
		Eventually(logBuf).Should(gbytes.Say("config-provider.*Stopping..."))
	})

	It("writes the default config to the config file path if it does not already exist", func() {
		Eventually(logBuf).Should(gbytes.Say("config-provider.*Writing default config"))
		cp.WaitUntilReady()
		configBytes, err := ioutil.ReadFile(configFile)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(configBytes)).To(ContainSubstring("api_port = 9000"))
		Expect(string(configBytes)).To(ContainSubstring("www.maps.com"))
		Expect(string(configBytes)).To(ContainSubstring("enabled = false"))
	})

	Context("if the config file already exists", func() {
		BeforeEach(func() {
			lines := []string{
				`api_port = 9000`,
				`[sources]`,
				`data_path = "/tmp"`,
				`maps.cache = "/tmp"`,
				`maps.api_path = "www.maps.com"`,
				`[adapters.hook]`,
				`enabled = false`,
				`[plugins]`,
				`"My Plugin" = "https://foo.com/my/plugin"`,
			}
			configString := strings.Join(lines, "\n")
			Expect(ioutil.WriteFile(configFile, []byte(configString), 0644)).To(Succeed())
		})

		It("does not overwrite the existing file", func() {
			Consistently(logBuf).ShouldNot(gbytes.Say("config-provider.*Writing default config"))
			cp.WaitUntilReady()
			cfg := cp.Config()
			Expect(cfg.Plugins).To(HaveKeyWithValue("My Plugin", "https://foo.com/my/plugin"))
		})
	})

	It("updates the saved config upon config file change", func() {
		cp.WaitUntilReady()
		originalCfg := cp.Config()

		sub, _ := cp.NotifyHub.Subscribe()

		Expect(appendToFile(configFile, "[plugins]\n"+`"Some Plugin" = "https://bar.com/some/plugin"`+"\n")).To(Succeed())

		Eventually(logBuf).Should(gbytes.Say("config-provider.*Detected config file change"))
		Eventually(logBuf).Should(gbytes.Say("config-provider.*Successfully applied config change"))
		cfg := cp.Config()
		Expect(cfg.Plugins).To(HaveKeyWithValue("Some Plugin", "https://bar.com/some/plugin"))

		Expect(sub).To(Receive(), "Config Provider should emit events upon config file change")

		Expect(cfg).ToNot(Equal(originalCfg))
	})

	Describe("AddPlugin", func() {
		It("adds the plugin and syncs changes to the config to disk", func() {
			cp.WaitUntilReady()

			sub, _ := cp.NotifyHub.Subscribe()

			Expect(cp.AddPlugin("Other Plugin", "https://foo.com/bar/plugin")).To(Succeed())

			Expect(sub).To(Receive(), "Config Provider should emit an event when a plugin is added")

			cfg := cp.Config()
			Expect(cfg.Plugins).To(HaveKeyWithValue("Other Plugin", "https://foo.com/bar/plugin"))

			configBytes, err := ioutil.ReadFile(configFile)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(configBytes)).To(ContainSubstring(`"Other Plugin" = "https://foo.com/bar/plugin"`))

			Consistently(logBuf).ShouldNot(gbytes.Say("config-provider.*Detected config file change"))
		})

		It("does not allow duplicate plugins to be added", func() {
			cp.WaitUntilReady()

			Expect(cp.AddPlugin("Other Plugin", "https://foo.com/bar/plugin")).To(Succeed())
			Expect(cp.AddPlugin("Other Plugin", "https://foo.com/bar/plugin")).To(MatchError(`plugin "Other Plugin" already exists`))
			Expect(cp.AddPlugin("Some Plugin", "https://foo.com/bar/plugin")).To(Succeed())
		})

		It("respects file changes", func() {
			cp.WaitUntilReady()

			Expect(cp.AddPlugin("Other Plugin", "https://foo.com/bar/plugin")).To(Succeed())
			Consistently(logBuf).ShouldNot(gbytes.Say("config-provider.*Detected config file change"))

			Expect(appendToFile(configFile, `  "Some Plugin" = "https://bar.com/some/plugin"`+"\n")).To(Succeed())

			Eventually(logBuf).Should(gbytes.Say("config-provider.*Detected config file change"))
			Eventually(logBuf).Should(gbytes.Say("config-provider.*Successfully applied config change"))

			Expect(cp.AddPlugin("Some Plugin", "https://something.com/foo/plugin")).To(MatchError(`plugin "Some Plugin" already exists`))
		})

		It("mutates the config without mutating references to it", func() {
			cp.WaitUntilReady()

			Expect(cp.AddPlugin("Some Plugin", "https://foo.com/some/plugin")).To(Succeed())

			cfg := cp.Config()
			Expect(cfg.Plugins).To(HaveKeyWithValue("Some Plugin", "https://foo.com/some/plugin"))
			Expect(len(cfg.Plugins)).To(Equal(1))

			Expect(cp.AddPlugin("Other Plugin", "https://foo.com/other/plugin")).To(Succeed())

			Expect(cfg.Plugins).ToNot(HaveKeyWithValue("Other Plugin", "https://foo.com/other/plugin"))
			Expect(len(cfg.Plugins)).To(Equal(1))
		})
	})

	Describe("RemovePlugin", func() {
		It("removes the plugin and syncs changes to the config to disk", func() {
			cp.WaitUntilReady()

			sub, _ := cp.NotifyHub.Subscribe()

			lines := []string{
				`api_port = 9000`,
				`[sources]`,
				`data_path = "/tmp"`,
				`maps.cache = "/tmp"`,
				`maps.api_path = "www.maps.com"`,
				`[adapters.hook]`,
				`enabled = false`,
				`[plugins]`,
				`"My Plugin" = "https://foo.com/my/plugin"`,
			}
			configString := strings.Join(lines, "\n")
			Expect(ioutil.WriteFile(configFile, []byte(configString), 0644)).To(Succeed())

			Eventually(logBuf).Should(gbytes.Say("config-provider.*Detected config file change"))
			Eventually(logBuf).Should(gbytes.Say("config-provider.*Successfully applied config change"))

			Expect(sub).To(Receive(), "Config Provider should emit an event when the config file is changed")

			Expect(cp.RemovePlugin("My Plugin")).To(Succeed())

			Expect(sub).To(Receive(), "Config Provider should emit an event when a plugin is removed")

			cfg := cp.Config()
			Expect(cfg.Plugins).ToNot(HaveKeyWithValue("My Plugin", "https://foo.com/my/plugin"))

			configBytes, err := ioutil.ReadFile(configFile)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(configBytes)).ToNot(ContainSubstring(`My Plugin`))

			Consistently(logBuf).ShouldNot(gbytes.Say("config-provider.*Detected config file change"))
		})

		It("does nothing if the plugin does not exist", func() {
			cp.WaitUntilReady()

			Expect(cp.AddPlugin("Other Plugin", "https://foo.com/bar/plugin")).To(Succeed())
			Expect(cp.RemovePlugin("Other Plugin")).To(Succeed())
			Expect(cp.RemovePlugin("Other Plugin")).To(Succeed())

			Consistently(logBuf).ShouldNot(gbytes.Say("config-provider.*Detected config file change"))
		})

		It("does nothing if no plugins were ever added", func() {
			cp.WaitUntilReady()

			Expect(cp.RemovePlugin("Other Plugin")).To(Succeed())

			Consistently(logBuf).ShouldNot(gbytes.Say("config-provider.*Detected config file change"))
		})

		It("respects file changes", func() {
			cp.WaitUntilReady()

			Expect(cp.AddPlugin("Some Plugin", "https://foo.com/bar/plugin")).To(Succeed())
			Expect(cp.RemovePlugin("Some Plugin")).To(Succeed())
			Consistently(logBuf).ShouldNot(gbytes.Say("config-provider.*Detected config file change"))

			Expect(appendToFile(configFile, `  "Some Plugin" = "https://bar.com/some/plugin"`+"\n")).To(Succeed())

			Eventually(logBuf).Should(gbytes.Say("config-provider.*Detected config file change"))
			Eventually(logBuf).Should(gbytes.Say("config-provider.*Successfully applied config change"))

			Expect(cp.AddPlugin("Some Plugin", "https://something.com/foo/plugin")).To(MatchError(`plugin "Some Plugin" already exists`))
		})

		It("mutates the config without mutating references to it", func() {
			cp.WaitUntilReady()

			Expect(cp.AddPlugin("Some Plugin", "https://foo.com/some/plugin")).To(Succeed())

			cfg := cp.Config()
			Expect(cfg.Plugins).To(HaveKeyWithValue("Some Plugin", "https://foo.com/some/plugin"))
			Expect(len(cfg.Plugins)).To(Equal(1))

			Expect(cp.RemovePlugin("Some Plugin")).To(Succeed())

			Expect(cfg.Plugins).To(HaveKeyWithValue("Some Plugin", "https://foo.com/some/plugin"))
			Expect(len(cfg.Plugins)).To(Equal(1))
		})
	})

	It("doesn't watch files anymore after shutdown", func() {
		sub, _ := cp.NotifyHub.Subscribe()

		supervisor.Stop()
		Eventually(logBuf).Should(gbytes.Say("config-provider.*Stopping..."))

		Expect(appendToFile(configFile, "[plugins]\n"+`"Some Plugin" = "https://bar.com/some/plugin"`+"\n")).To(Succeed())

		Consistently(logBuf).ShouldNot(gbytes.Say("config-provider.*Detected config file change"))

		Expect(sub).ToNot(Receive(), "Config Provider should not emit events after it's stopped")

		cfg := cp.Config()
		Expect(cfg.Plugins).To(BeEmpty())
	})
})

func appendToFile(filename string, stringToAppend string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(stringToAppend)
	if err != nil {
		return err
	}
	return nil
}
