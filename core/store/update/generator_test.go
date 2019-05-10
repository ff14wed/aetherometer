package update_test

import (
	"github.com/ff14wed/aetherometer/core/store/update"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generator", func() {
	var (
		testEnv = new(testVars)

		b         *xivnet.Block
		streamID  int
		generator update.Generator
	)

	BeforeEach(func() {
		*testEnv = genericSetup()
		b = testEnv.b
		streamID = testEnv.streamID
		generator = testEnv.generator

	})

	Context("when an ingress block arrives and the generator only process egress blocks", func() {
		BeforeEach(func() {
			movementData := &datatypes.Movement{Direction: 128}
			b.Data = movementData
		})

		It("does not generate an update", func() {
			u := generator.Generate(streamID, true, b)
			Expect(u).To(BeNil())
		})
	})

	Context("when an egress block arrives and the generator only process ingress blocks", func() {
		BeforeEach(func() {
			movementData := &datatypes.EgressMovement{Direction: 0}
			b.Data = movementData
		})

		It("does not generate an update", func() {
			u := generator.Generate(streamID, false, b)
			Expect(u).To(BeNil())
		})
	})

})
