package process_test

import (
	"time"

	"github.com/ff14wed/aetherometer/core/process"
	"github.com/ff14wed/aetherometer/core/process/processfakes"
	"github.com/thejerf/suture"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scanner", func() {
	var (
		pp     *processfakes.FakeEnumerator
		ticker chan time.Time

		scanner    *process.Scanner
		supervisor *suture.Supervisor
	)

	BeforeEach(func() {
		pp = new(processfakes.FakeEnumerator)
		pp.EnumerateProcessesStub = func() (map[uint32]string, error) {
			return map[uint32]string{
				1: "fooa.exe",
				2: "foob.exe",
				3: "fooc.exe",
			}, nil
		}
		ticker = make(chan time.Time)

		zapCfg := zap.NewDevelopmentConfig()
		zapCfg.OutputPaths = []string{}
		logger, err := zapCfg.Build()
		Expect(err).ToNot(HaveOccurred())

		scanner = process.NewScanner("foo", ticker, pp, 10, logger)

		supervisor = suture.New("test-scanner", suture.Spec{
			Log: func(line string) {
				_, _ = GinkgoWriter.Write([]byte(line))
			},
			FailureThreshold: 1,
		})
		supervisor.ServeBackground()
		_ = supervisor.Add(scanner)
	})

	AfterEach(func() {
		supervisor.Stop()
	})

	It("notifies clients of all running PIDs that match the string on startup", func() {
		var pid1, pid2, pid3 uint32
		Eventually(scanner.ProcessAddEventListener()).Should(Receive(&pid1))
		Eventually(scanner.ProcessAddEventListener()).Should(Receive(&pid2))
		Eventually(scanner.ProcessAddEventListener()).Should(Receive(&pid3))
		Expect([]uint32{pid1, pid2, pid3}).To(ConsistOf(
			uint32(1), uint32(2), uint32(3),
		))
	})

	Context("when the scanner has been running a while", func() {
		BeforeEach(func() {
			for i := 0; i < 3; i++ {
				Eventually(scanner.ProcessAddEventListener()).Should(Receive())
			}
		})

		It("notifies clients of any new PIDs that started running since", func() {
			pp.EnumerateProcessesStub = func() (map[uint32]string, error) {
				return map[uint32]string{
					1: "fooa.exe",
					2: "foob.exe",
					3: "fooc.exe",
					4: "food.exe",
				}, nil
			}
			ticker <- time.Now()
			var pid uint32
			Eventually(scanner.ProcessAddEventListener()).Should(Receive(&pid))
			Expect(pid).To(Equal(uint32(4)))
		})

		It("sends no notifications if nothing has changed", func() {
			ticker <- time.Now()
			Consistently(scanner.ProcessAddEventListener()).ShouldNot(Receive())
			Consistently(scanner.ProcessRemoveEventListener()).ShouldNot(Receive())
		})

		It("notifies clients of any PIDs that are no longer running", func() {
			pp.EnumerateProcessesStub = func() (map[uint32]string, error) {
				return map[uint32]string{1: "fooa.exe"}, nil
			}
			ticker <- time.Now()
			var pid1, pid2 uint32
			Eventually(scanner.ProcessRemoveEventListener()).Should(Receive(&pid1))
			Eventually(scanner.ProcessRemoveEventListener()).Should(Receive(&pid2))
			Expect([]uint32{pid1, pid2}).To(ConsistOf(
				uint32(2), uint32(3),
			))
		})

		It("notifies clients of both new running PIDs and no longer running PIDs", func() {
			pp.EnumerateProcessesStub = func() (map[uint32]string, error) {
				return map[uint32]string{
					1: "fooa.exe",
					2: "foob.exe",
					4: "food.exe",
				}, nil
			}
			ticker <- time.Now()
			var pid uint32
			Eventually(scanner.ProcessAddEventListener()).Should(Receive(&pid))
			Expect(pid).To(Equal(uint32(4)))
			Eventually(scanner.ProcessRemoveEventListener()).Should(Receive(&pid))
			Expect(pid).To(Equal(uint32(3)))
		})
	})
})
