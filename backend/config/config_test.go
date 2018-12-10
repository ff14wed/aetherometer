package config_test

import (
	"fmt"
	"os"
	"path/filepath"

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
				HookDLL:      dummyFile,
				FFXIVProcess: "ffxiv_dx11.exe",
				APIPort:      9000,
				Sources: config.SourceDirs{
					MapsDir: dummyPath,
					DataDir: dummyPath,
				},
			}
			Expect(c.Validate()).To(Succeed())
		})
		Describe("hook_dll", func() {
			var hookDLL string
			JustBeforeEach(func() {
				c = &config.Config{
					HookDLL: hookDLL,
				}
			})
			Context("when hook_dll is empty", func() {
				BeforeEach(func() {
					hookDLL = ""
				})
				It("errors", func() {
					Expect(c.Validate()).To(MatchError("config error: hook_dll must be provided"))
				})
			})
			Context("when hook_dll does not exist", func() {
				BeforeEach(func() {
					hookDLL = `Z:\foo\does\not\exist`
				})
				It("errors", func() {
					Expect(c.Validate()).To(MatchError(`config error: hook_dll file ("Z:\foo\does\not\exist") does not exist`))
				})
			})
			Context("when hook_dll is not a file", func() {
				BeforeEach(func() {
					hookDLL = dummyPath
				})
				It("errors", func() {
					Expect(c.Validate()).To(MatchError(fmt.Sprintf(`config error: hook_dll ("%s") must be a file`, dummyPath)))
				})
			})
		})
		Describe("ffxiv_process", func() {
			BeforeEach(func() {
				c = &config.Config{
					HookDLL: dummyFile,
				}
			})
			It("errors when empty", func() {
				Expect(c.Validate()).To(MatchError("config error: ffxiv_process must be provided"))
			})
		})
		Describe("api_port", func() {
			BeforeEach(func() {
				c = &config.Config{
					HookDLL:      dummyFile,
					FFXIVProcess: "ffxiv_dx11.exe",
				}
			})
			It("errors when zero", func() {
				Expect(c.Validate()).To(MatchError("config error: api_port must be provided"))
			})
		})
		Describe("maps_dir", func() {
			var mapsDir string
			JustBeforeEach(func() {
				c = &config.Config{
					HookDLL:      dummyFile,
					FFXIVProcess: "ffxiv_dx11.exe",
					APIPort:      9000,
					Sources: config.SourceDirs{
						MapsDir: mapsDir,
					},
				}
			})
			Context("when maps_dir is empty", func() {
				BeforeEach(func() {
					mapsDir = ""
				})
				It("errors", func() {
					Expect(c.Validate()).To(MatchError("config error in [sources]: maps_dir must be provided"))
				})
			})
			Context("when maps_dir does not exist", func() {
				BeforeEach(func() {
					mapsDir = `Z:\foo\does\not\exist`
				})
				It("errors", func() {
					Expect(c.Validate()).To(MatchError(`config error in [sources]: maps_dir directory ("Z:\foo\does\not\exist") does not exist`))
				})
			})
			Context("when maps_dir is not a directory", func() {
				BeforeEach(func() {
					mapsDir = dummyFile
				})
				It("errors", func() {
					Expect(c.Validate()).To(MatchError(fmt.Sprintf(`config error in [sources]: maps_dir ("%s") must be a directory`, dummyFile)))
				})
			})
		})
		Describe("data_dir", func() {
			var dataDir string
			JustBeforeEach(func() {
				c = &config.Config{
					HookDLL:      dummyFile,
					FFXIVProcess: "ffxiv_dx11.exe",
					APIPort:      9000,
					Sources: config.SourceDirs{
						MapsDir: dummyPath,
						DataDir: dataDir,
					},
				}
			})
			Context("when data_dir is empty", func() {
				BeforeEach(func() {
					dataDir = ""
				})
				It("errors", func() {
					Expect(c.Validate()).To(MatchError("config error in [sources]: data_dir must be provided"))
				})
			})
			Context("when data_dir does not exist", func() {
				BeforeEach(func() {
					dataDir = `Z:\foo\does\not\exist`
				})
				It("errors", func() {
					Expect(c.Validate()).To(MatchError(`config error in [sources]: data_dir directory ("Z:\foo\does\not\exist") does not exist`))
				})
			})
			Context("when data_dir is not a directory", func() {
				BeforeEach(func() {
					dataDir = dummyFile
				})
				It("errors", func() {
					Expect(c.Validate()).To(MatchError(fmt.Sprintf(`config error in [sources]: data_dir ("%s") must be a directory`, dummyFile)))
				})
			})
		})
	})
})
