package process_test

import (
	"errors"

	"github.com/ff14wed/aetherometer/core/process"
	"github.com/ff14wed/aetherometer/core/process/processfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListMatchingProcesses", func() {
	var pp *processfakes.FakeEnumerator

	BeforeEach(func() {
		pp = new(processfakes.FakeEnumerator)
		pp.EnumerateProcessesStub = func() (map[uint32]string, error) {
			return map[uint32]string{
				1: "fooa.exe",
				2: "foob.exe",
				3: "fooc.exe",
				4: "bara.exe",
				5: "BARB.exe",
				6: "barC.exe",
			}, nil
		}
	})

	It("returns the process ID list for which the process names match", func() {
		pids, err := process.ListMatchingProcesses("foo", pp)
		Expect(err).ToNot(HaveOccurred())
		Expect(pids).To(ConsistOf(uint32(1), uint32(2), uint32(3)))
	})

	It("returns matches in a case insensitive manner", func() {
		pids, err := process.ListMatchingProcesses("bar", pp)
		Expect(err).ToNot(HaveOccurred())
		Expect(pids).To(ConsistOf(uint32(4), uint32(5), uint32(6)))
	})

	It("returns an empty list if no process matches", func() {
		pids, err := process.ListMatchingProcesses("test", pp)
		Expect(err).ToNot(HaveOccurred())
		Expect(pids).To(BeEmpty())
	})

	Context("when the enumeration fails for some reason", func() {
		BeforeEach(func() {
			pp.EnumerateProcessesStub = func() (map[uint32]string, error) {
				return nil, errors.New("boom")
			}
		})

		It("returns an appropriate error", func() {
			_, err := process.ListMatchingProcesses("test", pp)
			Expect(err).To(MatchError("EnumerateProcesses error: boom"))
		})
	})
})
