package config_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ff14wed/sibyl/backend/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HookConfig", func() {
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
				Adapters: config.Adapters{
					Hook: config.HookConfig{
						Enabled:      true,
						DLLPath:      dummyFile,
						FFXIVProcess: "ffxiv_dx11.exe",
					},
				},
			}
			Expect(c.Validate()).To(Succeed())
		})

		Describe("hook_dll", func() {
			var hookDLL string

			JustBeforeEach(func() {
				c = &config.Config{
					APIPort:  9000,
					DataPath: dummyPath,
					Maps: config.MapConfig{
						Cache: dummyPath,
					},
					Adapters: config.Adapters{
						Hook: config.HookConfig{
							Enabled: true,
							DLLPath: hookDLL,
						},
					},
				}
			})

			Context("when hook_dll is empty", func() {
				BeforeEach(func() {
					hookDLL = ""
				})
				It("errors", func() {
					Expect(c.Validate()).To(MatchError("config error in [adapters.hook]: dll_path must be provided"))
				})
			})

			Context("when hook_dll does not exist", func() {
				BeforeEach(func() {
					hookDLL = `Z:\foo\does\not\exist`
				})
				It("errors", func() {
					Expect(c.Validate()).To(MatchError(`config error in [adapters.hook]: dll_path file ("Z:\foo\does\not\exist") does not exist`))
				})
			})

			Context("when hook_dll is not a file", func() {
				BeforeEach(func() {
					hookDLL = dummyPath
				})
				It("errors", func() {
					Expect(c.Validate()).To(MatchError(fmt.Sprintf(`config error in [adapters.hook]: dll_path ("%s") must be a file`, dummyPath)))
				})
			})
		})

		Describe("ffxiv_process", func() {
			BeforeEach(func() {
				c = &config.Config{
					APIPort:  9000,
					DataPath: dummyPath,
					Maps: config.MapConfig{
						Cache: dummyPath,
					},
					Adapters: config.Adapters{
						Hook: config.HookConfig{
							Enabled: true,
							DLLPath: dummyFile,
						},
					},
				}
			})

			It("errors when empty", func() {
				Expect(c.Validate()).To(MatchError("config error in [adapters.hook]: ffxiv_process must be provided"))
			})
		})
	})
})
